package ginX

import (
	"net/http"
)

// 路由对应实际执行的handler
type HandlerFunc func(*Context)

// 接管路由请求的实例
type Engine struct {
	// 路由映射
	router *router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

// 添加路由规则 内部使用，非导出
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 构造上下文
	c := newContext(w, req)
	engine.router.handle(c)
}
