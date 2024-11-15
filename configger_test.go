package zconf

import (
	"fmt"
	"testing"
	"time"
)

type AppConfig struct {
	App TestConfig
}

type TestConfig struct {
	Name string `default:"abc"`
	Addr string
	Port int
}

func Test(t *testing.T) {
	fmt.Println(time.Now())
}

func TestUnmarshalJson(t *testing.T) {
	var conf TestConfig
	filePath := "./test_config/conf.json"
	Load(&conf, filePath)
	fmt.Println(conf)
}

func TestUnmarshalYaml(t *testing.T) {
	var conf AppConfig
	filePath := "./test_config/conf.yml"
	Load(&conf, filePath)
	fmt.Println(conf)
}

func TestUnmarshalToml(t *testing.T) {
	var conf AppConfig
	filePath := "./test_config/conf.toml"
	Load(&conf, filePath)
	fmt.Println(conf)
}

func TestAutoload(t *testing.T) {
	var conf TestConfig
	filePath := "./test_config/conf.json"
	callback := func(c interface{}) {
		fmt.Println("callback")
		fmt.Println(c)
	}

	configger := NewWithOption(&Option{
		AutoReload:         true,
		AutoReloadInterval: 5 * time.Second,
		AutoReloadCallback: callback,
	})
	configger.Load(&conf, filePath)

	select {}
}
