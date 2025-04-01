package maxconns

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/thinkgos/limiter/limit"
)

// MaxConns returns a middleware that limit the concurrent connections.
func MaxConns(n int) gin.HandlerFunc {
	if n <= 0 {
		return func(c *gin.Context) {
			c.Next()
		}
	}
	latch := limit.NewLimit(n)
	return func(c *gin.Context) {
		if latch.TryBorrow() {
			defer func() {
				if err := latch.Return(); err != nil {
					fmt.Fprintf(os.Stderr, "maxconns: return conns failure, %v\n", err)
				}
			}()
			c.Next()
		} else {
			c.AbortWithStatus(http.StatusServiceUnavailable)
		}
	}
}
