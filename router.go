package ginX

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node // 请求方式及其对应的路由前缀树
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// parsePattern 解析路径 生成路由段集合
func parsePattern(pattern string) []string {
	ps := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, part := range ps {
		if part != "" {
			parts = append(parts, part)
			// 通配符 * 必须位于最后一个路由段，且只能出现一次
			if part[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	// 每种请求方式对应一棵路由前缀树
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path) // 路由段集合
	params := make(map[string]string) // 路径参数集合

	// root就是method请求方法对应的路由树
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for idx, part := range parts {
			if part[0] == ':' {
				// 记录路径参数
				params[part[1:]] = searchParts[idx]
			}
			// 记录通配符匹配结果
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[idx:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// getRouters 返回某请求方式下所有已注册的路由
func (r *router) getRouters(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Request.Method, c.Request.URL.Path)
	if n != nil {
		c.Params = params
		// 调用匹配到的handler
		key := c.Request.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 page not found: %s\n", c.Request.URL.Path)
	}
}
