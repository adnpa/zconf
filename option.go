package zconf

import "time"

type Option struct {
	AutoReload         bool
	AutoReloadInterval time.Duration
	AutoReloadCallback func(config interface{})
}
