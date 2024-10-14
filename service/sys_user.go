package service

import (
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"go-admin/global"
	"go-admin/model/common/request"
	"go-admin/model/system"
	"go-admin/utils"
	"gorm.io/gorm"
	"time"
)

type UserService struct {
}

func (us *UserService) Login(u *system.SysUser) (*system.SysUser, error) {
	if nil == global.GA_DB {
		return nil, fmt.Errorf("db not init")
	}

	var user system.SysUser
	err := global.GA_DB.Where("username = ?", u.Username).Preload("Authorities").Preload("Authority").First(&user).Error
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return nil, errors.New("密码错误")
		}
		menuServiceApp.UserAuthorityDefaultRouter(&user)
	}
	return &user, err
}

// Register
// 注册用户
func (us *UserService) Register(u system.SysUser) (userInter system.SysUser, err error) {
	var user system.SysUser
	err = global.GA_DB.
		Model(&system.SysUser{}).
		Where("username = ?", u.Username).
		First(&user).Error
	// 判断用户名是否被注册
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return userInter, errors.New("用户名已被注册")
	}
	// 密码加密
	u.Password = utils.BcryptHash(u.Password)
	// uuid
	u.UUID = uuid.Must(uuid.NewV4())
	err = global.GA_DB.
		Model(&system.SysUser{}).
		Create(&u).Error
	return u, err
}

// ChangePassword
// 设置密码
func (us *UserService) ChangePassword(u *system.SysUser, newPassword string) (userInter *system.SysUser, err error) {
	var user system.SysUser
	if err = global.GA_DB.Where("id = ?", u.ID).First(&user).Error; err != nil {
		return nil, err
	}
	if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
		return nil, errors.New("原密码错误")
	}
	user.Password = utils.BcryptHash(newPassword)
	err = global.GA_DB.Model(&user).Update("password", user.Password).Error
	return &user, err
}

// GetUserInfoList
// 分页获取用户列表
func (us *UserService) GetUserInfoList(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GA_DB.Model(&system.SysUser{})
	var userList []system.SysUser
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Preload("Authorities").Preload("Authority").Find(&userList).Error
	return userList, total, err
}

// SetUserAuthority
// 设置用户权限
func (us *UserService) SetUserAuthority(id uint, authorityId uint) error {
	assignErr := global.GA_DB.Where("sys_user_id = ? AND sys_authority_authority_id = ?", id, authorityId).First(&system.SysUserAuthority{}).Error
	if errors.Is(assignErr, gorm.ErrRecordNotFound) {
		return errors.New("该用户无此角色")
	}
	return global.GA_DB.Model(&system.SysUser{}).Where("id = ?", id).Update("authority_id", authorityId).Error
}

// SetUserAuthorities
// 设置一个用户的权限
func (us *UserService) SetUserAuthorities(id uint, authorityIds []uint) error {
	return global.GA_DB.Transaction(func(tx *gorm.DB) error {
		var user system.SysUser
		err := tx.Where("id = ?", id).First(&user).Error
		if err != nil {
			global.GA_LOG.Debug(err.Error())
			return errors.New("查询用户信息失败")
		}
		// 删除当前用户的所有权限记录
		err = tx.Delete(&system.SysUserAuthority{}, "sys_user_id = ?", id).Error
		if err != nil {
			global.GA_LOG.Debug(err.Error())
			return errors.New("删除用户权限失败")
		}
		// 创建一个空的SysUserAuthority切片，用于批量插入新的权限记录
		var userAuth []system.SysUserAuthority
		for _, v := range authorityIds {
			userAuth = append(userAuth, system.SysUserAuthority{
				SysUserId:               id,
				SysAuthorityAuthorityId: v,
			})
		}
		// 将新的权限记录批量插入数据库
		err = tx.Create(&userAuth).Error
		if err != nil {
			global.GA_LOG.Debug(err.Error())
			return errors.New("添加用户权限失败")
		}
		// 更新用户表中的authority_id字段为新的权限ID中的第一个
		err = tx.Model(&user).Update("authority_id", authorityIds[0]).Error
		if err != nil {
			global.GA_LOG.Debug(err.Error())
			return errors.New("更新用户主权限失败")
		}
		// 返回 nil 提交事务
		return nil
	})
}

// DeleteUser
// 删除用户
func (us *UserService) DeleteUser(id int) error {
	return global.GA_DB.Transaction(func(tx *gorm.DB) error {
		// 删除用户
		if err := tx.Where("id = ?", id).Delete(&system.SysUser{}).Error; err != nil {
			return err
		}
		// 删除用户 权限
		if err := tx.Delete(&[]system.SysUserAuthority{}, "sys_user_id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}

// SetUserInfo
// 设置用户信息
func (us *UserService) SetUserInfo(req system.SysUser) error {
	return global.GA_DB.Model(&system.SysUser{}).
		Select("updated_at", "nick_name", "header_img", "phone", "email", "sideMode", "enable").
		Where("id=?", req.ID).
		Updates(map[string]interface{}{
			"updated_at": time.Now(),
			"nick_name":  req.NickName,
			"header_img": req.HeaderImg,
			"phone":      req.Phone,
			"email":      req.Email,
			"side_mode":  req.SideMode,
			"enable":     req.Enable,
		}).Error
}

func (us *UserService) SetSelfInfo(req system.SysUser) error {
	return global.GA_DB.Model(&system.SysUser{}).
		Where("id=?", req.ID).
		Updates(req).Error
}

func (us *UserService) GetUserInfo(uuid uuid.UUID) (system.SysUser, error) {
	var reqUser system.SysUser
	err := global.GA_DB.Preload("Authorities").Preload("Authority").First(&reqUser, "uuid = ?", uuid).Error
	if err != nil {
		return reqUser, err
	}
	menuServiceApp.UserAuthorityDefaultRouter(&reqUser)
	return reqUser, err
}

func (us *UserService) ResetPassword(ID uint) (err error) {
	err = global.GA_DB.Model(&system.SysUser{}).Where("id = ?", ID).Update("password", utils.BcryptHash("123456")).Error
	return err
}
