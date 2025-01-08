package ginX

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		now := time.Now()
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Request.RequestURI, time.Since(now))
	}
}
