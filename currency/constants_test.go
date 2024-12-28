package currency

import (
	"testing"
)

func TestCurrencyUnit_ConvertTo(t *testing.T) {
	var amountInYuan CurrencyUnit = Yuan // 1元

	// 转换为其他单位并打印
	t.Logf("1元 = %d 角\n", amountInYuan.ConvertTo(UnitJiao))
	t.Logf("1元 = %d 分\n", amountInYuan.ConvertTo(UnitFen))
	t.Logf("1元 = %d 厘\n", amountInYuan.ConvertTo(UnitLi))
	t.Logf("1元 = %d 毫\n", amountInYuan.ConvertTo(UnitMao))
	t.Logf("1元 = %d 丝\n", amountInYuan.ConvertTo(UnitSi))

	// 将元转换为 CurrencyUnit 类型
	var floatAmount float64 = 123.456789
	currencyUnit := FromFloat64(floatAmount)
	t.Logf("%.6f 元转换为 CurrencyUnit: %s\n", floatAmount, currencyUnit.ToYuanString())

	// 转换回浮点数并打印
	t.Logf("再转换回浮点数: %.6f 元\n", currencyUnit.Float64())

	// 转换为其他单位并打印
	t.Log("\n123.456789元转换为其他单位:")
	t.Logf("123.456789元 = %d 角\n", currencyUnit.ConvertTo(UnitJiao))
	t.Logf("123.456789元 = %d 分\n", currencyUnit.ConvertTo(UnitFen))
	t.Logf("123.456789元 = %d 厘\n", currencyUnit.ConvertTo(UnitLi))
	t.Logf("123.456789元 = %d 毫\n", currencyUnit.ConvertTo(UnitMao))
	t.Logf("123.456789元 = %d 丝\n", currencyUnit.ConvertTo(UnitSi))
	t.Log(currencyUnit)

	// 从更小单位转换回元
	var amountInSi CurrencyUnit = Si * 1000000 // 1000000丝（等于100元）
	t.Log("\n1000000丝转换为元:")
	t.Logf("1000000丝 = %.6f 元\n", amountInSi.Float64())
}

func TestCurrencyUnit_ToYuanString(t *testing.T) {
	var amountInSi CurrencyUnit = Si * 1000000 // 1000000丝（等于100元）

	// 转换为元并打印
	t.Logf("\n1000000丝转换为元:")
	t.Logf("1000000丝 = %.6f 元\n", amountInSi.Float64())
}

func TestFromFloat64(t *testing.T) {
	t.Logf("%.6f 元转换为微元: %d 微元\n", 123.456789, FromFloat64(123.456789))
}

func TestFromYuan(t *testing.T) {
	amounts := []float64{
		1234.56789,
		123.456789,
		123.456,
		123.45,
		123.4,
		123.0,
		123,
		12.3,
		1.23,
		0.123,
		0.0123,
		0.00123,
		0.000123,
		0.0000123,
		0.00000123,
		0.000000123,
		0.0000000123,
		0.00000000123,
		0.000000000123,
		0.0000000000123,
	}

	for _, amount := range amounts {
		t.Logf("%.13f 元转换为微元: %d 微元\n", amount, FromFloat64(amount))
	}
}
