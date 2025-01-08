package ginX

import (
	"html/template"
	"net/http"
	"strings"
)

// HandlerFunc 路由对应实际执行的handler
type HandlerFunc func(*Context)

// Engine 接管路由请求的实例
type Engine struct {
	// 路由映射
	router *router

	// 分组控制
	*RouterGroup
	groups []*RouterGroup // 存储所有路由分组

	// 模板渲染
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// SetFuncMap 自定义处理函数 可在模板中调用
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// LoadHTMLGlob 加载模板集合
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// addRoute 添加路由规则 辅助函数，非导出
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
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		// 依据请求路径判断需要执行哪些分组中间件
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	// 构造上下文
	// 每收到一个并发请求，都会新建一个上下文，不需要关心参数覆盖的问题
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}
