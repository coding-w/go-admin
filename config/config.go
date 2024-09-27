package config

type Config struct {
	System    System    `mapstructure:"system" json:"system" yaml:"system"`          // 系统配置
	Zap       Zap       `mapstructure:"zap" json:"zap" yaml:"zap"`                   // 日志配置
	Captcha   Captcha   `mapstructure:"captcha" json:"captcha" yaml:"captcha"`       // 验证码配置
	Pgsql     Pgsql     `mapstructure:"pgsql" json:"pgsql" yaml:"pgsql"`             //数据库配置
	Mysql     Mysql     `mapstructure:"mysql" json:"mysql" yaml:"mysql"`             // mysql
	JWT       JWT       `mapstructure:"jwt" json:"jwt" yaml:"jwt"`                   // jwt 配置
	Local     Local     `mapstructure:"local" json:"local" yaml:"local"`             // 本地文件上传配置
	AliyunOSS AliyunOSS `mapstructure:"aliyunoss" json:"aliyunoss" yaml:"aliyunoss"` // 阿里对象存储 配置
	Redis     Redis     `mapstructure:"redis" json:"redis" yaml:"redis"`             // Redis 配置
}
