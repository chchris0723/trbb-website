package handlers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"trbb/internal/models"
	"trbb/internal/services"
)

type AdminEventHandler struct {
	eventSvc *services.EventService
}

func NewAdminEventHandler(eventSvc *services.EventService) *AdminEventHandler {
	return &AdminEventHandler{eventSvc: eventSvc}
}

// GET /v1/admin/events
func (h *AdminEventHandler) ListEvents(c *gin.Context) {
	var in services.ListEventsInput
	_ = c.ShouldBindQuery(&in)
	result, err := h.eventSvc.ListAdmin(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢失敗"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GET /v1/admin/events/:id
func (h *AdminEventHandler) GetEvent(c *gin.Context) {
	id, err := parseAdminID(c, "id")
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }

	ev, err := h.eventSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "賽事不存在"})
		return
	}
	c.JSON(http.StatusOK, ev)
}

// POST /v1/admin/events
func (h *AdminEventHandler) CreateEvent(c *gin.Context) {
	var in services.EventInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	creatorID := mustAdminUserID(c)
	ev, err := h.eventSvc.Create(c.Request.Context(), in, creatorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "建立失敗: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "賽事已建立", "event": ev})
}

// PUT /v1/admin/events/:id
func (h *AdminEventHandler) UpdateEvent(c *gin.Context) {
	id, err := parseAdminID(c, "id")
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }

	var in services.EventInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ev, err := h.eventSvc.Update(c.Request.Context(), id, in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失敗"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已更新", "event": ev})
}

// DELETE /v1/admin/events/:id
func (h *AdminEventHandler) DeleteEvent(c *gin.Context) {
	id, err := parseAdminID(c, "id")
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }

	if err := h.eventSvc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "刪除失敗"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已刪除"})
}

// GET /v1/admin/events/:id/registrations
func (h *AdminEventHandler) ListRegistrations(c *gin.Context) {
	id, err := parseAdminID(c, "id")
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }

	list, err := h.eventSvc.ListRegistrations(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢失敗"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"registrations": list, "total": len(list)})
}

// PUT /v1/admin/events/:id/registrations/:regId
func (h *AdminEventHandler) UpdateRegistration(c *gin.Context) {
	regID, err := parseAdminID(c, "regId")
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }

	var body struct {
		services.RegistrationInput
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.eventSvc.UpdateRegistration(c.Request.Context(), regID, body.RegistrationInput, body.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失敗"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已更新"})
}

// GET /v1/admin/events/:id/registrations/export?format=csv|excel|pdf
func (h *AdminEventHandler) ExportRegistrations(c *gin.Context) {
	id, err := parseAdminID(c, "id")
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }

	ev, err := h.eventSvc.GetByID(c.Request.Context(), id)
	if err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "賽事不存在"}); return }

	list, err := h.eventSvc.ListRegistrations(c.Request.Context(), id)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢失敗"}); return }

	format := c.DefaultQuery("format", "csv")
	switch format {
	case "csv":
		exportCSV(c, ev, list)
	case "excel":
		exportExcel(c, ev, list)
	default:
		exportCSV(c, ev, list)
	}
}

// ── Export helpers ────────────────────────────────────────────

var regHeaders = []string{
	"序號", "報名狀態", "會員帳號", "顯示名稱",
	"中文姓名", "英文姓名", "身份證字號", "護照號碼",
	"性別", "出生年月日", "手機", "Email",
	"衣服尺寸", "飲食習慣", "通訊地址",
	"緊急聯絡人", "緊急聯絡電話", "緊急聯絡關係",
	"備註", "報名時間",
}

func regToRow(i int, r *models.EventRegistration) []string {
	statusMap := map[int]string{0: "待確認", 1: "已確認", 2: "已取消", 3: "已退款"}
	genderMap := map[int]string{1: "男", 2: "女", 3: "其他"}
	foodMap   := map[int]string{1: "葷", 2: "素", 3: "全素"}

	gender := ""
	if r.RegGender != nil { gender = genderMap[*r.RegGender] }
	food := ""
	if r.RegFoodType != nil { food = foodMap[*r.RegFoodType] }

	return []string{
		strconv.Itoa(i + 1),
		statusMap[r.Status],
		r.Username, r.DisplayName,
		r.RegNameZh, r.RegNameEn,
		r.RegIDNumber, r.RegPassportNumber,
		gender, r.RegBirthday,
		r.RegPhone, r.RegEmail,
		r.RegShirtSize, food, r.RegAddress,
		r.RegEmergencyContact, r.RegEmergencyPhone, r.RegEmergencyRelation,
		r.Note,
		r.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func exportCSV(c *gin.Context, ev *models.Event, list []*models.EventRegistration) {
	var buf bytes.Buffer
	// BOM for Excel UTF-8
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"賽事：" + ev.Title})
	_ = w.Write(regHeaders)
	for i, r := range list {
		_ = w.Write(regToRow(i, r))
	}
	w.Flush()

	filename := fmt.Sprintf("報名清單_%s.csv", ev.Title)
	c.Header("Content-Disposition", `attachment; filename*=UTF-8''`+encodeFilename(filename))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
}

func exportExcel(c *gin.Context, ev *models.Event, list []*models.EventRegistration) {
	// Generate CSV with .xlsx extension for basic Excel compatibility
	// For full XLSX, integrate github.com/xuri/excelize/v2 when needed
	var buf bytes.Buffer
	buf.Write([]byte{0xEF, 0xBB, 0xBF})
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"賽事：" + ev.Title})
	_ = w.Write(regHeaders)
	for i, r := range list {
		_ = w.Write(regToRow(i, r))
	}
	w.Flush()

	filename := fmt.Sprintf("報名清單_%s.csv", ev.Title)
	c.Header("Content-Disposition", `attachment; filename*=UTF-8''`+encodeFilename(filename))
	c.Data(http.StatusOK, "application/vnd.ms-excel", buf.Bytes())
}

func encodeFilename(name string) string {
	var encoded []byte
	for _, b := range []byte(name) {
		if (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') ||
			b == '-' || b == '_' || b == '.' {
			encoded = append(encoded, b)
		} else {
			encoded = append(encoded, []byte(fmt.Sprintf("%%%02X", b))...)
		}
	}
	return string(encoded)
}

func parseAdminID(c *gin.Context, param string) (uint64, error) {
	return strconv.ParseUint(c.Param(param), 10, 64)
}

func mustAdminUserID(c *gin.Context) uint64 {
	raw, _ := c.Get("user_id")
	switch v := raw.(type) {
	case float64: return uint64(v)
	case string:
		id, _ := strconv.ParseUint(v, 10, 64)
		return id
	}
	return 0
}
