package ginX

import (
	"fmt"
	"strings"
)

// node 路由树节点
type node struct {
	pattern  string  // 待匹配路由
	part     string  // 当前路由段
	children []*node // 子节点集合
	isWild   bool    // 该节点是否包含模糊匹配
}

func (n *node) Info() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// matchChild 返回第一个成功匹配的子节点 用于插入完整路径
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren 返回所有成功匹配的子节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// insert 每注册一条路由规则 都要调用这个方法向对应的路由树中增加一组节点
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		// 叶子节点的pattern才赋值
		// 便于在匹配结束时，依据返回节点的pattern字段是否为空，来判定是否匹配成功
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		// 没有匹配到，则新增一个子节点`树干`
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	// 如果匹配成功（包含精确与模糊匹配），则继续递归插入下一级子节点
	child.insert(pattern, parts, height+1)
}

// search 逐级匹配，直到匹配到通配符或者叶子节点
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

// travel 返回所有pattern非空的叶子节点 也就是所有已注册的路由规则
func (n *node) travel(list *[]*node) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}
