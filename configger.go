package zconf

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"reflect"
	"strings"
	"time"
)

// exported func

func Load(config interface{}, files ...string) error {
	return NewWithOption(nil).Load(config, files...)
}

type Configger struct {
	*Option
	confModTimes map[string]time.Time

	Fs fs.FS
}

func NewWithOption(option *Option) *Configger {
	return &Configger{Option: option}
}

func (c *Configger) Load(config interface{}, files ...string) (err error) {
	defaultVal := reflect.Indirect(reflect.ValueOf(config))
	if !defaultVal.CanAddr() {
		return fmt.Errorf("Config %v should be addressable", config)
	}
	_, err = c.load(config, false, files...)

	if c.Option.AutoReload {
		go func() {
			//临时对象用于比较 是否修改
			reflectPtr := reflect.New(reflect.ValueOf(config).Elem().Type())
			reflectPtr.Elem().Set(defaultVal)

			timer := time.NewTimer(c.Option.AutoReloadInterval)
			for range timer.C {
				var changed bool
				changed, err = c.load(reflectPtr.Interface(), true, files...)
				fmt.Println(changed)
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
	// 用文件修改时间判断 配置是否修改
	configFiles, confModTimes := c.getConfigurationFiles(watch, files...)

	if watch {
		var changed = false
		for f, t := range confModTimes {
			if v, ok := c.confModTimes[f]; !ok || t.After(v) {
				changed = true
			}
		}
		if !changed {
			return false, nil
		}
	}

	for _, file := range configFiles {
		c.processFile(config, file, watch)
	}
	c.confModTimes = confModTimes
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
		return UnmarshalJson(data, config, watch)
	default:
		return errors.New("unknown config file type: " + filePath)
	}

}
