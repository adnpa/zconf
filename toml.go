package zconf

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type UnmatchedTomlKeysError struct {
	Keys []toml.Key
}

func (e *UnmatchedTomlKeysError) Error() string {
	return fmt.Sprintf("There are keys in the config file that do not match any field in the given struct: %v", e.Keys)
}

func UnmarshalToml(data []byte, config interface{}, errorOnUnmatchedKeys bool) error {
	metaData, err := toml.Decode(string(data), config)
	if err == nil && len(metaData.Undecoded()) > 0 && errorOnUnmatchedKeys {
		return &UnmatchedTomlKeysError{Keys: metaData.Undecoded()}
	}
	return nil
}
