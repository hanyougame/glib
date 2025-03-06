package translation

import (
	"fmt"
	"testing"
)

func TestTranslation(t *testing.T) {
	resp, err := TranslateText("en", "你好")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(resp)
}
