package core

import (
	"fmt"
	"go-admin/core/internal"
	"go-admin/global"
	"go-admin/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func Zap() (logger *zap.Logger) {
	ok, _ := utils.PathExists(global.GA_CONFIG.Zap.Director)
	if !ok {
		fmt.Printf("create %v directory\n", global.GA_CONFIG.Zap.Director)
		_ = os.Mkdir(global.GA_CONFIG.Zap.Director, os.ModePerm)
	}
	levels := global.GA_CONFIG.Zap.Levels()
	length := len(levels)
	cores := make([]zapcore.Core, 0, length)
	for i := 0; i < length; i++ {
		core := internal.NewZapCore(levels[i])
		cores = append(cores, core)
	}
	logger = zap.New(zapcore.NewTee(cores...))
	// 添加caller, 是否打印行数
	if global.GA_CONFIG.Zap.ShowLine {
		logger.WithOptions(zap.AddCaller())
	}
	return logger
}
