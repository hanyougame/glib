package config

type Mode int

const (
	Mysql Mode = iota + 1
	Postgres
	ClickHouse
	PostgresRead
)

type Config struct {
	Mode                   Mode     `json:"mode"`
	Separation             bool     `json:",default=false"`
	Trace                  bool     `json:"trace,default=false"`
	Master                 string   `json:"master,optional"`
	Sources                []string `json:"sources,optional"`
	Replicas               []string `json:"replicas,optional"`
	DSN                    string   `json:"dsn,optional"`
	Debug                  bool     `json:"debug,default=false"`
	MaxIdleConn            int      `json:"max_idle_conn"`
	MaxOpenConn            int      `json:"max_open_conn"`
	MaxLifetime            int      `json:"max_lifetime"`
	PrepareStmt            bool     `json:"prepare_stmt"`
	SkipDefaultTransaction bool     `json:"skip_default_transaction"`
}
