package utils

import (
	"errors"
	"os"
)

// PathExists 检查指定路径是否存在，并区分路径是否为目录
func PathExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		// 路径存在，判断是否为目录
		if fi.IsDir() {
			return true, nil
		}
		// 路径存在但不是目录
		return false, errors.New("存在同名文件")
	}
	if os.IsNotExist(err) {
		// 路径不存在
		return false, nil
	}
	// 其他错误
	return false, err
}
