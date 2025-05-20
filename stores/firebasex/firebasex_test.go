package firebasex

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"github.com/zeromicro/go-zero/core/jsonx"
	"testing"
	"time"
)

func TestSendBatch(t *testing.T) {
	const creat = ``
	var c Config
	if err := jsonx.UnmarshalFromString(creat, &c); err != nil {
		t.Fatal(err)
	}
	//fmt.Printf("c: %+v\n", c)
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	client, err := NewClient(ctx, c)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.SendEachForMulticast(ctx, &messaging.MulticastMessage{
		Tokens: []string{
			"token1",
			"token2",
		},
		Data: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
		Notification: &messaging.Notification{
			Title:    "this is title",
			Body:     "this is body",
			ImageURL: "https://k1-misc-dev.hanyouweb.com/siteadmin/default/img_hd_mn1.png",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("SuccessCount: %+v\n", resp.SuccessCount)
	fmt.Printf("FailureCount: %+v\n", resp.FailureCount)
	for _, r := range resp.Responses {
		fmt.Printf("response: %+v\n", r)
	}
}
