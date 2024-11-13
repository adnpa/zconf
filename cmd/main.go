package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type TestConfig struct {
	name string `json:"Name"`
	Addr string `json:"Addr"`
	Port int    `json:"Port"`
}

func main() {
	var conf TestConfig
	filePath := "./test_config/conf.json"
	file, _ := os.Open(filePath)
	defer file.Close()

	data, _ := io.ReadAll(file)
	json.Unmarshal(data, &conf)

	fmt.Println(conf)
}
