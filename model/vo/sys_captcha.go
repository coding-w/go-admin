package vo

// SysCaptchaResponse 验证码响应
type SysCaptchaResponse struct {
	CaptchaId     string `json:"captchaId"`
	Captcha       string `json:"captcha"`
	PicPath       string `json:"picPath"`
	CaptchaLength int    `json:"captchaLength"`
	OpenCaptcha   bool   `json:"openCaptcha"`
}
