package easierweb

import (
	"encoding/json"
	"encoding/xml"
	"gopkg.in/yaml.v3"
	"mime/multipart"
	"net/http"
	"strconv"
)

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

func (kv KV) Int(key string) int {
	i, _ := strconv.Atoi(kv[key])
	return i
}

func (kv KV) Int32(key string) int32 {
	i, _ := strconv.ParseInt(kv[key], 10, 32)
	return int32(i)
}

func (kv KV) Int64(key string) int64 {
	i, _ := strconv.ParseInt(kv[key], 10, 64)
	return i
}

func (kv KV) Float32(key string) float32 {
	f, _ := strconv.ParseFloat(key, 32)
	return float32(f)
}

func (kv KV) Float64(key string) float64 {
	f, _ := strconv.ParseFloat(key, 64)
	return f
}

func (kv KV) GetInt(key string) (int, error) {
	return strconv.Atoi(kv[key])
}

func (kv KV) GetInt32(key string) (int32, error) {
	i, err := strconv.ParseInt(kv[key], 10, 32)
	return int32(i), err
}

func (kv KV) GetInt64(key string) (int64, error) {
	return strconv.ParseInt(kv[key], 10, 64)
}

func (kv KV) GetFloat32(key string) (float32, error) {
	f, err := strconv.ParseFloat(key, 32)
	return float32(f), err
}

func (kv KV) GetFloat64(key string) (float64, error) {
	return strconv.ParseFloat(key, 64)
}

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

func (c *Context) GetFile(key string) (multipart.File, error) {
	file, _, err := c.Request.FormFile(key)
	if err != nil {
		return nil, err
	}
	return file, nil
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
	marshal, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	c.Write(code, marshal)
}

func (c *Context) WriteYaml(code int, obj any) {
	marshal, err := yaml.Marshal(obj)
	if err != nil {
		panic(err)
	}
	c.Write(code, marshal)
}

func (c *Context) WriteXml(code int, obj any) {
	marshal, err := xml.Marshal(obj)
	if err != nil {
		panic(err)
	}
	c.Write(code, marshal)
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
