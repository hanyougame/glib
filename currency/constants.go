package currency

import (
	"fmt"
	"math"
)

const (
	OneUnit = 1000000 // 1元 = 1,000,000 微元（精度6位小数）

	Yuan CurrencyUnit = OneUnit * 1      // 1元 = 1,000,000 微元
	Jiao CurrencyUnit = OneUnit / 10     // 1角 = 100,000 微元
	Fen  CurrencyUnit = OneUnit / 100    // 1分 = 10,000 微元
	Li   CurrencyUnit = OneUnit / 1000   // 1厘 = 1,000 微元
	Mao  CurrencyUnit = OneUnit / 10000  // 1毫 = 100 微元
	Si   CurrencyUnit = OneUnit / 100000 // 1丝 = 10 微元
)

// ToYuanString 转换为元字符串
func (c CurrencyUnit) ToYuanString() string {
	yuan := float64(c) / float64(OneUnit)
	return fmt.Sprintf("%.6f 元", yuan)
}

// Float64 转换为浮点数
func (c CurrencyUnit) Float64() float64 {
	return float64(c) / float64(OneUnit)
}

// ConvertTo 转换为指定单位
func (c CurrencyUnit) ConvertTo(target Unit) int64 {
	switch target {
	case UnitYuan:
		return int64(math.Round(float64(c) / float64(OneUnit)))
	case UnitJiao:
		return int64(math.Round(float64(c) / float64(Jiao)))
	case UnitFen:
		return int64(math.Round(float64(c) / float64(Fen)))
	case UnitLi:
		return int64(math.Round(float64(c) / float64(Li)))
	case UnitMao:
		return int64(math.Round(float64(c) / float64(Mao)))
	case UnitSi:
		return int64(math.Round(float64(c) / float64(Si)))
	default:
		return 0
	}
}

// FromFloat64 从浮点数创建货币单位
func FromFloat64(value float64) CurrencyUnit {
	return CurrencyUnit(value * float64(OneUnit))
}
