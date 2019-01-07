package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

var fundinfo *FundPrice
var rpcFile string

func jsonInit(msg json.RawMessage) interface{} {
	var params RpcInitParams
	if err := json.Unmarshal(msg, &params); err != nil {
	}

	rpcFile = fmt.Sprintf("%s/%s", params.Configuration.LightningDir, params.Configuration.RpcFile)

	fundinfo = NewFundPrice("https://apiv2.bitcoinaverage.com/indices/local/ticker/short",
		params.Options.Crypto,
		params.Options.Fiat)
	return "ok"
}

func jsonGetManifest(msg json.RawMessage) interface{} {
	response := RpcInit{
		Options: []RpcInitOptions{
			{
				Name:        "crypto",
				Type:        "string",
				Default:     "BTC",
				Description: "Ticker symbol for crypto currency.",
			},
			{
				Name:        "fiat",
				Type:        "string",
				Default:     "USD",
				Description: "Ticker symbol for fiat currency.",
			},
		},
		Rpcmethods: []RpcMethods{
			{
				Name:        "fundprice",
				Description: "Returns a summerized fund data with price",
			},
		},
	}
	return response
}

func jsonFundPrice(msg json.RawMessage) interface{} {
	c, err := net.Dial("unix", rpcFile)
	if err != nil {
	}

	c.Write([]byte(`{"jsonrpc":"2.0","id":98,"method":"listfunds","params":[]}`))
	buf := make([]byte, 1024)
	n, err := c.Read(buf[:])
	if err != nil {
	}

	var funds RpcFundsResult
	rpcfunds := RpcResult{
		Result: &funds,
	}
	if err := json.Unmarshal(buf[0:n], &rpcfunds); err != nil {
	}

	pricebuf := make([]byte, 1024)

	response, err := http.Get(fundinfo.ApiRequest()) // https://apiv2.bitcoinaverage.com/indices/local/ticker/short?crypto=BTC&fiat=USD
	n, err = response.Body.Read(pricebuf)

	var m map[string]ApiResult
	var totals Funds
	errp := json.Unmarshal(pricebuf[0:n], &m)
	if errp == nil {
		price := m[fundinfo.ResponseSymbol()].Bid

		var channelFunds int64 = 0
		for _, ch := range funds.Channels {
			channelFunds += ch.ChannelSat
		}
		var chainFunds int64 = 0
		for _, o := range funds.Outputs {
			chainFunds += o.Value
		}

		conversion := FundConvert{
			Fiat:    price,
			Divisor: float32(100000000),
		}

		totals = Funds{
			Chain:   *NewFund(chainFunds, conversion),
			Channel: *NewFund(channelFunds, conversion),
		}
	}

	return totals
}

func main() {
	commands := map[string]rpcfun{
		"init":        jsonInit,
		"fundprice":   jsonFundPrice,
		"getmanifest": jsonGetManifest,
	}

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	for {
		var msg json.RawMessage
		cmd := RpcCommand{
			Params: &msg,
		}
		err := json.NewDecoder(reader).Decode(&cmd)

		if err != nil {
		}
		method, ok := commands[cmd.Method]
		if ok {
			rpcResponse := RpcResult{
				Id:      cmd.Id,
				Jsonrpc: "2.0",
				Result:  method(msg),
			}

			json.NewEncoder(writer).Encode(rpcResponse)
			writer.Flush()
			writer.Reset(os.Stdout)
			reader.Reset(os.Stdin)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

/*
 {"jsonrpc":"2.0","id":99,"method":"init","params":{ "configuration": {"lightning-dir": "/home/richard/.lightning","rpc-file":"lightning-rpc"    }  }}}
 {"jsonrpc":"2.0","id":99,"method":"fundprice","params":[]}
*/
