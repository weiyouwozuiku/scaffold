package engine

import (
	"net/http"
	"strings"
)

type root map[string]*node

type router struct {
	root
}

func newRouter() *router {
	return &router{make(root)}
}

func parsePattern(pattern string) []string {
	info := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range info {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRouter(method, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	_, ok := r.root[method]
	if !ok {
		r.root[method] = &node{}
	}
	r.root[method].insert(pattern, parts, 0, handler)
}
func (r *router) getRouter(method, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.root[method]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
		}
	}
}
func (r *router) getRouters(method, path string) (*node, map[string]string) {
	root, ok := r.root[method]
	if !ok {
		return nil, nil
	}
	searchParts := parsePattern(path)
	n := root.search(searchParts, 0)
	if n != nil {

	}

}

func (r *router) handler(c *Context) {
	n, params := r.getRouter(c.Method, c.Path)
	if n != nil {
		c.Params = params
		n.handler(c)
	} else {
		c.STRING(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
