package engine

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type H map[string]any

type traceKey struct{}

var TraceKey traceKey

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

	TimeStamp time.Time
	Ctx       context.Context
}

// generateTraceID 生成一个 OpenTelemetry 规范的 Trace ID
func generateTraceID() string {
	// 生成 16 字节的随机数
	traceIDBytes := make([]byte, 16)
	rand.Read(traceIDBytes)

	// 将随机数转换为十六进制字符串
	traceID := hex.EncodeToString(traceIDBytes)
	return traceID
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	trace := ""
	if req.Header.Get("trace") != "" {
		trace = req.Header.Get("trace")
	} else {
		trace = generateTraceID()
	}
	return &Context{
		Req:       req,
		Writer:    w,
		Path:      req.URL.Path,
		Method:    req.Method,
		TimeStamp: time.Now(),
		Ctx:       context.WithValue(context.Background(), TraceKey, trace),
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
