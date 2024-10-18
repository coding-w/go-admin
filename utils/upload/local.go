package upload

import (
	"errors"
	"fmt"
	"go-admin/global"
	"go-admin/utils"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var mu sync.Mutex

type Local struct {
	OSS
}

func (l *Local) UploadFile(file *multipart.FileHeader) (string, string, error) {
	// 读取文件后缀
	ext := filepath.Ext(file.Filename)
	// 读取文件名并加密
	name := utils.MD5V([]byte(strings.TrimSuffix(file.Filename, ext)))
	// 拼接目标路径
	targetPath := filepath.Join(global.GA_CONFIG.Local.StorePath, time.Now().Format("2006/01/02"))
	filename := name + ext
	fullPath := filepath.Join(targetPath, filename)
	// 尝试创建此路径
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		global.GA_LOG.Error("failed to create directory", zap.String("path", targetPath), zap.Error(err))
		return "", "", fmt.Errorf("failed to create directory: %w", err)
	}
	// 打开上传的文件
	srcFile, err := file.Open()
	if err != nil {
		global.GA_LOG.Error("failed to open uploaded file", zap.String("filename", file.Filename), zap.Error(err))
		return "", "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer srcFile.Close()
	// 创建目标文件
	destFile, err := os.Create(fullPath)
	if err != nil {
		global.GA_LOG.Error("failed to create destination file", zap.String("path", fullPath), zap.Error(err))
		return "", "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()
	// 传输（拷贝）文件
	if _, err := io.Copy(destFile, srcFile); err != nil {
		global.GA_LOG.Error("failed to copy file", zap.String("path", fullPath), zap.Error(err))
		return "", "", fmt.Errorf("failed to copy file: %w", err)
	}
	return fullPath, filename, nil
}

func (l *Local) DeleteFile(key string) error {
	// 检查 key 是否为空
	if key == "" {
		return errors.New("key不能为空")
	}

	// 验证 key 是否包含非法字符或尝试访问存储路径之外的文件
	if strings.Contains(key, "..") || strings.ContainsAny(key, `\/:*?"<>|`) {
		return errors.New("非法的key")
	}

	p := filepath.Join(global.GA_CONFIG.Local.StorePath, key)

	// 检查文件是否存在
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return errors.New("文件不存在")
	}

	// 使用文件锁防止并发删除
	mu.Lock()
	defer mu.Unlock()

	err := os.Remove(p)
	if err != nil {
		return errors.New("文件删除失败: " + err.Error())
	}
	return nil
}
