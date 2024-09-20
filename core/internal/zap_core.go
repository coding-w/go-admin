package internal

import (
	"go-admin/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

// ZapCore 自定义日志核心 实现了 Core
type ZapCore struct {
	level zapcore.Level
	zapcore.Core
}

var _ zapcore.Core = (*ZapCore)(nil)

// Enabled 方法判断日志级别是否启用
func (z *ZapCore) Enabled(level zapcore.Level) bool {
	return z.level == level
}

// With 方法添加字段
// 将字段附加到日志条目中，返回一个新的zapcore.Core实例
func (z *ZapCore) With(fields []zapcore.Field) zapcore.Core {
	return z.Core.With(fields)
}

// Check 方法检查日志级别是否启用，如果启用，则返回一个新的zapcore.CheckedEntry实例
func (z *ZapCore) Check(entry zapcore.Entry, check *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if z.Enabled(entry.Level) {
		return check.AddCore(entry, z)
	}
	return check
}

// Write 方法将日志条目写入日志核心
// 遍历字段，如果存在business、folder或directory字段，根据字段值更新同步器和日志核心
// 调用Core.Write方法将条目写入
func (z *ZapCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	for i := 0; i < len(fields); i++ {
		if fields[i].Key == "business" || fields[i].Key == "folder" || fields[i].Key == "directory" {
			syncer := z.WriteSyncer(fields[i].String)
			z.Core = zapcore.NewCore(global.GA_CONFIG.Zap.Encoder(), syncer, z.level)
		}
	}
	return z.Core.Write(entry, fields)
}

// Sync 方法同步日志  调用Core.Sync方法，确保所有缓冲区的数据写入磁盘
func (z *ZapCore) Sync() error {
	return z.Core.Sync()
}

func NewZapCore(level zapcore.Level) *ZapCore {
	entry := &ZapCore{
		level: level,
	}
	// 调用WriteSyncer方法获取日志的写入同步器
	syncer := entry.WriteSyncer()
	// 创建一个函数，用于判断是否启用指定日志级别
	levelEnabler := zap.LevelEnablerFunc(
		func(l zapcore.Level) bool {
			return l == level
		},
	)
	// 创建一个新的日志核心，传入编码器、同步器和级别启用器
	entry.Core = zapcore.NewCore(global.GA_CONFIG.Zap.Encoder(), syncer, levelEnabler)
	return entry
}

func (z *ZapCore) WriteSyncer(formats ...string) zapcore.WriteSyncer {
	lc := NewLogCutter(
		global.GA_CONFIG.Zap.Director,
		z.level.String(),
		global.GA_CONFIG.Zap.RetentionDay,
		LogCutterWithLayout(time.Now().Format(time.DateOnly)),
		LogCutterWithFormats(formats...),
	)
	if global.GA_CONFIG.Zap.LogInConsole {
		multiWriteSyncer := zapcore.NewMultiWriteSyncer(os.Stdout, lc)
		return zapcore.AddSync(multiWriteSyncer)
	}
	return zapcore.AddSync(lc)
}
