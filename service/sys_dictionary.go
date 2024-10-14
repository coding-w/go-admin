package service

import (
	"errors"
	"go-admin/global"
	"go-admin/model/system"
	"gorm.io/gorm"
)

type DictionaryService struct{}

func (s *DictionaryService) CreateSysDictionary(sysDictionary system.SysDictionary) error {
	if (!errors.Is(global.GA_DB.First(&system.SysDictionary{}, "type = ?", sysDictionary.Type).Error, gorm.ErrRecordNotFound)) {
		return errors.New("存在相同的type，不允许创建")
	}
	err := global.GA_DB.Create(&sysDictionary).Error
	return err
}

func (s *DictionaryService) DeleteSysDictionary(sysDictionary system.SysDictionary) (err error) {
	err = global.GA_DB.Where("id = ?", sysDictionary.ID).Preload("SysDictionaryDetails").First(&sysDictionary).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("请不要搞事")
	}
	if err != nil {
		return err
	}
	err = global.GA_DB.Delete(&sysDictionary).Error
	if err != nil {
		return err
	}

	if sysDictionary.SysDictionaryDetails != nil {
		return global.GA_DB.Where("sys_dictionary_id=?", sysDictionary.ID).Delete(sysDictionary.SysDictionaryDetails).Error
	}
	return
}

func (s *DictionaryService) UpdateSysDictionary(sysDictionary *system.SysDictionary) (err error) {
	var dict system.SysDictionary
	sysDictionaryMap := map[string]interface{}{
		"Name":   sysDictionary.Name,
		"Type":   sysDictionary.Type,
		"Status": sysDictionary.Status,
		"Desc":   sysDictionary.Desc,
	}
	err = global.GA_DB.Where("id = ?", sysDictionary.ID).First(&dict).Error
	if err != nil {
		global.GA_LOG.Debug(err.Error())
		return errors.New("查询字典数据失败")
	}
	if dict.Type != sysDictionary.Type {
		if !errors.Is(global.GA_DB.First(&system.SysDictionary{}, "type = ?", sysDictionary.Type).Error, gorm.ErrRecordNotFound) {
			return errors.New("存在相同的type，不允许创建")
		}
	}
	err = global.GA_DB.Model(&dict).Updates(sysDictionaryMap).Error
	return err
}

func (s *DictionaryService) GetSysDictionary(Type string, Id uint, status *bool) (sysDictionary system.SysDictionary, err error) {
	var flag = false
	if status == nil {
		flag = true
	} else {
		flag = *status
	}
	err = global.GA_DB.Where("(type = ? OR id = ?) and status = ?", Type, Id, flag).Preload("SysDictionaryDetails", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", true).Order("sort")
	}).First(&sysDictionary).Error
	return
}

func (s *DictionaryService) GetSysDictionaryInfoList() (list interface{}, err error) {
	var sysDictionarys []system.SysDictionary
	err = global.GA_DB.Find(&sysDictionarys).Error
	return sysDictionarys, err
}
