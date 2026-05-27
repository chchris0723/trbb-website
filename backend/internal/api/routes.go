package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apiHandlers "trbb/internal/api/handlers"
	"trbb/internal/config"
	"trbb/internal/middleware"
	"trbb/internal/services"
	"trbb/pkg/cache"
	"trbb/pkg/database"
	"trbb/pkg/logger"
	"trbb/pkg/storage"
)

func RegisterRoutes(r *gin.RouterGroup, db *database.DB, rdb *cache.Cache,
	minio *storage.Storage, cfg *config.Config, log *logger.Logger) {

	userSvc  := services.NewUserService(db, cfg.App.SecretKey)
	eventSvc := services.NewEventService(db)
	authH    := apiHandlers.NewAuthHandler(userSvc)
	eventH   := apiHandlers.NewEventHandler(eventSvc)

	// ── Auth (public) ────────────────────────────────────────
	auth := r.Group("/auth")
	{
		auth.POST("/register",        authH.Register)
		auth.POST("/login",           authH.Login)
		auth.POST("/refresh",         todo("refresh"))
		auth.POST("/forgot-password", todo("forgot-password"))
		auth.POST("/reset-password",  todo("reset-password"))
	}

	// ── Public GET（任何人可讀）──────────────────────────────
	// 原則：同一路徑只能有一個 GET 登記
	r.GET("/events",             eventH.ListEvents)
	r.GET("/events/:id",         eventH.GetEvent)
	r.GET("/announcements",      todo("announcements"))
	r.GET("/announcements/:id",  todo("announcement-detail"))
	r.GET("/products",           todo("products-list"))
	r.GET("/products/:id",       todo("products-detail"))
	r.GET("/secondhand",         todo("secondhand-list"))
	r.GET("/secondhand/:id",     todo("secondhand-detail"))
	r.GET("/training",           todo("training-list-public"))
	r.GET("/training/:id",       todo("training-detail-public"))

	// ── Protected（需登入）──────────────────────────────────
	// 僅掛 POST / PUT / DELETE，不重複掛已有的 GET 路徑
	p := r.Group("", middleware.JWT(cfg.App.SecretKey))
	{
		// 個人資料
		me := p.Group("/me")
		{
			me.GET("",                        authH.GetProfile)
			me.PUT("",                        authH.UpdateProfile)
			me.PUT("/password",               authH.ChangePassword)
			me.GET("/registration-profile",   authH.GetRegistrationProfile)
			me.POST("/avatar",                todo("upload-avatar"))
			me.GET("/notifications",          todo("notifications"))
			me.PUT("/notifications/:id/read", todo("mark-read"))
		}

		// 賽事（只有需要登入的動作）
		p.GET("/events/:id/register",             eventH.GetMyRegistration)
		p.POST("/events/:id/register",            eventH.Register)
		p.DELETE("/events/:id/register",          eventH.CancelRegistration)
		p.POST("/events/:id/group",               todo("group-create"))
		p.POST("/events/:id/group/:gid/join",     todo("group-join"))
		p.GET("/events/:id/transport",            todo("transport-list"))
		p.POST("/events/:id/transport/:tid/book", todo("transport-book"))

		// 訂單（專屬於登入用戶）
		p.GET("/orders",          todo("orders-list"))
		p.GET("/orders/:id",      todo("orders-detail"))
		p.POST("/orders",         todo("orders-create"))
		p.POST("/orders/:id/pay", todo("orders-pay"))

		// 二手（僅寫入需要登入，讀取已在 public GET 處理）
		p.POST("/secondhand",       todo("secondhand-create"))
		p.PUT("/secondhand/:id",    todo("secondhand-update"))
		p.DELETE("/secondhand/:id", todo("secondhand-delete"))

		// 訓練日記（登入後才能寫）
		p.POST("/training",         todo("training-create"))
		p.PUT("/training/:id",      todo("training-update"))
		p.DELETE("/training/:id",   todo("training-delete"))
		p.POST("/training/garmin/sync", todo("garmin-sync"))

		// 上傳
		p.POST("/upload/image", todo("upload-image"))
	}
}

func todo(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"msg": "TODO: " + name})
	}
}
