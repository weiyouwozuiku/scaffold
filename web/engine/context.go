package engine

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net"
	"net/http"
	"os"
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

const (
	HEADER_RID = "Header-Rid"
)

func GetLocalIP() string {
	ip := "127.0.0.1"
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ip
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	}

	return ip
}

// generateTraceID 生成一个 OpenTelemetry 规范的 Trace ID
func genTraceID() string {

	ip := GetLocalIP()

	now := time.Now()
	timestamp := uint32(now.Unix())
	timeNano := now.UnixNano()
	pid := os.Getpid()
	b := bytes.Buffer{}

	b.WriteString(hex.EncodeToString(net.ParseIP(ip).To4()))
	b.WriteString(fmt.Sprintf("%x", timestamp&0xffffffff))
	b.WriteString(fmt.Sprintf("%04x", timeNano&0xffff))
	b.WriteString(fmt.Sprintf("%04x", pid&0xffff))
	b.WriteString(fmt.Sprintf("%06x", rand.Int31n(1<<24)))
	b.WriteString("b0")

	return b.String()
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	trace := ""
	if req.Header.Get("trace") != "" {
		trace = req.Header.Get("trace")
	} else {
		trace = genTraceID()
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
