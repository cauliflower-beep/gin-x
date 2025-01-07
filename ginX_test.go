package ginX

import (
	"net/http"
	"testing"
)

func TestGinX(t *testing.T) {
	r := New()
	r.GET("/", func(c *Context) {
		c.HTML(http.StatusOK, "<h1>Hello GinX</h1>")
	})

	r.GET("/hello", func(c *Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Request.URL.Path)
	})

	r.GET("/hello/:name", func(c *Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Request.URL.Path)
	})

	r.GET("/assets/*filepath", func(c *Context) {
		c.JSON(http.StatusOK, H{"filepath": c.Param("filepath")})
	})

	_ = r.Run("0.0.0.0:9999")
}
