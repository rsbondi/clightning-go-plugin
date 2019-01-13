package main

import (
	"clplugin"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
)

var fundinfo *FundPrice
var plug *clplugin.Plugin

func jsonFundPrice(msg json.RawMessage) interface{} {
	c, err := net.Dial("unix", plug.RpcFile())
	if err != nil {
	}

	plug.Log("info", "Calling listfunds plugin method")
	c.Write([]byte(`{"jsonrpc":"2.0","id":98,"method":"listfunds","params":[]}`))
	buf := make([]byte, 1024)
	n, err := c.Read(buf[:])
	if err != nil {
		plug.Log("warn", fmt.Sprintf("Calling listfunds plugin method failed: %s", err.Error()))
	}

	var funds RpcFundsResult
	rpcfunds := clplugin.RpcResult{
		Result: &funds,
	}
	if err := json.Unmarshal(buf[0:n], &rpcfunds); err != nil {
		plug.Log("warn", fmt.Sprintf("Parsing listfunds response failed: %s", err.Error()))
	}

	pricebuf := make([]byte, 1024)

	plug.Log("info", fmt.Sprintf("Making request to bitcoinaverage.com for pricing info fiat=%s, crypto=%s", fundinfo.Fiat, fundinfo.Crypto))
	response, err := http.Get(fundinfo.ApiRequest()) // https://apiv2.bitcoinaverage.com/indices/local/ticker/short?crypto=BTC&fiat=USD
	n, err = response.Body.Read(pricebuf)

	var m map[string]ApiResult
	var totals Funds
	errp := json.Unmarshal(pricebuf[0:n], &m)
	var price float32
	if errp != nil {
		price = 0
		plug.Log("warn", fmt.Sprintf("Parsing pai response failed: %s", errp.Error()))
	} else {
		price = m[fundinfo.ResponseSymbol()].Bid
	}

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

	return totals
}

// Example plugin createion.
// 1) create a new instance with `clplugin.NewPlugin()`
// 2) create methods called by plugin
// 3) add methods with call to `AddMethod`, you can call multiple times
// 4) add any options used by plugin with calls to `AddOption`
// 5) optionally handle additional initialization with `AddInit`(only once)
//    this will help to get any command line arguments passed to the daemon.
//    In this case, the "fiat" and "crypto" values
// 6) Call `Run()`
func main() {

	plug = clplugin.NewPlugin()
	plug.AddMethod("fundprice", "show fund summary with price", jsonFundPrice)

	plug.AddOption("fiat", "USD", "Ticker symbol for fiat currency.")
	plug.AddOption("crypto", "BTC", "Ticker symbol for crypto currency.")

	plug.AddInit(func(msg json.RawMessage) {
		plug.Log("info", "Additional init options added for fundprice plugin")
		var options RpcOptions
		if err := json.Unmarshal(msg, &options); err != nil {
			// additional handling
		}

		fundinfo = NewFundPrice("https://apiv2.bitcoinaverage.com/indices/local/ticker/short",
			options.Crypto,
			options.Fiat)

	})

	plug.Run()

}

/*
 {"jsonrpc":"2.0","id":99,"method":"init","params":{ "configuration": {"lightning-dir": "/home/richard/.lightning","rpc-file":"lightning-rpc"    }  }}}
 {"jsonrpc":"2.0","id":99,"method":"fundprice","params":[]}
*/
