package zconf

import (
	"bytes"
	"gopkg.in/yaml.v3"
)

func UnmarshalYaml(data []byte, config interface{}, errorOnUnmatchedKeys bool) (err error) {
	if errorOnUnmatchedKeys {
		dec := yaml.NewDecoder(bytes.NewBuffer(data))
		dec.KnownFields(true)
		return dec.Decode(config)
	}
	return yaml.Unmarshal(data, config)
}
