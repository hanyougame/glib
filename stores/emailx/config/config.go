package config

// Config redis配置
type Config struct {
	Debug      bool   `json:"debug,default=false"`
	Trace      bool   `json:"trace,default=false"`
	SmtpServer string `json:"smtp_server,required"` // smtp服务器地址
	SmtpPort   int    `json:"smtp_port,required"`   // smtp服务器地址
	Password   string `json:"password,required"`    // 发送密码
	SendEmail  string `json:"send_email,optional"`  // 发送地址
	SendName   string `json:"send_name,optional"`   // 发送人名称
}
