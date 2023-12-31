package easierweb

import (
	"encoding/json"
	"encoding/xml"
	"gopkg.in/yaml.v3"
	"strconv"
)

// 表单参数

type Params map[string]string

func (kv Params) Set(key string, value string) Params {
	if kv == nil {
		return kv
	}
	kv[key] = value
	return kv
}

func (kv Params) Get(key string) string {
	if kv == nil {
		return ""
	}
	return kv[key]
}

func (kv Params) Del(key string) Params {
	if kv == nil {
		return kv
	}
	delete(kv, key)
	return kv
}

func (kv Params) Has(key string) bool {
	if kv == nil {
		return false
	}
	_, has := kv[key]
	return has
}

func (kv Params) Keys() []string {
	if kv == nil {
		return nil
	}
	var ks = make([]string, 0, len(kv))
	for k := range kv {
		ks = append(ks, k)
	}
	return ks
}

func (kv Params) Values() []string {
	if kv == nil {
		return nil
	}
	var vs = make([]string, 0, len(kv))
	for _, v := range kv {
		vs = append(vs, v)
	}
	return vs
}

func (kv Params) Int(key string) int {
	i, err := kv.ParseInt(key)
	if err != nil {
		panic(err)
	}
	return i
}

func (kv Params) Int32(key string) int32 {
	i, err := kv.ParseInt32(key)
	if err != nil {
		panic(err)
	}
	return i
}

func (kv Params) Int64(key string) int64 {
	i, err := kv.ParseInt64(key)
	if err != nil {
		panic(err)
	}
	return i
}

func (kv Params) Float32(key string) float32 {
	f, err := kv.ParseFloat32(key)
	if err != nil {
		panic(err)
	}
	return f
}

func (kv Params) Float64(key string) float64 {
	f, err := kv.ParseFloat64(key)
	if err != nil {
		panic(err)
	}
	return f
}

func (kv Params) GetInt(key string) int {
	i, _ := kv.ParseInt(key)
	return i
}

func (kv Params) GetInt32(key string) int32 {
	i, _ := kv.ParseInt32(key)
	return i
}

func (kv Params) GetInt64(key string) int64 {
	i, _ := kv.ParseInt64(key)
	return i
}

func (kv Params) GetFloat32(key string) float32 {
	f, _ := kv.ParseFloat32(key)
	return f
}

func (kv Params) GetFloat64(key string) float64 {
	f, _ := kv.ParseFloat64(key)
	return f
}

func (kv Params) ParseInt(key string) (int, error) {
	return strconv.Atoi(kv.Get(key))
}

func (kv Params) ParseInt32(key string) (int32, error) {
	i, err := strconv.ParseInt(kv.Get(key), 10, 32)
	return int32(i), err
}

func (kv Params) ParseInt64(key string) (int64, error) {
	return strconv.ParseInt(kv.Get(key), 10, 64)
}

func (kv Params) ParseFloat32(key string) (float32, error) {
	f, err := strconv.ParseFloat(kv.Get(key), 32)
	return float32(f), err
}

func (kv Params) ParseFloat64(key string) (float64, error) {
	return strconv.ParseFloat(kv.Get(key), 64)
}

type Data []byte

func (d Data) String() string {
	return string(d)
}

func (d Data) ParseJSON(obj any) error {
	return json.Unmarshal(d, obj)
}

func (d Data) ParseYAML(obj any) error {
	return yaml.Unmarshal(d, obj)
}

func (d Data) ParseXML(obj any) error {
	return xml.Unmarshal(d, obj)
}
