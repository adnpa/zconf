package zconf

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"os"
	"reflect"
	"strings"
	"time"
)

func Load(config interface{}, files ...string) error {
	return NewWithOption(nil).Load(config, files...)
}

type Configger struct {
	*Option                           //选项
	confModTimes map[string]time.Time //配置修改时间
	Fs           fs.FS                //文件系统 需要访问其他磁盘的情况下
}

func NewWithOption(option *Option) *Configger {
	if option == nil {
		option = &Option{}
	}
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
			//临时对象 用于比较修改时间
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

	err := c.processDefault(config)
	if err != nil {
		fmt.Println("default err", err)
		return false, err
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
	case strings.HasSuffix(filePath, ".yaml") || strings.HasSuffix(filePath, ".yml"):
		return UnmarshalYaml(data, config, watch)
	case strings.HasSuffix(filePath, ".json"):
		return UnmarshalJson(data, config, watch)
	case strings.HasSuffix(filePath, ".toml"):
		return UnmarshalToml(data, config, watch)
	default:
		return errors.New("unknown config file type: " + filePath)
	}

}

func (c *Configger) processDefault(config interface{}) (err error) {
	configValue := reflect.Indirect(reflect.ValueOf(config))
	if configValue.Kind() != reflect.Struct {
		return errors.New("invalid config, should be struct")
	}

	configType := configValue.Type()
	for i := 0; i < configType.NumField(); i++ {
		var (
			fieldStruct = configType.Field(i)
			field       = configValue.Field(i)
		)

		if !field.CanAddr() || !field.CanInterface() {
			continue
		}

		isEmpty := reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface())
		if isEmpty {
			//如果配置为空则填充默认值
			if defaultVal := fieldStruct.Tag.Get("default"); defaultVal != "" {
				if err = yaml.Unmarshal([]byte(defaultVal), field.Addr().Interface()); err != nil {
					return
				}
			}
		}

		for field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		//处理嵌套struct
		switch field.Kind() {
		case reflect.Struct:
			if err = c.processDefault(field.Addr().Interface()); err != nil {
				return
			}
		case reflect.Slice:
			for i := 0; i < field.Len(); i++ {
				if reflect.Indirect(field.Index(i)).Interface() == reflect.String {
					if err = c.processDefault(field.Index(i).Addr().Interface()); err != nil {
						return
					}
				}
			}
		}

	}
	return
}
