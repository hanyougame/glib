package config

// Config redis配置
type Config struct {
	Addrs      []string `json:"addrs"`
	Debug      bool     `json:"debug,default=false"`
	Trace      bool     `json:"trace,default=false"`
	MasterName string   `json:"master_name,optional"`
	Username   string   `json:"username,optional"`
	Password   string   `json:"password,optional"`
	DB         int      `json:"db,default=0"`
}
