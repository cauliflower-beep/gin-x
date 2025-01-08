package ginX

import (
	"fmt"
	"html/template"
	"net/http"
	"testing"
	"time"
)

func print4V1() HandlerFunc {
	return func(c *Context) {
		fmt.Println("here is v1's middleware.")
	}
}

func TestGinX(t *testing.T) {

	// 自定义模板渲染函数，可在模板中调用
	FormatAsDate := func(t time.Time) string {
		year, month, day := t.Date()
		return fmt.Sprintf("%d-%02d-%02d", year, month, day)
	}

	r := New()
	// 全局中间件
	r.Use(Logger(), Recovery())

	// 注册模板渲染函数
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./static")

	r.GET("/index", func(c *Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>", nil)
	})

	r.GET("/panic", func(c *Context) {
		arr := []int{1, 2, 3}
		c.String(http.StatusOK, "%#v", arr[3])
	})

	r.GET("/hello", func(c *Context) {
		c.HTML(http.StatusOK, "demo.tmpl", H{
			"title": "gin-x",
			"now":   time.Now(),
		})
	})

	v1 := r.Group("/v1")
	v1.Use(print4V1())
	{
		v1.GET("/", func(c *Context) {
			c.HTML(http.StatusOK, "<h1>Hello GinX</h1>", nil)
		})

		v1.GET("/hello", func(c *Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Request.URL.Path)
		})
	}
	v2 := r.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Request.URL.Path)
		})
		v2.POST("/login", func(c *Context) {
			time.Sleep(time.Second)
			c.JSON(http.StatusOK, H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	_ = r.Run("0.0.0.0:9999")
}
