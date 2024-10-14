package service

import (
	"errors"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go-admin/global"
	"go-admin/model/dto"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"sync"
)

type CasbinService struct{}

var casbinServiceApp = new(CasbinService)

var (
	syncedCachedEnforcer *casbin.SyncedCachedEnforcer
	once                 sync.Once
)

// Casbin 创建一个管理访问控制策略对象
func (s *CasbinService) Casbin() *casbin.SyncedCachedEnforcer {
	// 单例模式
	once.Do(func() {
		db, err := gormadapter.NewAdapterByDB(global.GA_DB)
		if err != nil {
			global.GA_LOG.Error("适配数据库失败请检查casbin表是否为InnoDB引擎!", zap.Error(err))
			return
		}
		// request_definition 请求定义
		// r 表示请求的格式，包含三个元素：sub (subject, 主体)、obj (object, 对象) 和 act (action, 动作)

		// policy_definition 策略定义
		// p 表示策略的格式，包含三个元素：sub (subject, 主体)、obj (object, 对象) 和 act (action, 动作)

		// role_definition 角色定义
		// g 表示角色的格式，包含两个元素：_ (表示泛指主体)

		// policy_effect 策略效果
		// e 表示策略的效果，定义了当某个策略匹配时的结果。在这里，表示如果策略中有某条规则的效果 p.eft 为 allow，则允许该请求

		// matchers 匹配器
		// m 定义了如何匹配请求和策略。
		// 在这里，表示请求的主体 r.sub 必须与策略中的主体 p.sub 相同，且请求的对象 r.obj 必须与策略中的对象 p.obj 匹配（使用 keyMatch2 函数），
		// 并且请求的动作 r.act 必须与策略中的动作 p.act 相同
		text := `
		[request_definition]
		r = sub, obj, act

		[policy_definition]
		p = sub, obj, act

		[role_definition]
		g = _, _

		[policy_effect]
		e = some(where (p.eft == allow))

		[matchers]
		m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
		`
		// 从字符串 text 中创建 Casbin 的模型
		m, err := model.NewModelFromString(text)
		if err != nil {
			global.GA_LOG.Error("字符串加载模型失败!", zap.Error(err))
			return
		}
		// 创建一个 casbin.SyncedCachedEnforcer 对象，将模型 m 和数据库适配器 db 传入；这个对象支持策略的缓存和同步加载
		syncedCachedEnforcer, _ = casbin.NewSyncedCachedEnforcer(m, db)
		// 设置缓存的过期时间为 1 小时
		syncedCachedEnforcer.SetExpireTime(60 * 60)
		// 加载策略到 syncedCachedEnforcer 对象中
		_ = syncedCachedEnforcer.LoadPolicy()
	})
	return syncedCachedEnforcer
}

// AddPolicies
// 添加匹配的权限
func (s *CasbinService) AddPolicies(db *gorm.DB, rules [][]string) error {
	var casbinRules []gormadapter.CasbinRule
	for _, rule := range rules {
		casbinRules = append(casbinRules, gormadapter.CasbinRule{
			Ptype: "p",
			V0:    rule[0],
			V1:    rule[1],
			V2:    rule[2],
		})
	}
	return db.Create(&casbinRules).Error
}

// FreshCasbin
// 刷新权限
func (s *CasbinService) FreshCasbin() error {
	e := s.Casbin()
	err := e.LoadPolicy()
	return err
}

// GetPolicyPathByAuthorityId 根据authId获取 策略 path
func (s *CasbinService) GetPolicyPathByAuthorityId(AuthorityID uint) (pathMaps []dto.CasbinInfo) {
	e := s.Casbin()
	authorityId := strconv.Itoa(int(AuthorityID))
	list := e.GetFilteredPolicy(0, authorityId)
	for _, v := range list {
		pathMaps = append(pathMaps, dto.CasbinInfo{
			Path:   v[1],
			Method: v[2],
		})
	}
	return pathMaps
}

// UpdateCasbin 更新策略
func (s *CasbinService) UpdateCasbin(AuthorityID uint, casbinInfos []dto.CasbinInfo) error {
	authorityId := strconv.Itoa(int(AuthorityID))
	s.ClearCasbin(0, authorityId)
	rules := [][]string{}
	// 做权限去重处理
	deduplicateMap := make(map[string]bool)
	for _, v := range casbinInfos {
		key := authorityId + v.Path + v.Method
		if _, ok := deduplicateMap[key]; !ok {
			deduplicateMap[key] = true
			rules = append(rules, []string{authorityId, v.Path, v.Method})
		}
	}
	e := s.Casbin()
	success, _ := e.AddPolicies(rules)
	if !success {
		return errors.New("存在相同api,添加失败,请联系管理员")
	}
	return nil
}

// ClearCasbin 移除 Casbin 中符合条件的策略
func (s *CasbinService) ClearCasbin(v int, p ...string) bool {
	e := s.Casbin()
	success, _ := e.RemoveFilteredPolicy(v, p...)
	return success
}

// RemoveFilteredPolicy 删除策略
func (s *CasbinService) RemoveFilteredPolicy(db *gorm.DB, authorityId string) error {
	return db.Delete(&gormadapter.CasbinRule{}, "v0 = ?", authorityId).Error
}

// UpdateCasbinApi 更新策略
func (s *CasbinService) UpdateCasbinApi(oldPath string, newPath string, oldMethod string, newMethod string) error {
	err := global.GA_DB.Model(&gormadapter.CasbinRule{}).Where("v1 = ? AND v2 = ?", oldPath, oldMethod).Updates(map[string]interface{}{
		"v1": newPath,
		"v2": newMethod,
	}).Error
	e := s.Casbin()
	err = e.LoadPolicy()
	if err != nil {
		return err
	}
	return err
}
