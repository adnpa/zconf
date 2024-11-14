package zconf

import (
	"fmt"
	"testing"
	"time"
)

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

//func TestDefault(t *testing.T) {
//	var conf TestConfig
//	file
//}
