package ginX

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H map[string]any 的别名 完全仿照 gin 进行的简化
type H map[string]any

// Context 封装context 包含整个请求->响应的数据集 也便于定义方法简化部分响应重复代码
type Context struct {
	// 原始对象
	Writer  http.ResponseWriter
	Request *http.Request

	// 请求相关
	Params map[string]string // 解析出来的路径参数
	// 响应相关
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: req,
	}
}

// Param 访问请求路径中的参数
func (c *Context) Param(key string) string {
	val, _ := c.Params[key]
	return val
}

func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader 设置响应报文头部信息
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// String 字符串响应的便捷构造方式
func (c *Context) String(code int, format string, values ...any) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	_, _ = c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON json响应的便捷构造方式
func (c *Context) JSON(code int, obj any) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	// 构造一个 Encoder 将输出写入到 c.Writer
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	_, _ = c.Writer.Write(data)
}

// HTML html响应的便捷构造方式
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	_, _ = c.Writer.Write([]byte(html))
}
