package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type H map[string]any

type Context struct {
	// 请求相关
	Req *http.Request
	// 响应相关
	Writer http.ResponseWriter

	// 请求信息
	Path   string
	Method string
	Params map[string]string

	// 响应信息
	StatusCode int

	Ctx context.Context
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Req:    req,
		Writer: w,
		Path:   req.URL.Path,
		Method: req.Method,
		Ctx:    context.Background(),
	}
}
func (c *Context) Param(key string) string {
	value, ok := c.Params[key]
	if !ok {
		slog.Log(c.Ctx, slog.LevelError, "context Param error||key:%s is not exit")
	}
	return value
}
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}
func (c *Context) JSON(code int, obj any) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		slog.Log(c.Ctx, slog.LevelError, "context JSON error")
		http.Error(c.Writer, err.Error(), 500)
	}
}
func (c *Context) DATA(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
func (c *Context) STRING(code int, format string, values ...any) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))

}
