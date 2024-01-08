package easierweb

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/websocket"
	"gopkg.in/yaml.v3"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Context struct {
	Route          string
	Header         Params
	Path           Params
	Query          Params
	Form           Params
	Body           Data
	Code           int
	Result         Data
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	WebsocketConn  *websocket.Conn
	Flusher        http.Flusher
	Logger         *slog.Logger
	index          int
	handles        []Handle
	written        bool
	closed         bool
}

func (c *Context) Next() {
	c.index++
	for c.index < len(c.handles) {
		c.handles[c.index](c)
		c.index++
	}
}

func (c *Context) Abort() {
	c.index = len(c.handles) + 1
}

// POST Form File

func (c *Context) FileKeys() []string {
	files := c.Request.MultipartForm.File
	var ks = make([]string, 0, len(files))
	for k := range files {
		ks = append(ks, k)
	}
	return ks
}

func (c *Context) GetFile(key string) (multipart.File, error) {
	file, _, err := c.Request.FormFile(key)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Query/Form/Path/Header Params Bind

func (c *Context) BindQuery(obj any) error {
	return c.Query.Bind(obj)
}

func (c *Context) BindForm(obj any) error {
	return c.Form.Bind(obj)
}

func (c *Context) BindPath(obj any) error {
	return c.Path.Bind(obj)
}

func (c *Context) BindHeader(obj any) error {
	return c.Header.Bind(obj)
}

// POST Body Bind

func (c *Context) BindJSON(obj any) error {
	return c.Body.ParseJSON(obj)
}

func (c *Context) BindYAML(obj any) error {
	return c.Body.ParseYAML(obj)
}

func (c *Context) BindXML(obj any) error {
	return c.Body.ParseXML(obj)
}

// Result Write

func (c *Context) WriteJSON(code int, obj any) {
	if c.written {
		return
	}
	marshal, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	c.Write(code, marshal)
}

func (c *Context) WriteYAML(code int, obj any) {
	if c.written {
		return
	}
	marshal, err := yaml.Marshal(obj)
	if err != nil {
		panic(err)
	}
	c.Write(code, marshal)
}

func (c *Context) WriteXML(code int, obj any) {
	if c.written {
		return
	}
	marshal, err := xml.Marshal(obj)
	if err != nil {
		panic(err)
	}
	c.Write(code, marshal)
}

func (c *Context) Redirect(code int, url string) {
	if c.written {
		return
	}
	http.Redirect(c.ResponseWriter, c.Request, url, code)
}

func (c *Context) WriteLocalFile(contentType, fileName, localFilePath string) {
	if c.written {
		return
	}
	fileBytes, err := os.ReadFile(localFilePath)
	if err != nil {
		panic(err)
	}
	c.WriteFile(contentType, fileName, fileBytes)
}

func (c *Context) WriteFile(contentType, fileName string, fileBytes []byte) {
	if c.written {
		return
	}
	if len(fileName) > 0 {
		c.SetContentDisposition(fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	} else {
		c.SetContentDisposition(fmt.Sprintf("attachment; filename=\"%v\"", time.Now().Unix()))
	}
	if len(contentType) == 0 {
		c.SetContentType("application/octet-stream")
	} else {
		c.SetContentType(contentType)
	}
	c.Write(http.StatusOK, fileBytes)
}

func (c *Context) WriteString(code int, text string) {
	if c.written {
		return
	}
	c.Write(code, []byte(text))
}

func (c *Context) NoContent(code int) {
	c.Write(code, nil)
}

func (c *Context) Write(code int, data []byte) {
	if c.written {
		return
	}
	c.ResponseWriter.WriteHeader(code)
	_, err := c.ResponseWriter.Write(data)
	if err != nil {
		panic(err)
	}
	c.Code = code
	c.Result = data
	c.written = true
}

func (c *Context) SetContentType(value string) {
	c.SetHeader("Content-Type", value)
}

func (c *Context) SetContentDisposition(value string) {
	c.SetHeader("Content-Disposition", value)
}

func (c *Context) SetHeader(key, value string) {
	c.ResponseWriter.Header().Set(key, value)
}

func (c *Context) AddHeader(key, value string) {
	c.ResponseWriter.Header().Add(key, value)
}

// WS Receive

func (c *Context) ReceiveJSON(obj any) error {
	return websocket.JSON.Receive(c.WebsocketConn, obj)
}

func (c *Context) ReceiveYAML(obj any) error {
	var buf string
	err := websocket.Message.Receive(c.WebsocketConn, &buf)
	if err != nil {
		return err
	}
	return yaml.Unmarshal([]byte(buf), obj)
}

func (c *Context) ReceiveXML(obj any) error {
	var buf string
	err := websocket.Message.Receive(c.WebsocketConn, &buf)
	if err != nil {
		return err
	}
	return xml.Unmarshal([]byte(buf), obj)
}

func (c *Context) ReceiveString() (string, error) {
	var buf string
	err := websocket.Message.Receive(c.WebsocketConn, &buf)
	if err != nil {
		return "", err
	}
	return buf, nil
}

func (c *Context) Receive() ([]byte, error) {
	var buf []byte
	err := websocket.Message.Receive(c.WebsocketConn, &buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// WS Send

func (c *Context) SendJSON(obj any) error {
	marshal, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return c.Send(marshal)
}

func (c *Context) SendYAML(obj any) error {
	marshal, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	return c.Send(marshal)
}

func (c *Context) SendXML(obj any) error {
	marshal, err := xml.Marshal(obj)
	if err != nil {
		return err
	}
	return c.Send(marshal)
}

func (c *Context) SendString(text string) error {
	return c.Send([]byte(text))
}

func (c *Context) Send(msg []byte) error {
	_, err := c.WebsocketConn.Write(msg)
	if err != nil {
		return err
	}
	return nil
}

// WS Close

func (c *Context) Close() error {
	if c.closed {
		return nil
	}
	err := c.WebsocketConn.Close()
	if err != nil {
		return err
	}
	c.closed = true
	return nil
}

// SSE Push

func (c *Context) PushJSON(obj any, split string) error {
	marshal, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return c.Push(string(marshal), split)
}

func (c *Context) PushYAML(obj any, split string) error {
	marshal, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	return c.Push(string(marshal), split)
}

func (c *Context) PushXML(obj any, split string) error {
	marshal, err := xml.Marshal(obj)
	if err != nil {
		return err
	}
	return c.Push(string(marshal), split)
}

func (c *Context) Push(msg string, split string) error {
	_, err := fmt.Fprintf(c.ResponseWriter, fmt.Sprintf("%s%s", msg, split))
	if err != nil {
		return err
	}
	c.Flusher.Flush()
	return nil
}

// Other request parameters

func (c *Context) GetCookie(name string) (*http.Cookie, error) {
	return c.Request.Cookie(name)
}

func (c *Context) Cookies() []*http.Cookie {
	return c.Request.Cookies()
}

func (c *Context) URI() string {
	return c.Request.RequestURI
}

func (c *Context) Method() string {
	return c.Request.Method
}

func (c *Context) URL() *url.URL {
	return c.Request.URL
}

func (c *Context) RemoteAddr() string {
	return c.Request.RemoteAddr
}

func (c *Context) Host() string {
	return c.Request.Host
}

func (c *Context) Proto() string {
	return c.Request.Proto
}

// Set

func setContext(ctx *Context, router *Router, route string, res http.ResponseWriter, req *http.Request, par httprouter.Params, ws *websocket.Conn, middlewares ...Handle) error {

	defer func() {
		err := recover()
		if err != nil {
			router.logger.Error(fmt.Sprintf("set context error: %s", err))
		}
	}()

	handles := append([]Handle(nil), router.middlewares...)
	handles = append(handles, middlewares...)
	ctx.Route = route
	ctx.index = 0
	ctx.handles = handles
	ctx.Header = nil
	ctx.Path = nil
	ctx.Query = nil
	ctx.Form = nil
	ctx.Body = nil
	ctx.Request = req
	ctx.ResponseWriter = res
	ctx.WebsocketConn = ws
	ctx.Flusher = nil
	ctx.Logger = router.logger
	ctx.Code = 0
	ctx.Result = nil
	ctx.written = false
	ctx.closed = false

	if strings.Contains(strings.ToLower(req.Header.Get("Content-Type")), "multipart/form-data") ||
		strings.Contains(strings.ToLower(req.Header.Get("content-type")), "multipart/form-data") {
		err := req.ParseMultipartForm(router.multipartFormMaxMemory)
		if err != nil {
			return err
		}
	} else if strings.Contains(strings.ToLower(req.Header.Get("Content-Type")), "application/x-www-form-urlencoded") ||
		strings.Contains(strings.ToLower(req.Header.Get("content-type")), "application/x-www-form-urlencoded") {
		err := req.ParseForm()
		if err != nil {
			return err
		}
	} else {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}
		ctx.Body = bodyBytes
	}

	if len(req.Header) > 0 {
		ctx.Header = make(map[string]string, len(req.Header))
		for k, v := range req.Header {
			if len(v) > 0 {
				ctx.Header[k] = v[0]
			}
		}
	}

	if len(par) > 0 {
		ctx.Path = make(map[string]string, len(par))
		for _, v := range par {
			ctx.Path[v.Key] = v.Value
		}
	}

	if len(req.URL.Query()) > 0 {
		ctx.Query = make(map[string]string, len(req.URL.Query()))
		for k, v := range req.URL.Query() {
			if len(v) > 0 {
				ctx.Query[k] = v[0]
			}
		}
	}

	if len(req.PostForm) > 0 {
		ctx.Form = make(map[string]string, len(req.PostForm))
		for k, v := range req.PostForm {
			if len(v) > 0 {
				ctx.Form[k] = v[0]
			}
		}
	}

	return nil
}
