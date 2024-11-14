package zconf

import (
	"bytes"
	"encoding/json"
)

func UnmarshalJson(data []byte, config interface{}, errorOnUnmatchedKeys bool) (err error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	if errorOnUnmatchedKeys {
		dec.DisallowUnknownFields()
	}
	err = dec.Decode(&config)
	return
}
