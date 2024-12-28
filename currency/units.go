package currency

// CurrencyUnit 货币单位类型
type CurrencyUnit int64

// Unit 货币单位
type Unit int

const (
	UnitYuan Unit = iota // 元
	UnitJiao             // 角
	UnitFen              // 分
	UnitLi               // 厘
	UnitMao              // 毫
	UnitSi               // 丝
)
