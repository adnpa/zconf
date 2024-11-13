package zconf

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// exported func

func Load(config interface{}, files ...string) error {
	return NewWithOption(nil).Load(config, files...)
}

type Configger struct {
	*Option
}

func NewWithOption(option *Option) *Configger {
	return &Configger{option}
}

func (c *Configger) Load(config interface{}, files ...string) (err error) {
	_, err = c.load(config, false, files...)

	if c.Option.AutoReload {
		go func() {
			timer := time.NewTimer(c.Option.AutoReloadInterval)
			for range timer.C {
				var changed bool
				changed, err = c.load(config, true, files...)
				if changed && c.Option.AutoReloadCallback != nil && err == nil {
					c.Option.AutoReloadCallback(config)
				} else if err != nil {
					fmt.Printf("load config file failed, err:%v\n", err)
				}
				timer.Reset(c.Option.AutoReloadInterval)
			}
		}()
	}
	return
}

func (c *Configger) load(config interface{}, watch bool, files ...string) (bool, error) {

	for _, file := range files {
		c.processFile(config, file, watch)
	}

	return true, nil
}

func (c *Configger) processFile(config interface{}, filePath string, watch bool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	data, _ := io.ReadAll(file)

	switch {
	//case strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml"):
	//	return yaml.Unmarshal(data, config)
	case strings.HasSuffix(filePath, ".json"):
		fmt.Print("jsonjson")
		return UnmarshalJson(data, config, watch)
	default:

		return errors.New("unknown config file type: " + filePath)
	}

}
