package utils

import (
	"encoding/json"
	"fmt"
	"github.com/hanyougame/glib/utils/codec"
	"github.com/samber/lo"
	"github.com/samber/lo/mutable"
	"log"
	"strconv"
)

// PrettyJSON 美化打印
func PrettyJSON(v interface{}) {
	// 使用 json.MarshalIndent 进行格式化和美化打印
	prettyJSON, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatalf("JSON marshalling failed: %s", err)
	}
	// 打印格式化后的 JSON 字符串
	fmt.Println(string(prettyJSON))
}

// Ternary 三元运算符的模拟函数
func Ternary[T any](condition bool, value1, value2 T) T {
	if condition {
		return value1
	}
	return value2
}

func GenerateKeyAndIv(fileName string) (string, string) {
	hash := []rune(codec.SHA256.Encode(fileName))
	mutable.Reverse(hash)
	hashStr := string(hash)
	str := lo.Substring(hashStr, 0, 8)

	index, _ := strconv.ParseInt(str, 16, 64)
	index = index % 2

	var key, iv string
	if index != 0 {
		key = lo.Substring(hashStr, 6, 32)
		iv = lo.Substring(hashStr, 6, 16)
	} else {
		key = lo.Substring(hashStr, 10, 32)
		iv = lo.Substring(hashStr, 16, 16)
	}
	keyRuneSlice := []rune(key)
	mutable.Reverse(keyRuneSlice)
	return string(keyRuneSlice), iv
}
