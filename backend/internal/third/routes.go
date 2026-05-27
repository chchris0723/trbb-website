package third

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"trbb/internal/config"
	"trbb/pkg/cache"
	"trbb/pkg/database"
	"trbb/pkg/logger"
)

// RegisterRoutes mounts all /v1/third/* routes.
// These are callback/webhook endpoints from third-party providers.
func RegisterRoutes(
	r *gin.RouterGroup,
	db *database.DB,
	rdb *cache.Cache,
	cfg *config.Config,
	log *logger.Logger,
) {
	h := &handler{db: db, rdb: rdb, cfg: cfg, log: log}

	// ── Payment Callbacks ───────────────────────────────────
	payment := r.Group("/payment")
	{
		// ECPay / GreenWorld
		payment.POST("/callback", h.PaymentCallback)   // Synchronous return
		payment.POST("/notify", h.PaymentNotify)       // Asynchronous notify
		// Stripe webhook
		payment.POST("/stripe/webhook", h.StripeWebhook)
		// LINE Pay confirm
		payment.GET("/linepay/confirm", h.LinePayConfirm)
		payment.GET("/linepay/cancel", h.LinePayCancel)
	}

	// ── Garmin Connect ──────────────────────────────────────
	garmin := r.Group("/garmin")
	{
		garmin.GET("/auth", h.GarminAuthRedirect)       // Step 1: redirect user
		garmin.GET("/callback", h.GarminCallback)       // Step 2: OAuth callback
		garmin.POST("/deregister", h.GarminDeregister)  // User deregistration push
		garmin.POST("/activities", h.GarminActivitiesPush) // Activity push (webhook)
		garmin.POST("/user-metrics", h.GarminUserMetrics)
		garmin.DELETE("/user", h.GarminUserDelete)
	}

	// ── Google OAuth ────────────────────────────────────────
	r.GET("/google/callback", h.GoogleCallback)

	// ── Facebook OAuth ──────────────────────────────────────
	r.GET("/facebook/callback", h.FacebookCallback)

	// ── LINE OAuth / Notify ─────────────────────────────────
	r.GET("/line/callback", h.LineCallback)

	// ── SMS / Email Webhooks ────────────────────────────────
	r.POST("/sms/status", h.SMSStatusCallback)
	r.POST("/email/event", h.EmailEventCallback)
}

type handler struct {
	db  *database.DB
	rdb *cache.Cache
	cfg *config.Config
	log *logger.Logger
}

// Stub handlers – implement actual provider logic per file
func (h *handler) PaymentCallback(c *gin.Context)      { c.JSON(http.StatusNotImplemented, gin.H{"msg": "TODO"}) }
func (h *handler) PaymentNotify(c *gin.Context)        { c.JSON(http.StatusNotImplemented, gin.H{"msg": "TODO"}) }
func (h *handler) StripeWebhook(c *gin.Context)        { c.JSON(http.StatusNotImplemented, gin.H{"msg": "TODO"}) }
func (h *handler) LinePayConfirm(c *gin.Context)       { c.JSON(http.StatusNotImplemented, gin.H{"msg": "TODO"}) }
func (h *handler) LinePayCancel(c *gin.Context)        { c.JSON(http.StatusNotImplemented, gin.H{"msg": "TODO"}) }
func (h *handler) GarminAuthRedirect(c *gin.Context)   { c.JSON(http.StatusNotImplemented, gin.H{"msg": "TODO"}) }
func (h *handler) GarminCallback(c *gin.Context)       { c.JSON(http.StatusNotImplemented, gin.H{"msg": "TODO"}) }
func (h *handler) GarminDeregister(c *gin.Context)     { c.JSON(http.StatusNotImplemented, gin.H{"msg": "TODO"}) }
func (h *handler) GarminActivitiesPush(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"msg": "TODO"}) }
func (h *handler) GarminUserMetrics(c *gin.Context)    { c.JSON(http.StatusNotImplemented, gin.H{"msg": "TODO"}) }
func (h *handler) GarminUserDelete(c *gin.Context)     { c.JSON(http.StatusNotImplemented, gin.H{"msg": "TODO"}) }
func (h *handler) GoogleCallback(c *gin.Context)       { c.JSON(http.StatusNotImplemented, gin.H{"msg": "TODO"}) }
func (h *handler) FacebookCallback(c *gin.Context)     { c.JSON(http.StatusNotImplemented, gin.H{"msg": "TODO"}) }
func (h *handler) LineCallback(c *gin.Context)         { c.JSON(http.StatusNotImplemented, gin.H{"msg": "TODO"}) }
func (h *handler) SMSStatusCallback(c *gin.Context)    { c.JSON(http.StatusNotImplemented, gin.H{"msg": "TODO"}) }
func (h *handler) EmailEventCallback(c *gin.Context)   { c.JSON(http.StatusNotImplemented, gin.H{"msg": "TODO"}) }
