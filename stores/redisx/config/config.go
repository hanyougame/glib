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
	IsCluster  bool     `json:"is_cluster,optional"`

	// 连接池配置。默认值面向高并发消费场景(mq-server 类服务)，可按需在配置中覆盖。
	// 背景：引入 automaxprocs 后，容器内 GOMAXPROCS 被限制为 cgroup CPU 配额，
	// go-redis 默认连接池(10*GOMAXPROCS)偏小，高并发下易打满排队。

	// PoolSize 连接池最大连接数，默认 100。
	PoolSize int `json:"pool_size,default=100"`
	// MinIdleConns 最小空闲连接数，用于预热连接，避免突发流量时集中建连，默认 20。
	MinIdleConns int `json:"min_idle_conns,default=20"`
	// DialTimeoutMs 拨号(建连)超时，单位毫秒，默认 3000(3s)。
	DialTimeoutMs int `json:"dial_timeout_ms,default=3000"`
	// ReadTimeoutMs 读超时，单位毫秒，默认 3000(3s)。
	ReadTimeoutMs int `json:"read_timeout_ms,default=3000"`
	// WriteTimeoutMs 写超时，单位毫秒，默认 3000(3s)。
	WriteTimeoutMs int `json:"write_timeout_ms,default=3000"`
	// PoolTimeoutMs 从连接池获取连接的等待超时，单位毫秒，默认 3000(3s)。
	// 连接池打满时请求在该时间内排队等待空闲连接，超时返回 redis: connection pool timeout。
	PoolTimeoutMs int `json:"pool_timeout_ms,default=3000"`
}
