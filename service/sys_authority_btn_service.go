package service

import (
	"errors"
	"go-admin/global"
	"go-admin/model/dto"
	"go-admin/model/system"
	"go-admin/model/vo"
	"gorm.io/gorm"
)

type AuthorityBtnService struct{}

func (abs *AuthorityBtnService) GetAuthorityBtn(req dto.SysAuthorityBtnReq) (res vo.SysAuthorityBtnRes, err error) {
	var authorityBtn []system.SysAuthorityBtn
	err = global.GA_DB.Find(&authorityBtn, "authority_id = ? and sys_menu_id = ?", req.AuthorityId, req.MenuID).Error
	if err != nil {
		return
	}
	var selected []uint
	for _, v := range authorityBtn {
		selected = append(selected, v.SysBaseMenuBtnID)
	}
	res.Selected = selected
	return res, err
}

func (abs *AuthorityBtnService) SetAuthorityBtn(req dto.SysAuthorityBtnReq) (err error) {
	return global.GA_DB.Transaction(func(tx *gorm.DB) error {
		var authorityBtn []system.SysAuthorityBtn
		err = tx.Delete(&[]system.SysAuthorityBtn{}, "authority_id = ? and sys_menu_id = ?", req.AuthorityId, req.MenuID).Error
		if err != nil {
			return err
		}
		for _, v := range req.Selected {
			authorityBtn = append(authorityBtn, system.SysAuthorityBtn{
				AuthorityId:      req.AuthorityId,
				SysMenuID:        req.MenuID,
				SysBaseMenuBtnID: v,
			})
		}
		if len(authorityBtn) > 0 {
			err = tx.Create(&authorityBtn).Error
		}
		if err != nil {
			return err
		}
		return err
	})
}

func (abs *AuthorityBtnService) CanRemoveAuthorityBtn(id string) (err error) {
	fErr := global.GA_DB.First(&system.SysAuthorityBtn{}, "sys_base_menu_btn_id = ?", id).Error
	if errors.Is(fErr, gorm.ErrRecordNotFound) {
		return nil
	}
	return errors.New("此按钮正在被使用无法删除")
}
