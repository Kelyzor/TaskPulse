package middleware

import (
	"strconv"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
)

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	seconds := int(time.Until(info.ResetTime).Seconds())
	c.String(429, "Too many requests. Try again in "+strconv.Itoa(seconds)+"s")
}

func RateLimitMiddleware(requestsPerMinute int) gin.HandlerFunc {
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Minute,
		Limit: uint(requestsPerMinute),
	})
	return ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})
}
