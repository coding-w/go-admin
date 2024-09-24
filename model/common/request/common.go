package request

import "gorm.io/gorm"

// PageInfo 分页
type PageInfo struct {
	Page     int    `json:"page" form:"page"`         // 页码
	PageSize int    `json:"pageSize" form:"pageSize"` // 每页大小
	Keyword  string `json:"keyword" form:"keyword"`   // 关键字
}

// PageScopes Scopes
func (pi *PageInfo) PageScopes() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pi.Page <= 0 {
			pi.Page = 1
		}
		switch {
		case pi.PageSize > 100:
			pi.PageSize = 100
		case pi.PageSize <= 0:
			pi.PageSize = 10
		}
		offset := (pi.Page - 1) * pi.PageSize
		return db.Offset(offset).Limit(pi.PageSize)
	}
}

// GetById 根据id查询
type GetById struct {
	ID int `json:"id" form:"id"` // 主键ID
}

func (r *GetById) Uint() uint {
	return uint(r.ID)
}

// IdsReq id 数组 struct
type IdsReq struct {
	Ids []int `json:"ids" form:"ids"`
}

// GetAuthorityId 权限id
type GetAuthorityId struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"` // 角色ID
}

type Empty struct{}
