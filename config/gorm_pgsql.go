package config

import (
	"fmt"
	"gorm.io/gorm/logger"
	"strings"
)

type Pgsql struct {
	Host        string `mapstructure:"host" json:"host" yaml:"host"`                            // 数据库地址
	Port        int    `mapstructure:"port" json:"port" yaml:"port"`                            // 数据库端口
	Config      string `mapstructure:"config" json:"config" yaml:"config"`                      // 高级配置
	Dbname      string `mapstructure:"db-name" json:"db-name" yaml:"db-name"`                   // 数据库名
	Username    string `mapstructure:"username" json:"username" yaml:"username"`                // 用户名
	Password    string `mapstructure:"password" json:"password" yaml:"password"`                // 密码
	MaxIdleConn int    `mapstructure:"max-idle-conn" json:"max-idle-conn" yaml:"max-idle-conn"` // 空闲连接数
	MaxOpenConn int    `mapstructure:"max-open-conn" json:"max-open-conn" yaml:"max-open-conn"` // 打开连接数
	LogMode     string `mapstructure:"log-mode" json:"log-mode" yaml:"log-mode"`                // 是否开启Gorm全局日志
	LogZap      bool   `mapstructure:"log-zap" json:"log-zap" yaml:"log-zap"`                   // 是否通过zap写入日志文件
	Prefix      string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                      // 数据库前缀
	Singular    bool   `mapstructure:"singular" json:"singular" yaml:"singular"`                // 是否开启全局禁用复数，true表示开启
}

func (p *Pgsql) Dsn() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d %s", p.Host, p.Username, p.Password, p.Dbname, p.Port, p.Config)
}

func (p *Pgsql) LogLevel() logger.LogLevel {
	switch strings.ToLower(p.LogMode) {
	case "silent", "Silent":
		return logger.Silent
	case "error", "Error":
		return logger.Error
	case "warn", "Warn":
		return logger.Warn
	case "info", "Info":
		return logger.Info
	default:
		return logger.Info
	}
}
