package zconf

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func UnmarshalJson(data []byte, config interface{}, errorOnUnmatchedKeys bool) (err error) {
	fmt.Println("data", data)
	dec := json.NewDecoder(bytes.NewReader(data))
	if errorOnUnmatchedKeys {
		dec.DisallowUnknownFields()
	}
	err = dec.Decode(&config)

	fmt.Println("end", config)
	return
}
