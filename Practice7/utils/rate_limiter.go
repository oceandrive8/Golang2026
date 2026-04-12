package utils

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type clientData struct {
	count     int
	firstSeen time.Time
}

var (
	mu      sync.Mutex
	clients = make(map[string]*clientData)
)

const (
	limit      = 5
	timeWindow = 1 * time.Minute
)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {

		var key string

		// 1. Identify user (JWT or IP)
		if userID, exists := c.Get("userID"); exists {
			key = "user_" + userID.(string)
		} else {
			key = "ip_" + c.ClientIP()
		}

		// 2. Lock for thread safety
		mu.Lock()
		defer mu.Unlock()
		now := time.Now()
		data, exists := clients[key]

		if !exists {
			clients[key] = &clientData{
				count:     1,
				firstSeen: now,
			}
			c.Next()
			return
		}

		// reset window
		if now.Sub(data.firstSeen) > timeWindow {
			data.count = 1
			data.firstSeen = now
			c.Next()
			return
		}

		data.count++

		// block if exceeded
		if data.count > limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			return
		}

		c.Next()
	}
}
