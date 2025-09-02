package game_center

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hanyougame/glib/utils"
)

func TestGameCenter_Register(t *testing.T) {
	g := NewGameCenter(GameCenterConfig{
		Username:   "90fae8ef04e68c6bc8dfdce7b1b1b3bc",
		Password:   "84363c525fc8f8b9703065772a918eb8",
		RequestURL: "https://center-openapi-gateway-dev.ogweb.org",
	})

	openid, err := g.Register(context.Background(), RegisterAccountReq{
		UserID: "test",
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(openid)
}

func TestGameCenter_GameLink(t *testing.T) {
	g := NewGameCenter(GameCenterConfig{
		Username:   "90fae8ef04e68c6bc8dfdce7b1b1b3bc",
		Password:   "84363c525fc8f8b9703065772a918eb8",
		RequestURL: "https://center-openapi-gateway-dev.ogweb.org",
	})

	openid, err := g.GetGameLink(context.Background(), GetGameLinkReq{
		GameID:  128641,
		OpenID:  140552033009666,
		TaskID:  uuid.NewString(),
		Balance: "1000",
		Extend: Extend{
			DeviceOS: "ios",
			DeviceID: "123456",
			ClientIP: "127.0.0.1",
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(openid)
}

func TestGameCenter_GameList(t *testing.T) {
	g := NewGameCenter(GameCenterConfig{
		Username:   "90fae8ef04e68c6bc8dfdce7b1b1b3bc",
		Password:   "84363c525fc8f8b9703065772a918eb8",
		RequestURL: "https://center-openapi-gateway-dev.ogweb.org",
	})

	list, err := g.GetGameList(context.Background(), GetGameListReq{
		PlatformID: 1,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	utils.PrettyJSON(list)
}

func TestGameCenter_GetPlatformList(t *testing.T) {
	g := NewGameCenter(GameCenterConfig{
		Username:   "90fae8ef04e68c6bc8dfdce7b1b1b3bc",
		Password:   "84363c525fc8f8b9703065772a918eb8",
		RequestURL: "https://center-openapi-gateway-dev.ogweb.org",
	})
	list, err := g.GetPlatformList(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	utils.PrettyJSON(list)
}

func TestGameCenter_GetBetList(t *testing.T) {
	g := NewGameCenter(GameCenterConfig{
		Username:   "90fae8ef04e68c6bc8dfdce7b1b1b3bc",
		Password:   "84363c525fc8f8b9703065772a918eb8",
		RequestURL: "https://center-openapi-gateway-dev.ogweb.org",
	})
	list, err := g.GetBetList(context.Background(), GetBetListReq{
		StartTime: 1756801567,
		EndTime:   1756802107,
		OpenID:    140552033009666,
		Page:      1,
		PageSize:  10,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	utils.PrettyJSON(list)
}

func TestGameCenter_CheckStatus(t *testing.T) {
	g := NewGameCenter(GameCenterConfig{
		Username:   "90fae8ef04e68c6bc8dfdce7b1b1b3bc",
		Password:   "84363c525fc8f8b9703065772a918eb8",
		RequestURL: "https://center-openapi-gateway-dev.ogweb.org",
	})
	resp, err := g.CheckStatus(context.Background(), CheckTransferStatusReq{
		OpenID:  140552033009666,
		TaskIds: []string{"65d9edf2-bb97-48d5-b439-09ea271e8730"},
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	utils.PrettyJSON(resp)
}

func TestGameCenter_TransferOut(t *testing.T) {
	g := NewGameCenter(GameCenterConfig{
		Username:   "90fae8ef04e68c6bc8dfdce7b1b1b3bc",
		Password:   "84363c525fc8f8b9703065772a918eb8",
		RequestURL: "https://center-openapi-gateway-dev.ogweb.org",
	})
	resp, err := g.TransferOut(context.Background(), TransferOutReq{
		OpenID:     140552033009666,
		TaskID:     uuid.NewString(),
		PlatformID: 100028,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	utils.PrettyJSON(resp)
}
func TestGameCenter_GetBalance(t *testing.T) {
	g := NewGameCenter(GameCenterConfig{
		Username:   "90fae8ef04e68c6bc8dfdce7b1b1b3bc",
		Password:   "84363c525fc8f8b9703065772a918eb8",
		RequestURL: "https://center-openapi-gateway-dev.ogweb.org",
	})
	resp, err := g.GetBalance(context.Background(), GetBalanceReq{
		OpenID:      140552033009666,
		PlatformIds: []int64{100028},
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	utils.PrettyJSON(resp)
}
