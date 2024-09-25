package global

import (
	"time"

	"gorm.io/gorm"
)

type GA_Model struct {
	ID        uint           `gorm:"primarykey;column:id" json:"id"`    // 主键ID
	CreatedAt time.Time      `gorm:"column:created_at" json:"createAt"` // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updateAt"` // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at" json:"-"`  // 删除时间
}
