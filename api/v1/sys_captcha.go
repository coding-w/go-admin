package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go-admin/global"
	"go-admin/model/common/response"
	"go-admin/model/vo"
	"go.uber.org/zap"
	"time"
)

// 生成图片
var store = base64Captcha.DefaultMemStore

type CaptchaApi struct{}

func (ca *CaptchaApi) Captcha(c *gin.Context) {
	// 判断验证码是否开启
	openCaptcha := global.GA_CONFIG.Captcha.OpenCaptcha               // 是否开启防爆次数
	openCaptchaTimeOut := global.GA_CONFIG.Captcha.OpenCaptchaTimeOut // 缓存超时时间
	key := c.ClientIP()
	v, ok := global.LocalCache.Get(key)
	if ok {
		global.LocalCache.Set(key, 1, time.Second*time.Duration(openCaptchaTimeOut))
	}
	var oc bool
	if openCaptcha == 0 || openCaptcha < interfaceToInt(v) {
		oc = true
	}
	// 字符,公式,验证码配置
	// 生成默认数字的driver
	driver := base64Captcha.NewDriverDigit(
		global.GA_CONFIG.Captcha.ImgHeight,
		global.GA_CONFIG.Captcha.ImgWidth,
		global.GA_CONFIG.Captcha.KeyLong,
		0.7,
		80)
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, res, err := cp.Generate()
	if err != nil {
		global.GA_LOG.Error("验证码获取失败!", zap.Error(err))
		response.FailWithMessage("验证码获取失败", c)
		return
	}
	response.OkWithDetailed(vo.SysCaptchaResponse{
		CaptchaId:     id,
		PicPath:       b64s,
		CaptchaLength: global.GA_CONFIG.Captcha.KeyLong,
		OpenCaptcha:   oc,
		Captcha:       res,
	}, "验证码获取成功", c)
}

// interfaceToInt 类型转换
func interfaceToInt(v interface{}) (i int) {
	switch v := v.(type) {
	case int:
		i = v
	default:
		i = 0
	}
	return
}
