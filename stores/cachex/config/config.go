package config

// Config 缓存配置
type Config struct {
	StrongConsistency bool `json:"strong_consistency,default=true"`
	DisableCacheRead  bool `json:"disable_cache_read,default=false"`
}
