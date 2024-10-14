package service

import (
	"errors"
	"go-admin/global"
	"go-admin/model/system"
	"gorm.io/gorm"
)

type BaseMenuService struct{}

func (bms BaseMenuService) DeleteBaseMenu(id int) (err error) {
	err = global.GA_DB.First(&system.SysBaseMenu{}, "parent_id = ?", id).Error
	if err != nil {
		return global.GA_DB.Transaction(func(tx *gorm.DB) error {
			err = tx.Delete(&system.SysBaseMenu{}, "id = ?", id).Error
			if err != nil {
				return err
			}
			err = tx.Delete(&system.SysBaseMenuParameter{}, "sys_base_menu_id = ?", id).Error
			if err != nil {
				return err
			}
			err = tx.Delete(&system.SysBaseMenuBtn{}, "sys_base_menu_id = ?", id).Error
			if err != nil {
				return err
			}
			err = tx.Delete(&system.SysAuthorityBtn{}, "sys_menu_id = ?", id).Error
			if err != nil {
				return err
			}
			err = tx.Delete(&system.SysAuthorityMenu{}, "sys_base_menu_id = ?", id).Error
			if err != nil {
				return err
			}
			return nil
		})
	}
	return errors.New("此菜单存在子菜单不可删除")
}

// UpdateBaseMenu
// 更新路由
func (bms BaseMenuService) UpdateBaseMenu(menu system.SysBaseMenu) error {
	var oldMenu system.SysBaseMenu
	upDateMap := make(map[string]interface{})
	upDateMap["keep_alive"] = menu.KeepAlive
	upDateMap["close_tab"] = menu.CloseTab
	upDateMap["default_menu"] = menu.DefaultMenu
	upDateMap["parent_id"] = menu.ParentId
	upDateMap["path"] = menu.Path
	upDateMap["name"] = menu.Name
	upDateMap["hidden"] = menu.Hidden
	upDateMap["component"] = menu.Component
	upDateMap["title"] = menu.Title
	upDateMap["active_name"] = menu.ActiveName
	upDateMap["icon"] = menu.Icon
	upDateMap["sort"] = menu.Sort
	err := global.GA_DB.Transaction(func(tx *gorm.DB) error {
		tx.Where("id = ?", menu.ID).Find(&oldMenu)
		if oldMenu.Name != menu.Name {
			if !errors.Is(tx.Where("id <> ? AND name = ?", menu.ID, menu.Name).First(&system.SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
				global.GA_LOG.Debug("存在相同name修改失败")
				return errors.New("存在相同name修改失败")
			}
		}
		txErr := tx.Unscoped().Delete(&system.SysBaseMenuParameter{}, "sys_base_menu_id = ?", menu.ID).Error
		if txErr != nil {
			global.GA_LOG.Debug(txErr.Error())
			return txErr
		}
		txErr = tx.Unscoped().Delete(&system.SysBaseMenuBtn{}, "sys_base_menu_id = ?", menu.ID).Error
		if txErr != nil {
			global.GA_LOG.Debug(txErr.Error())
			return txErr
		}

		if len(menu.Parameters) > 0 {
			for k := range menu.Parameters {
				menu.Parameters[k].SysBaseMenuID = menu.ID
			}
			txErr = tx.Create(&menu.Parameters).Error
			if txErr != nil {
				global.GA_LOG.Debug(txErr.Error())
				return txErr
			}
		}

		if len(menu.MenuBtn) > 0 {
			for k := range menu.MenuBtn {
				menu.MenuBtn[k].SysBaseMenuID = menu.ID
			}
			txErr = tx.Create(&menu.MenuBtn).Error
			if txErr != nil {
				global.GA_LOG.Debug(txErr.Error())
				return txErr
			}
		}
		txErr = tx.Model(&oldMenu).Updates(upDateMap).Error
		if txErr != nil {
			global.GA_LOG.Debug(txErr.Error())
			return txErr
		}
		return nil
	})
	return err
}

func (bms BaseMenuService) GetBaseMenuById(id int) (menu system.SysBaseMenu, err error) {
	err = global.GA_DB.Preload("MenuBtn").Preload("Parameters").Where("id = ?", id).First(&menu).Error
	return
}
