package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"trbb/internal/services"
)

type EventHandler struct {
	eventSvc *services.EventService
}

func NewEventHandler(eventSvc *services.EventService) *EventHandler {
	return &EventHandler{eventSvc: eventSvc}
}

// GET /v1/api/events
func (h *EventHandler) ListEvents(c *gin.Context) {
	var in services.ListEventsInput
	_ = c.ShouldBindQuery(&in)
	result, err := h.eventSvc.ListPublic(c.Request.Context(), in)
	if err != nil {
		_ = c.Error(err) // 讓 middleware logger 印出來
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢賽事失敗"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GET /v1/api/events/:id
func (h *EventHandler) GetEvent(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	ev, err := h.eventSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "賽事不存在"})
		} else {
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢失敗"})
		}
		return
	}
	if ev.Status != 1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "賽事不存在"})
		return
	}
	c.JSON(http.StatusOK, ev)
}

// POST /v1/api/events/:id/register
func (h *EventHandler) Register(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var in services.RegistrationInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	reg, err := h.eventSvc.Register(c.Request.Context(), id, mustUserID(c), in)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrAlreadyRegistered):
			c.JSON(http.StatusConflict, gin.H{"error": "您已報名此賽事"})
		case errors.Is(err, services.ErrRegClosed):
			c.JSON(http.StatusForbidden, gin.H{"error": "報名時間已結束"})
		case errors.Is(err, services.ErrEventFull):
			c.JSON(http.StatusForbidden, gin.H{"error": "報名人數已額滿"})
		default:
			_ = c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "報名失敗"})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "報名成功！", "registration": reg})
}

// DELETE /v1/api/events/:id/register
func (h *EventHandler) CancelRegistration(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.eventSvc.CancelRegistration(c.Request.Context(), id, mustUserID(c)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到報名記錄"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "已取消報名"})
}

// GET /v1/api/events/:id/register
func (h *EventHandler) GetMyRegistration(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	reg, err := h.eventSvc.GetMyRegistration(c.Request.Context(), id, mustUserID(c))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"registered": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"registered": true, "registration": reg})
}


