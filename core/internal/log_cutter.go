package internal

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

// LogCutter 日志切割器 实现了 WriteSyncer
type LogCutter struct {
	level        string        // 日志级别(debug, info, warn, error, dpanic, panic, fatal)
	layout       string        // 时间格式 2006-01-02 15:04:05
	formats      []string      // 自定义参数([]string{Director,"2006-01-02", "business"(此参数可不写), level+".log"}
	director     string        // 日志文件夹
	retentionDay int           //日志保留天数
	file         *os.File      // 文件句柄
	mutex        *sync.RWMutex // 读写锁
}

type LogCutterOption func(*LogCutter)

// LogCutterWithLayout 时间格式
func LogCutterWithLayout(layout string) LogCutterOption {
	return func(lc *LogCutter) {
		lc.layout = layout
	}
}

// LogCutterWithFormats 格式化参数
func LogCutterWithFormats(format ...string) LogCutterOption {
	return func(lc *LogCutter) {
		if len(format) > 0 {
			lc.formats = format
		}
	}
}

func NewLogCutter(director string, level string, retentionDay int, options ...LogCutterOption) *LogCutter {
	cutter := &LogCutter{
		level:        level,
		director:     director,
		retentionDay: retentionDay,
		mutex:        new(sync.RWMutex),
	}
	for i := 0; i < len(options); i++ {
		options[i](cutter)
	}
	return cutter
}

func (lc *LogCutter) Write(bytes []byte) (n int, err error) {
	lc.mutex.Lock()
	defer func() {
		if lc.file != nil {
			_ = lc.file.Close()
			lc.file = nil
		}
		lc.mutex.Unlock()
	}()
	length := len(lc.formats)
	values := make([]string, 0, length+3)
	// 日志文件夹
	values = append(values, lc.director)
	if lc.layout != "" {
		values = append(values, lc.layout)
	}
	for i := 0; i < length; i++ {
		values = append(values, lc.level)
	}
	values = append(values, lc.level+".log")
	filename := filepath.Join(values...)
	dir := filepath.Dir(filename)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return 0, err
	}
	// 删除旧文件
	err = removeFolders(lc.director, lc.retentionDay)
	if err != nil {
		return 0, err
	}
	lc.file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return 0, err
	}

	return lc.file.Write(bytes)

}

// Sync 方法用于同步文件数据，将内存中的数据写入磁盘
func (lc *LogCutter) Sync() error {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	if lc.file != nil {
		return lc.file.Sync()
	}
	return nil
}

// removeFolders 增加日志目录文件清理 小于等于零的值默认忽略不再处理
func removeFolders(dir string, days int) error {
	if days <= 0 {
		return nil
	}
	cutoff := time.Now().AddDate(0, 0, -days)
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.ModTime().Before(cutoff) && path != dir {
			err = os.RemoveAll(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
