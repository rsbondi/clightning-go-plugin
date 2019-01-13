package main

type RpcOptions struct {
	Crypto string `json:"crypto"`
	Fiat   string `json:"fiat"`
}

type RpcInfo struct {
	Result RpcInfoResult
}

type RpcInfoResult struct {
	Crypto string `json:"crypto"`
	Fiat   string `json:"fiat"`
}
