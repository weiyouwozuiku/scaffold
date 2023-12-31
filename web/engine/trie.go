package engine

import (
	"fmt"
	"strings"
)

type node struct {
	pattern  string      // 待匹配路由，例如 /p/:lang
	part     string      // 路由中的一部分，例如 :lang
	children []*node     // 子节点，例如 [doc, tutorial, intro]
	isWild   bool        // 是否精确匹配，part 含有 : 或 * 时为true
	handler  HandlerFunc // 处理方法
}

func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s, part=%s, isWild=%t}", n.pattern, n.part, n.isWild)
}

// search for insert
func (n *node) matchChild(part string) *node {
	for _, it := range n.children {
		if it.part == part || it.isWild {
			return it
		}
	}
	return nil
}
func (n *node) insert(pattern string, parts []string, height int, handler HandlerFunc) {
	// 退出递归条件
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	// url中的/分割组中的一个
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1, handler)
}
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
}
func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}
