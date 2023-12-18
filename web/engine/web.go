package engine

import "net/http"

type HandlerFunc func(ctx *Context)

type Engine struct {
	*router
}

func New() *Engine {
	return &Engine{newRouter()}
}
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}
func (e *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c := newContext(writer, request)
	e.router.handler(c)
}
