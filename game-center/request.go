package game_center

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/hanyougame/glib/utils/httpc"
)

func postRequest(ctx context.Context, config GameCenterConfig, url string, body any) (resp *resty.Response, err error) {
	resp, err = httpc.Do(ctx).
		SetBasicAuth(config.Username, config.Password).
		SetBody(body).
		Post(config.RequestURL + url)
	if err != nil {
		return
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(resp.String())
	}
	return resp, nil
}

func getRequest(ctx context.Context, config GameCenterConfig, url string, params url.Values) (resp *resty.Response, err error) {
	return httpc.Do(ctx).
		SetBasicAuth(config.Username, config.Password).
		SetQueryParamsFromValues(params).
		Get(config.RequestURL + url)
}
