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
	const creat = `{
  "type": "service_account",
  "project_id": "k1-test-66ec7",
  "private_key_id": "4414d6c1c99a8c2df3a51e773bc87822a183158c",
  "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCsN9U9QIXRYczC\n4r5z+m2+i9kzqzBXje9G0y7PwW83f90MQ6JWgnchYepcixWR+uethoYzThoMdK/k\nI2YoMw1XJe6LbfwqrygYObdFL9pkPBvYKWOg1vi7750/sNiGRPkJm/MyDeKk0Q1m\nA6cjgASNdKi7knfraZCBzttFXmmUuL8qhUcZz86M/R/XMFaZ0Wl1ipaDhLqgK73F\nFAgFLjioc0DGJhck966I4lkAsDseh0GS3tzmp6xM1TZ6e/spLU2qGqg6DGw/tJTr\n5LG+U55xcQzf5m2BSa2ri5JNI6vJ/t4PrEtfPR3mdEPc+Issk7VZ/xEP/GU/BibM\nBNhV5ImnAgMBAAECggEABsxRa5CPd7tvlMXUlMlu6ZWlpkd8qggdSONvNp5BIKS4\nehJ1ZwrRXwz5blvDJzcLTCtCQuYzWwiEJPJMxwt03FrP9ssJqDmv0Hz+wgloiopg\n4wCHQUMbIHgN8gSuF2bDmeCBEhnAN5dU+xDr2RfllhvjyTu35iANp0IeH7fhXRB5\nMabkO3p3XENSUyLrbwRH2+sFYbcWZ/Xi63vrPtrpqIz17T+qO6L39SokDecQiy4R\nvcSUumZ84h9VSd7xLsTVy6+ldnKJsrwvsGZLbKNOrmn4HOTSRCL+i03Zk950lXQk\n3fU/wyQX/tGFvTSZPrdVRhEvH2hi8rshXhDse/eKeQKBgQDsdbbIeRM+uNysbkbS\nm9VvHPKzSr7E7Y9U3h/WS1GE/UhZII9xxVEa57aGAeeSWVfj2I5KqFbcGRRDw9Bn\niiwc0oalqu6xRYm9E9wN4Um/QImDUabHWSCkeSs4keDyGDXVHzAryjS+5de2i2Jw\nazKIM1rDxJ2mpLgVIPDpics6BQKBgQC6cxdWg5MEIQwbSxolTBbOr8kjXx8007t+\ncB+QWcKGy6m4Od2MzXjmJV35MCa4+ZI7NUkZt2q2NxjH2GQ5TXZ/Q0xAfxXMqdOF\njZ7GNTw1gi6igZa2+81VwmrzVuppWdcOEzUAX1BMgZlt4KTN7yOgZTxhBvDeKH9W\n6xLjhu0IuwKBgAZPeepNuOLCojkD1GYIfKCzgBKi2wZ8ZIfaSbU5W1qWO9kNEmbx\n4iE6r9dRs7FnKv5MqDd72J2VGcJZNnWc7WJzh95h+m7GsU9XeVbxfTtaaJEVm8LU\nMRfSmAGxevRbCwy+AVYZ1mCF18TqYV0orzeNT03MMHzjBnqh32wMHl1VAoGAdjgX\naiGPDCUwGfYqo890/qRy8lyT2tkNnwPU36AqOXHVMRxvn6GAsNyskx22z75eu+/z\nY7zqipKwwbYxBixZoCO5q4No/RYJKjs7mw2SGZNrolpsh6Vs8p2NVrGbfCUcuK3x\nA0VsErZx5Dt2A6VHf6HxAMx32UTwSNpQZbfdf5UCgYEAg71djxdTL2bBZ1gkaI7e\n2Jl2rThAYcKxdwRt9Vaxg/LQk0XCN13IXdLzQBFqvhJjvoxjDI/bX/UIwxdmYzyB\nidvRVD9lul1WtjaVZh0WwL1Oz7vXBkA6lx1ZuG+A46iL6mo9a0a8zFW+uADX8VxW\not8geGiOks6YroBbD+6Hy2o=\n-----END PRIVATE KEY-----\n",
  "client_email": "firebase-adminsdk-fbsvc@k1-test-66ec7.iam.gserviceaccount.com",
  "client_id": "106184005349660568843",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/firebase-adminsdk-fbsvc%40k1-test-66ec7.iam.gserviceaccount.com",
  "universe_domain": "googleapis.com"
}`

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
