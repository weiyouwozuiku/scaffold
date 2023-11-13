package engine

type HandlerFunc func(ctx *Context)

type Engine struct {
	router
}
