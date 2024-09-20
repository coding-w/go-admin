package utils

import (
	"strconv"
	"strings"
	"time"
)

// ParseDuration 格式化时间
func ParseDuration(d string) (time.Duration, error) {
	// 去除字符串两端的空白字符
	d = strings.TrimSpace(d)

	// 尝试使用标准的 time.ParseDuration 函数解析
	dr, err := time.ParseDuration(d)
	if err == nil {
		return dr, nil
	}

	// 如果字符串包含 'd'，则手动处理天数部分
	if strings.Contains(d, "d") {
		index := strings.Index(d, "d")
		// 提取天数部分并转换为整数
		hour, err := strconv.Atoi(d[:index])
		if err != nil {
			return 0, err
		}

		dr = time.Duration(hour) * 24 * time.Hour

		// 解析 'd' 后面的部分（小时、分钟等）
		ndr, err := time.ParseDuration(d[index+1:])
		if err != nil {
			return dr, nil
		}

		return dr + ndr, nil
	}

	// 尝试将整个字符串解析为整数，并将其转换为 time.Duration 类型
	dv, err := strconv.ParseInt(d, 10, 64)
	return time.Duration(dv), err
}
