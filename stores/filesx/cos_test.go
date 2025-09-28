package filesx

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestCosUpload(t *testing.T) {
	cos, _ := NewCosStorage(CosStorageConfig{
		BucketURL: "https://k1-cos-1333271963.cos.ap-hongkong.myqcloud.com",
		SecretKey: "",
		SecretID:  "",
	})
	var ctx = context.Background()
	f, _ := os.Open("1.html")
	_, err := cos.Upload(ctx, f, "upload/1.html")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestCosDelete(t *testing.T) {
	cos, _ := NewCosStorage(CosStorageConfig{
		BucketURL: "https://k1-cos-1333271963.cos.ap-hongkong.myqcloud.com",
		SecretKey: "",
		SecretID:  "",
	})

	var ctx = context.Background()
	var file = "upload/1.html"
	err := cos.Delete(ctx, file)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestGetPkgUploadSignUrl(t *testing.T) {
	url, err := GetPkgUploadSignUrl(CosStorageConfig{
		PkgBucketURL: "",
		SecretID:     "",
		SecretKey:    "",
	}, "app/k1_0.7.42.apk")
	if err != nil {
		fmt.Println("GetPkgUploadSignUrl err: ", err)
		return
	}
	fmt.Println("url: ", url)
}
