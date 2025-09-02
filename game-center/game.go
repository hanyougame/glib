package game_center

import (
	"context"
	"errors"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/logx"
)

type GameCenter struct {
	Config GameCenterConfig
}

func NewGameCenter(config GameCenterConfig) *GameCenter {
	return &GameCenter{
		config,
	}
}

// Register 注册账号
func (g *GameCenter) Register(ctx context.Context, req RegisterAccountReq) (openid string, err error) {
	httpResp, err := postRequest(ctx, g.Config, RegisterURL, req)
	if err != nil {
		logx.WithContext(ctx).Errorf("[game-center] register request error: %s", err.Error())
		return
	}

	registerResp := &CommonResp[RegisterAccountResp]{}
	err = jsonx.Unmarshal(httpResp.Body(), registerResp)
	if err != nil {
		logx.WithContext(ctx).Errorf("[game-center] unmarshal register response error: %s", err.Error())
		return
	}

	if registerResp.Code != 0 {
		logx.WithContext(ctx).Errorf("[game-center] register request error: %s", registerResp.Message)
		err = errors.New(registerResp.Message)
		return
	}

	return registerResp.Data.OpenID, nil
}

// GetGameLink 获取游戏地址
func (g *GameCenter) GetGameLink(ctx context.Context, req GetGameLinkReq) (resp GetGameLinkResp, err error) {
	var httpResp *resty.Response
	httpResp, err = postRequest(ctx, g.Config, GetGameLinkURL, req)
	if err != nil {
		logx.WithContext(ctx).Errorf("[game-center] get game link request error: %s", err.Error())
		return
	}

	getGameLinkResp := &CommonResp[GetGameLinkResp]{}
	err = jsonx.Unmarshal(httpResp.Body(), getGameLinkResp)
	if err != nil {
		logx.WithContext(ctx).Errorf("[game-center] unmarshal get game link response error: %s", err.Error())
		return
	}

	if getGameLinkResp.Code != 0 {
		logx.WithContext(ctx).Errorf("[game-center] get game link request error: %s", getGameLinkResp.Message)
		err = errors.New(getGameLinkResp.Message)
		return
	}

	resp = getGameLinkResp.Data
	return
}

// GetGameList 获取游戏列表
func (g *GameCenter) GetGameList(ctx context.Context, req GetGameListReq) (resp GetGameListResp, err error) {
	var httpResp *resty.Response
	httpResp, err = getRequest(ctx, g.Config, GetGameListURL, url.Values{"platform_id": {cast.ToString(req.PlatformID)}})
	if err != nil {
		logx.WithContext(ctx).Errorf("[game-center] get game list request error: %s", err.Error())
		return
	}
	getGameListResp := &CommonResp[GetGameListResp]{}
	err = jsonx.Unmarshal(httpResp.Body(), getGameListResp)
	if err != nil {
		logx.WithContext(ctx).Errorf("[game-center] unmarshal get game list response error: %s", err.Error())
		return
	}
	if getGameListResp.Code != 0 {
		logx.WithContext(ctx).Errorf("[game-center] get game list request error: %s", getGameListResp.Message)
		err = errors.New(getGameListResp.Message)
		return
	}
	resp = getGameListResp.Data
	return
}

// GetPlatformList 获取厂商列表
func (g *GameCenter) GetPlatformList(ctx context.Context) (resp GetGamePlatformListResp, err error) {
	var httpResp *resty.Response
	httpResp, err = getRequest(ctx, g.Config, GetPlatformList, nil)
	if err != nil {
		logx.WithContext(ctx).Errorf("[game-center] get platform list request error: %s", err.Error())
		return
	}
	getGamePlatformListResp := &CommonResp[GetGamePlatformListResp]{}
	err = jsonx.Unmarshal(httpResp.Body(), getGamePlatformListResp)
	if err != nil {
		logx.WithContext(ctx).Errorf("[game-center] unmarshal get platform list response error: %s", err.Error())
		return
	}
	if getGamePlatformListResp.Code != 0 {
		logx.WithContext(ctx).Errorf("[game-center] get platform list request error: %s", getGamePlatformListResp.Message)
		err = errors.New(getGamePlatformListResp.Message)
		return
	}
	resp = getGamePlatformListResp.Data
	return
}

// GetBetList 获取投注记录
func (g *GameCenter) GetBetList(ctx context.Context, req GetBetListReq) (resp GetBetListResp, err error) {
	// 结构体转url.Values
	values := url.Values{}
	values.Add("start_time", cast.ToString(req.StartTime))
	values.Add("end_time", cast.ToString(req.EndTime))
	if req.OpenID != 0 {
		values.Add("open_id", cast.ToString(req.OpenID))
	}
	values.Add("page", cast.ToString(req.Page))
	values.Add("page_size", cast.ToString(req.PageSize))
	values.Add("merchant_user_id", req.MerchantUserID)
	var httpResp *resty.Response
	httpResp, err = getRequest(ctx, g.Config, GetBetList, values)
	if err != nil {
		logx.WithContext(ctx).Errorf("[game-center] get bet list request error: %s", err.Error())
		return
	}
	getGameBetListResp := &CommonResp[GetBetListResp]{}
	err = jsonx.Unmarshal(httpResp.Body(), getGameBetListResp)
	if err != nil {
		logx.WithContext(ctx).Errorf("[game-center] get bet list request error: %s", err.Error())
		return
	}
	if getGameBetListResp.Code != 0 {
		logx.WithContext(ctx).Errorf("[game-center] get bet list request error: %s", getGameBetListResp.Message)
		err = errors.New(getGameBetListResp.Message)
		return
	}
	resp = getGameBetListResp.Data
	return
}
