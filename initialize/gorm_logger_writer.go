package initialize

import (
	"fmt"
	"go-admin/config"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type Writer struct {
	config config.GeneralDB
	writer logger.Writer
}

var _ logger.Writer = (*Writer)(nil)

func NewWriter(config config.GeneralDB, writer logger.Writer) *Writer {
	return &Writer{config: config, writer: writer}
}

// Printf 格式化打印日志
func (w *Writer) Printf(message string, data ...interface{}) {
	if w.config.LogZap {
		switch w.config.LogLevel() {
		case logger.Silent:
			zap.L().Debug(fmt.Sprintf(message, data...))
		case logger.Error:
			zap.L().Error(fmt.Sprintf(message, data...))
		case logger.Warn:
			zap.L().Warn(fmt.Sprintf(message, data...))
		case logger.Info:
			zap.L().Info(fmt.Sprintf(message, data...))
		default:
			zap.L().Info(fmt.Sprintf(message, data...))
		}
		return
	}
	w.writer.Printf(message, data...)
}
