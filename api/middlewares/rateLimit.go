package middlewares

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/models/response"
)

func RateLimitMiddleware(r rate.Limit, b int) gin.HandlerFunc {

	limiter := rate.NewLimiter(r, b)
	return func(c *gin.Context) {
		res := response.Gin{Res: c}
		if !limiter.Allow() {
			res.Response(c, statecode.TokenErr, nil)
			c.Abort()
			return
		}
	}
}
