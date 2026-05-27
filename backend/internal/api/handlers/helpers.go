package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// mustUserID extracts user ID from JWT claims set by middleware.
func mustUserID(c *gin.Context) uint64 {
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

// parseID parses a uint64 path param.
func parseID(c *gin.Context, param string) (uint64, error) {
	return strconv.ParseUint(c.Param(param), 10, 64)
}
