package currency

import (
	"fmt"
	"testing"
)

func TestCurrencyUnit_ConvertTo(t *testing.T) {
	var c Unit = 12345678
	fmt.Println(c.ConvertYuan())
}
