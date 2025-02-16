package currency

import (
	"fmt"
	"testing"
)

func TestCurrencyUnit_ConvertTo(t *testing.T) {

	var c int64 = 12345678
	fmt.Println(WeiToYuan(c))
}
