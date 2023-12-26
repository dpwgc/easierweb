package easierweb

import (
	"encoding/json"
	"encoding/xml"
	"gopkg.in/yaml.v3"
	"net/http"
	"strconv"
)

type Method func(ctx *Context)

type Context struct {
	Header         KV
	Path           KV
	Query          KV
	Form           KV
	Body           []byte
	Code           int
	Result         []byte
	CustomCache    any
	Request        *http.Request
	ResponseWriter http.ResponseWriter
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

type KV map[string]string

func (kv KV) Has(key string) bool {
	_, has := kv[key]
	return has
}

func (kv KV) Keys() []string {
	var ks = make([]string, 0, len(kv))
	for k := range kv {
		ks = append(ks, k)
	}
	return ks
}

func (kv KV) Values() []string {
	var vs = make([]string, 0, len(kv))
	for _, v := range kv {
		vs = append(vs, v)
	}
	return vs
}

func (kv KV) GetInt(key string) int {
	i, _ := strconv.Atoi(kv[key])
	return i
}

func (kv KV) GetInt32(key string) int32 {
	i, _ := strconv.ParseInt(kv[key], 10, 32)
	return int32(i)
}

func (kv KV) GetFloat32(key string) float32 {
	f, _ := strconv.ParseFloat(key, 32)
	return float32(f)
}

func (kv KV) GetInt64(key string) int64 {
	i, _ := strconv.ParseInt(kv[key], 10, 64)
	return i
}

func (kv KV) GetFloat64(key string) float64 {
	f, _ := strconv.ParseFloat(key, 64)
	return f
}

func (c *Context) BindJsonBody(obj any) error {
	return json.Unmarshal(c.Body, obj)
}

func (c *Context) BindYamlBody(obj any) error {
	return yaml.Unmarshal(c.Body, obj)
}

func (c *Context) BindXmlBody(obj any) error {
	return xml.Unmarshal(c.Body, obj)
}

func (c *Context) WriteJsonResult(code int, obj any) {
	marshal, _ := json.Marshal(obj)
	c.WriteResult(code, marshal)
}

func (c *Context) WriteYamlResult(code int, obj any) {
	marshal, _ := yaml.Marshal(obj)
	c.WriteResult(code, marshal)
}

func (c *Context) WriteXmlResult(code int, obj any) {
	marshal, _ := xml.Marshal(obj)
	c.WriteResult(code, marshal)
}

func (c *Context) WriteResult(code int, data []byte) {
	c.Code = code
	c.Result = data
}
