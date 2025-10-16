package whatsappx

type Config struct {
	BaseURL           string
	APIVersion        string
	AccessToken       string
	PhoneNumberID     string
	BusinessAccountID string
	AppSecret         string
	AppID             string
	SecureRequests    bool

	VerificationCodeTmplName string   // 验证码模板名称
	VerificationCodeLanguage Language // 验证码语言
}

type Language int

const (
	_ Language = iota
	English
)

func (t Language) ToWhatsAppLang() string {
	switch t {
	case English:
		return "en_US"
	default:
		return "en_US"
	}
}
