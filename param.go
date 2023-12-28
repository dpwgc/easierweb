package easierweb

import (
	"encoding/json"
	"encoding/xml"
	"gopkg.in/yaml.v3"
	"strconv"
)

// 表单参数

type FormKV map[string]string

func (kv FormKV) Set(key string, value string) FormKV {
	if kv == nil {
		return kv
	}
	kv[key] = value
	return kv
}

func (kv FormKV) Get(key string) string {
	if kv == nil {
		return ""
	}
	return kv[key]
}

func (kv FormKV) Del(key string) FormKV {
	if kv == nil {
		return kv
	}
	delete(kv, key)
	return kv
}

func (kv FormKV) Has(key string) bool {
	if kv == nil {
		return false
	}
	_, has := kv[key]
	return has
}

func (kv FormKV) Keys() []string {
	if kv == nil {
		return nil
	}
	var ks = make([]string, 0, len(kv))
	for k := range kv {
		ks = append(ks, k)
	}
	return ks
}

func (kv FormKV) Values() []string {
	if kv == nil {
		return nil
	}
	var vs = make([]string, 0, len(kv))
	for _, v := range kv {
		vs = append(vs, v)
	}
	return vs
}

func (kv FormKV) Int(key string) int {
	i, err := kv.ParseInt(key)
	if err != nil {
		panic(err)
	}
	return i
}

func (kv FormKV) Int32(key string) int32 {
	i, err := kv.ParseInt32(key)
	if err != nil {
		panic(err)
	}
	return i
}

func (kv FormKV) Int64(key string) int64 {
	i, err := kv.ParseInt64(key)
	if err != nil {
		panic(err)
	}
	return i
}

func (kv FormKV) Float32(key string) float32 {
	f, err := kv.ParseFloat32(key)
	if err != nil {
		panic(err)
	}
	return f
}

func (kv FormKV) Float64(key string) float64 {
	f, err := kv.ParseFloat64(key)
	if err != nil {
		panic(err)
	}
	return f
}

func (kv FormKV) GetInt(key string) int {
	i, _ := kv.ParseInt(key)
	return i
}

func (kv FormKV) GetInt32(key string) int32 {
	i, _ := kv.ParseInt32(key)
	return i
}

func (kv FormKV) GetInt64(key string) int64 {
	i, _ := kv.ParseInt64(key)
	return i
}

func (kv FormKV) GetFloat32(key string) float32 {
	f, _ := kv.ParseFloat32(key)
	return f
}

func (kv FormKV) GetFloat64(key string) float64 {
	f, _ := kv.ParseFloat64(key)
	return f
}

func (kv FormKV) ParseInt(key string) (int, error) {
	return strconv.Atoi(kv.Get(key))
}

func (kv FormKV) ParseInt32(key string) (int32, error) {
	i, err := strconv.ParseInt(kv.Get(key), 10, 32)
	return int32(i), err
}

func (kv FormKV) ParseInt64(key string) (int64, error) {
	return strconv.ParseInt(kv.Get(key), 10, 64)
}

func (kv FormKV) ParseFloat32(key string) (float32, error) {
	f, err := strconv.ParseFloat(kv.Get(key), 32)
	return float32(f), err
}

func (kv FormKV) ParseFloat64(key string) (float64, error) {
	return strconv.ParseFloat(kv.Get(key), 64)
}

// 缓存参数

type CacheKV map[string]any

func (kv CacheKV) Set(key string, value any) CacheKV {
	kv[key] = value
	return kv
}

func (kv CacheKV) Get(key string) any {
	return kv[key]
}

func (kv CacheKV) Del(key string) CacheKV {
	delete(kv, key)
	return kv
}

func (kv CacheKV) Has(key string) bool {
	_, has := kv[key]
	return has
}

func (kv CacheKV) Keys() []string {
	var ks = make([]string, 0, len(kv))
	for k := range kv {
		ks = append(ks, k)
	}
	return ks
}

func (kv CacheKV) Values() []any {
	var vs = make([]any, 0, len(kv))
	for _, v := range kv {
		vs = append(vs, v)
	}
	return vs
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
