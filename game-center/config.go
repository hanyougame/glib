package game_center

import "strings"

type GameCenterConfig struct {
	Username     string         `json:"username"`
	Password     string         `json:"password"`
	RequestURL   string         `json:"request_url"`
	CurrencyConf []CurrencyItem `json:"currency_conf"`
}

type CurrencyItem struct {
	Currency string `json:"currency"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c GameCenterConfig) GetCurrencyConf(currency string) CurrencyItem {
	for _, item := range c.CurrencyConf {
		if strings.ToLower(item.Currency) == strings.ToLower(currency) {
			return item
		}
	}
	return CurrencyItem{
		Currency: currency,
		Username: c.Username,
		Password: c.Password,
	}
}
