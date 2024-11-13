package zconf

import (
	"fmt"
	"testing"
	"time"
)

type TestConfig struct {
	Name string
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

}
