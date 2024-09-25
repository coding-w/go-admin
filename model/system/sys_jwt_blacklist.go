package system

import "go-admin/global"

type JwtBlacklist struct {
	global.GA_Model
	Jwt string `gorm:"type:text;comment:jwt"`
}
