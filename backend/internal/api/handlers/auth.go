package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"trbb/internal/services"
)

type AuthHandler struct {
	userSvc *services.UserService
}

func NewAuthHandler(userSvc *services.UserService) *AuthHandler {
	return &AuthHandler{userSvc: userSvc}
}

// POST /v1/api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var in services.RegisterInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.userSvc.Register(c.Request.Context(), in)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmailExists):
			c.JSON(http.StatusConflict, gin.H{"error": "此 Email 已被註冊"})
		case errors.Is(err, services.ErrUsernameExists):
			c.JSON(http.StatusConflict, gin.H{"error": "此會員 ID 已被使用"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "註冊失敗，請稍後再試"})
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "註冊成功，請等待管理員審核",
		"user":    user.PublicProfile(),
	})
}

// POST /v1/api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
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
		case errors.Is(err, services.ErrAccountPending):
			c.JSON(http.StatusForbidden, gin.H{"error": "帳號尚待管理員審核，請耐心等候"})
		case errors.Is(err, services.ErrAccountSuspended):
			c.JSON(http.StatusForbidden, gin.H{"error": "帳號已被停用，請聯繫管理員"})
		case errors.Is(err, services.ErrAccountRejected):
			c.JSON(http.StatusForbidden, gin.H{"error": "帳號申請已被拒絕，請聯繫管理員"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "登入失敗，請稍後再試"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token":         result.Token,
		"refresh_token": result.Refresh,
		"user":          result.User.FullProfile(),
	})
}

// GET /v1/api/me
func (h *AuthHandler) GetProfile(c *gin.Context) {
	user, err := h.userSvc.FindByID(c.Request.Context(), mustUserID(c))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "使用者不存在"})
		return
	}
	c.JSON(http.StatusOK, user.FullProfile())
}

// PUT /v1/api/me
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var in services.UpdateProfileInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.userSvc.UpdateProfile(c.Request.Context(), mustUserID(c), in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失敗，請稍後再試"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "個人資料已更新",
		"user":    user.FullProfile(),
	})
}

// PUT /v1/api/me/password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var in struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.userSvc.ChangePassword(c.Request.Context(), mustUserID(c), in.OldPassword, in.NewPassword); err != nil {
		if errors.Is(err, services.ErrInvalidPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "舊密碼錯誤"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "修改失敗"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "密碼修改成功"})
}

// GET /v1/api/me/registration-profile
// 給賽事報名頁帶入的個人資料
func (h *AuthHandler) GetRegistrationProfile(c *gin.Context) {
	user, err := h.userSvc.FindByID(c.Request.Context(), mustUserID(c))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "使用者不存在"})
		return
	}
	c.JSON(http.StatusOK, user.RegistrationProfile())
}

