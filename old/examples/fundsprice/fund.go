package main

import (
	"strings"
)

type FundPrice struct {
	Url    string `json:"url"`
	Crypto string `json:"crypto"`
	Fiat   string `json:"fiat"`
}

func NewFundPrice(url string, crypto string, fiat string) *FundPrice {
	var cryptoSymbol string
	if crypto != "" {
		cryptoSymbol = crypto
	} else {
		cryptoSymbol = "BTC"
	}

	var fiatSymbol string
	if fiat != "" {
		fiatSymbol = fiat
	} else {
		fiatSymbol = "USD"
	}

	fundinfo = &FundPrice{
		Url:    url,
		Crypto: cryptoSymbol,
		Fiat:   fiatSymbol,
	}

	return fundinfo
}

func (fp *FundPrice) ResponseSymbol() string {
	s := []string{strings.ToUpper(fp.Crypto), strings.ToUpper(fp.Fiat)}
	return strings.Join(s, "")
}

func (fp *FundPrice) ApiRequest() string {
	s := []string{fp.Url, "?crypto=", strings.ToUpper(fp.Crypto), "&fiat=", strings.ToUpper(fp.Fiat)}
	return strings.Join(s, "")
}

type Output struct {
	Txid    string `json:"txid"`
	Output  int    `json:"output"`
	Value   int64  `json:"value"`
	Address string `json:"address"`
	Status  string `json:"string"`
}

type Channel struct {
	PeerId          string `json:"peer_id"`
	ShortChannelId  string `json:"short_channel_id"`
	ChannelSat      int64  `json:"channel_sat"`
	ChannelTotalSat int64  `json:"channel_total_sat"`
	FundingTxid     string `json:"funding_txid"`
}

type RpcFundsResult struct {
	Outputs  []Output  `json:"outputs"`
	Channels []Channel `json:"channels"`
}

type ApiResult struct {
	Bid float32 `json:"bid"`
}

type Fund struct {
	Amount int64
	Value  float32
}

type Funds struct {
	Chain   Fund
	Channel Fund
}

type FundConvert struct {
	Fiat    float32
	Divisor float32
}

func NewFund(amt int64, conv FundConvert) *Fund {
	fund := &Fund{
		Amount: amt,
		Value:  float32(amt) / conv.Divisor * conv.Fiat,
	}

	return fund
}
