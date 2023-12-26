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
	Header      KV
	Path        KV
	Query       KV
	Body        []byte
	CustomCache any
	Request     *http.Request
	code        int
	response    []byte
	index       int
	methods     []Method
}

// Next
// go to the next processing method
func (c *Context) Next() {
	c.index++
	for c.index < len(c.methods) {
		c.methods[c.index](c)
		c.index++
	}
}

// Abort
// stop continuing down execution
func (c *Context) Abort() {
	c.index = len(c.methods) + 1
}

type KV map[string]string

func (kv KV) Has(key string) bool {
	_, has := kv[key]
	return has
}

func (kv KV) Keys() []string {
	var ks []string
	for k, _ := range kv {
		ks = append(ks, k)
	}
	return ks
}

func (kv KV) Values() []string {
	var vs []string
	for _, v := range kv {
		vs = append(vs, v)
	}
	return vs
}

func (kv KV) GetString(key string) string {
	return kv[key]
}

func (kv KV) GetInt(key string) int {
	i, _ := strconv.Atoi(kv[key])
	return i
}

func (kv KV) GetInt64(key string) int64 {
	i, _ := strconv.ParseInt(kv[key], 10, 64)
	return i
}

func (kv KV) GetFloat32(key string) float32 {
	f, _ := strconv.ParseFloat(key, 32)
	return float32(f)
}

func (kv KV) GetFloat64(key string) float64 {
	f, _ := strconv.ParseFloat(key, 64)
	return f
}

func (c *Context) BindJson(obj any) error {
	return json.Unmarshal(c.Body, obj)
}

func (c *Context) BindYaml(obj any) error {
	return yaml.Unmarshal(c.Body, obj)
}

func (c *Context) BindXml(obj any) error {
	return xml.Unmarshal(c.Body, obj)
}

func (c *Context) WriteJson(code int, obj any) {
	marshal, _ := json.Marshal(obj)
	c.Write(code, marshal)
}

func (c *Context) WriteYaml(code int, obj any) {
	marshal, _ := yaml.Marshal(obj)
	c.Write(code, marshal)
}

func (c *Context) WriteXml(code int, obj any) {
	marshal, _ := xml.Marshal(obj)
	c.Write(code, marshal)
}

func (c *Context) WriteString(code int, obj string) {
	c.Write(code, []byte(obj))
}

func (c *Context) Write(code int, data []byte) {
	c.code = code
	c.response = data
}
