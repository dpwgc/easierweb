package easierweb

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/websocket"
	"gopkg.in/yaml.v3"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type Context struct {
	Header         FormKV
	Path           FormKV
	Query          FormKV
	Form           FormKV
	Body           Data
	Code           int
	Result         Data
	CustomCache    CacheKV
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	WebsocketConn  *websocket.Conn
	index          int
	methods        []Method
}

func (c *Context) Next() {
	c.index++
	for c.index < len(c.methods) {
		c.methods[c.index](c)
		c.index++
	}
}

func (c *Context) Abort() {
	c.index = len(c.methods) + 1
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
	marshal, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	c.Write(code, marshal)
}

func (c *Context) WriteYAML(code int, obj any) {
	marshal, err := yaml.Marshal(obj)
	if err != nil {
		panic(err)
	}
	c.Write(code, marshal)
}

func (c *Context) WriteXML(code int, obj any) {
	marshal, err := xml.Marshal(obj)
	if err != nil {
		panic(err)
	}
	c.Write(code, marshal)
}

func (c *Context) Redirect(code int, url string) {
	http.Redirect(c.ResponseWriter, c.Request, url, code)
}

func (c *Context) WriteLocalFile(contentType, fileName, localFilePath string) {
	fileBytes, err := os.ReadFile(localFilePath)
	if err != nil {
		panic(err)
	}
	c.WriteFile(contentType, fileName, fileBytes)
}

func (c *Context) WriteFile(contentType, fileName string, fileBytes []byte) {
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

func (c *Context) WriteHTML(code int, html string) {
	c.ResponseWriter.Header().Set("Content-Type", "text/html")
	c.Write(code, []byte(html))
}

func (c *Context) WriteString(code int, text string) {
	c.ResponseWriter.Header().Set("Content-Type", "text/plain")
	c.Write(code, []byte(text))
}

func (c *Context) NoContent(code int) {
	c.Write(code, nil)
}

func (c *Context) Write(code int, data []byte) {
	c.Code = code
	c.Result = data
	c.ResponseWriter.WriteHeader(code)
	_, err := c.ResponseWriter.Write(data)
	if err != nil {
		panic(err)
	}
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

func (c *Context) Send(msg []byte) error {
	_, err := c.WebsocketConn.Write(msg)
	if err != nil {
		return err
	}
	return nil
}