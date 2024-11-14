package zconf

import (
	"io/fs"
	"os"
	"time"
)

// return (config files, config-mod map)
func (c *Configger) getConfigurationFiles(watchMode bool, files ...string) ([]string, map[string]time.Time) {
	var resFiles []string
	var resFileModTimes = make(map[string]time.Time)

	stat := os.Stat
	if c.Fs != nil {
		stat = func(name string) (os.FileInfo, error) {
			return fs.Stat(c.Fs, name)
		}
	}

	for _, file := range files {
		if fileInfo, err := stat(file); err == nil && fileInfo.Mode().IsRegular() {
			resFiles = append(resFiles, file)
			resFileModTimes[file] = fileInfo.ModTime()
		}
	}

	return resFiles, resFileModTimes
}
