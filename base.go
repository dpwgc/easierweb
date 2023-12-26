package easierweb

import (
	"encoding/json"
	"encoding/xml"
	"gopkg.in/yaml.v3"
	"os"
)

func ReadJsonConfig(file string, obj any) error {
	configBytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(configBytes, obj)
}

func ReadYamlConfig(file string, obj any) error {
	configBytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(configBytes, obj)
}

func ReadXmlConfig(file string, obj any) error {
	configBytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return xml.Unmarshal(configBytes, obj)
}
