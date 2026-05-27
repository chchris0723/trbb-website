package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"trbb/internal/models"
	"trbb/internal/services"
)

type AdminUserHandler struct {
	userSvc *services.UserService
}

func NewAdminUserHandler(userSvc *services.UserService) *AdminUserHandler {
	return &AdminUserHandler{userSvc: userSvc}
}

// POST /v1/admin/auth/login
func (h *AdminUserHandler) Login(c *gin.Context) {
	var in services.LoginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.userSvc.Login(c.Request.Context(), in)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidPassword):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "帳號或密碼錯誤"})
		default:
			c.JSON(http.StatusForbidden, gin.H{"error": "帳號狀態異常或登入失敗"})
		}
		return
	}
	if !result.User.IsAdmin() {
		c.JSON(http.StatusForbidden, gin.H{"error": "權限不足"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token":         result.Token,
		"refresh_token": result.Refresh,
		"user":          result.User.PublicProfile(),
	})
}

// ── 會員列表 (role=1) ─────────────────────────────────────────

// GET /v1/admin/members
func (h *AdminUserHandler) ListMembers(c *gin.Context) {
	var in services.ListUsersInput
	_ = c.ShouldBindQuery(&in)
	role := models.RoleMember
	in.Role = &role
	h.listUsers(c, in)
}

// ── 管理員列表 (role>=8) ──────────────────────────────────────

// GET /v1/admin/admins
func (h *AdminUserHandler) ListAdmins(c *gin.Context) {
	var in services.ListUsersInput
	_ = c.ShouldBindQuery(&in)
	// 取得 role 8 + 9
	in.Role = nil // service 裡加 minRole 過濾
	result, err := h.userSvc.ListAdmins(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢失敗"})
		return
	}
	list := make([]map[string]any, 0, len(result.Users))
	for _, u := range result.Users {
		list = append(list, u.FullProfile())
	}
	c.JSON(http.StatusOK, gin.H{
		"users": list, "total": result.Total,
		"page": result.Page, "pages": result.Pages,
	})
}

// ── 共用 ──────────────────────────────────────────────────────

func (h *AdminUserHandler) listUsers(c *gin.Context, in services.ListUsersInput) {
	result, err := h.userSvc.ListUsers(c.Request.Context(), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢失敗"})
		return
	}
	list := make([]map[string]any, 0, len(result.Users))
	for _, u := range result.Users {
		list = append(list, u.FullProfile())
	}
	c.JSON(http.StatusOK, gin.H{
		"users": list, "total": result.Total,
		"page": result.Page, "pages": result.Pages,
	})
}

// GET /v1/admin/users/:id
func (h *AdminUserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	user, err := h.userSvc.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "使用者不存在"})
		return
	}
	c.JSON(http.StatusOK, user.FullProfile())
}

// PUT /v1/admin/users/:id/profile — 修改資料（不含帳號）
func (h *AdminUserHandler) UpdateUserProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var in services.AdminUpdateProfileInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 取得被修改者資料以確認 role
	target, err := h.userSvc.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "使用者不存在"})
		return
	}
	if err := h.userSvc.AdminUpdateProfile(
		c.Request.Context(), id, in,
		mustOperatorRole(c), target.Role,
	); err != nil {
		if errors.Is(err, services.ErrPermissionDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": "權限不足"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失敗"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "資料已更新"})
}

// PUT /v1/admin/users/:id/password — 直接設定密碼
func (h *AdminUserHandler) SetUserPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		Password string `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	target, err := h.userSvc.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "使用者不存在"})
		return
	}
	if err := h.userSvc.AdminSetPassword(
		c.Request.Context(), id, body.Password,
		mustOperatorRole(c), target.Role,
	); err != nil {
		if errors.Is(err, services.ErrPermissionDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": "權限不足"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "修改失敗"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "密碼已更新"})
}

// PUT /v1/admin/users/:id/status
func (h *AdminUserHandler) UpdateUserStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		Status int `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	operatorRole := mustOperatorRole(c)
	if err := h.userSvc.UpdateStatus(c.Request.Context(), id, body.Status, operatorRole); err != nil {
		if errors.Is(err, services.ErrPermissionDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": "權限不足"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失敗"})
		return
	}
	msg := map[int]string{
		models.StatusActive:    "已核准",
		models.StatusRejected:  "已拒絕",
		models.StatusSuspended: "已停用",
		models.StatusPending:   "已設回待審",
	}[body.Status]
	if msg == "" {
		msg = "已更新"
	}
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

// POST /v1/admin/members — 後台新增一般會員（管理員以上可用）
func (h *AdminUserHandler) CreateMember(c *gin.Context) {
	var in services.CreateMemberInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.userSvc.CreateMember(c.Request.Context(), in, mustOperatorRole(c))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmailExists):
			c.JSON(http.StatusConflict, gin.H{"error": "此 Email 已存在"})
		case errors.Is(err, services.ErrUsernameExists):
			c.JSON(http.StatusConflict, gin.H{"error": "此會員 ID 已存在"})
		case errors.Is(err, services.ErrPermissionDenied):
			c.JSON(http.StatusForbidden, gin.H{"error": "權限不足"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "建立失敗"})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "會員建立成功", "user": user.PublicProfile()})
}

// POST /v1/admin/admins — 超級管理員新增管理員
func (h *AdminUserHandler) CreateAdmin(c *gin.Context) {
	var in services.CreateAdminInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.userSvc.CreateAdmin(c.Request.Context(), in, mustOperatorRole(c))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmailExists):
			c.JSON(http.StatusConflict, gin.H{"error": "此 Email 已存在"})
		case errors.Is(err, services.ErrPermissionDenied):
			c.JSON(http.StatusForbidden, gin.H{"error": "僅超級管理員可新增管理員"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "建立失敗"})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "管理員建立成功", "user": user.PublicProfile()})
}

// DELETE /v1/admin/admins/:id — 刪除管理員（超級管理員限定，不能刪超級）
func (h *AdminUserHandler) DeleteAdmin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	target, err := h.userSvc.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "使用者不存在"})
		return
	}
	if target.Role >= models.RoleSuper {
		c.JSON(http.StatusForbidden, gin.H{"error": "超級管理員不可刪除"})
		return
	}
	if err := h.userSvc.SoftDelete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "刪除失敗"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "管理員已刪除"})
}

func mustOperatorRole(c *gin.Context) int {
	raw, _ := c.Get("user_role")
	if v, ok := raw.(float64); ok {
		return int(v)
	}
	return 0
}

func mustOperatorID(c *gin.Context) uint64 {
	raw, _ := c.Get("user_id")
	switch v := raw.(type) {
	case float64:
		return uint64(v)
	case string:
		id, _ := strconv.ParseUint(v, 10, 64)
		return id
	}
	return 0
}
