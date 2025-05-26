package firebasex

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"github.com/zeromicro/go-zero/core/jsonx"
	"google.golang.org/api/option"
)

func NewClient(ctx context.Context, c Config) (*messaging.Client, error) {
	var (
		err error
		b   []byte
		app *firebase.App
	)
	if b, err = jsonx.Marshal(c); err != nil {
		return nil, fmt.Errorf("encode config err: %v", err)
	}
	if app, err = firebase.NewApp(ctx, nil, option.WithCredentialsJSON(b)); err != nil {
		return nil, fmt.Errorf("new firebase app err: %+v", err)
	}
	var client *messaging.Client
	if client, err = app.Messaging(ctx); err != nil {
		return nil, fmt.Errorf("new firebase messaging client err: %+v", err)
	}

	return client, nil
}
