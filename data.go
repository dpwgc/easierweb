package easierweb

import (
	"encoding/json"
	"encoding/xml"
	"gopkg.in/yaml.v3"
)

type Data []byte

func (d *Data) ParseJSON(obj any) error {
	return json.Unmarshal(*d, obj)
}

func (d *Data) ParseYAML(obj any) error {
	return yaml.Unmarshal(*d, obj)
}

func (d *Data) ParseXML(obj any) error {
	return xml.Unmarshal(*d, obj)
}

func (d *Data) SaveJSON(obj any) error {
	marshal, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	d.Save(marshal)
	return nil
}

func (d *Data) SaveYAML(obj any) error {
	marshal, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}
	d.Save(marshal)
	return nil
}

func (d *Data) SaveXML(obj any) error {
	marshal, err := xml.Marshal(obj)
	if err != nil {
		return err
	}
	d.Save(marshal)
	return nil
}

func (d *Data) Save(new []byte) {
	*d = new
}
