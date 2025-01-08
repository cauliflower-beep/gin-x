package ginX

import "log"

type RouterGroup struct {
	prefix      string        // 分组前缀 gin也是用共同前缀区分分组的
	middlewares []HandlerFunc // 作用在分组上的中间件
	parent      *RouterGroup
	engine      *Engine // 通过engine示例访问 Router 所有路由分组共享一个全局engine实例
}

// Group 基于当前路由分组创建一个新子分组
// 共同前缀由框架用户提供
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("add route: %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// Use 将中间件应用到路由分组上 跟gin的做法相同
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}
