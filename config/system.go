package config

type System struct {
	DbType   string `mapstructure:"db-type" json:"db-type" yaml:"db-type"`
	OssType  string `mapstructure:"oss-type" json:"oss-type" yaml:"oss-type"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	Env      string `mapstructure:"env" json:"env" yaml:"env"`
	UseRedis bool   `mapstructure:"use-redis" json:"use-redis" yaml:"use-redis"`
	UesMongo bool   `mapstructure:"use-mongo" json:"use-mongo" yaml:"use-mongo"`
	// 多点登录拦截
	UseMultipoint bool   `mapstructure:"use-multipoint" json:"use-multipoint" yaml:"use-multipoint"`
	IpLimitCount  int    `mapstructure:"ip-limit-count" json:"ip-limit-count" yaml:"ip-limit-count"`
	IpLimitTime   int    `mapstructure:"ip-limit-time" json:"ip-limit-time" yaml:"ip-limit-time"`
	RouterPrefix  string `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"`
}
