package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/niftynei/glightning/glightning"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

const VERSION = "0.0.2"

var plugin *glightning.Plugin
var remote *RemoteRPC
var local net.Conn

type RemoteRPC struct {
	Username string
	Password string
	Port     string
	RPCFile  string
}

type RpcResult struct {
	Jsonrpc string      `json:"jsonrpc"`
	Id      int         `json:"id"`
	Result  interface{} `json:"result"`
}

func NewRemoteRPC(options map[string]glightning.Option, rpcfile string) *RemoteRPC {
	return &RemoteRPC{
		Username: options["remote-user"].GetValue().(string),
		Password: options["remote-password"].GetValue().(string),
		Port:     options["remote-port"].GetValue().(string),
		RPCFile:  rpcfile,
	}
}

func main() {
	plugin = glightning.NewPlugin(onInit)

	registerOptions(plugin)
	err := plugin.Start(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func auth(req *http.Request) bool {
	authHeader := strings.SplitN(req.Header.Get("Authorization"), " ", 2)

	if len(authHeader) != 2 || authHeader[0] != "Basic" {
		return false
	}

	basic, _ := base64.StdEncoding.DecodeString(authHeader[1])
	userpass := strings.SplitN(string(basic), ":", 2)

	if userpass[0] == remote.Username && userpass[1] == remote.Password {
		return true
	}
	return false
}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	if !auth(req) {
		// TODO: id
		rpcerr := &RpcResult{
			Jsonrpc: "2.0",
			Result:  "Not Authorized",
		}
		rpcResponse, _ := json.Marshal(rpcerr)
		w.Write(rpcResponse)
		return
	}
	local, err := net.Dial("unix", remote.RPCFile)
	if err != nil {
		log.Fatal("unable to connect to clightning")
	}
	defer local.Close()

	var unix2http = make([]byte, 1024)
	var responseBuf = make([]byte, 0)

	_, errc := io.Copy(local, req.Body)
	if errc != nil && errc != io.EOF {
		log.Printf("Copy error: %s", errc)
	}

	for {
		r, err := local.Read(unix2http)
		if err != nil {
			if err != io.EOF {
				log.Printf("RPC error to clightning: %s", err.Error())
			}
			break
		}
		responseBuf = append(responseBuf, unix2http[:r]...)
		if unix2http[r-2] == '\n' && unix2http[r-1] == '\n' {
			break
		}
	}
	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.Write(responseBuf)
}

func registerOptions(p *glightning.Plugin) {
	p.RegisterOption(glightning.NewOption("remote-user", "User used for authentication", "remoteuser"))
	p.RegisterOption(glightning.NewOption("remote-password", "Password used for authentication", "remotepass"))
	p.RegisterOption(glightning.NewOption("remote-port", "Port to listen for remote requests", "9222"))
	p.RegisterOption(glightning.NewOption("remote-cert", "Server certificate", " ")) // crashes on empty
	p.RegisterOption(glightning.NewOption("remote-key", "Server key", " "))
}

func onInit(plugin *glightning.Plugin, options map[string]glightning.Option, config *glightning.Config) {
	log.Printf("versiion: %s initialized for port %s\n", VERSION, options["remote-port"])
	rpcfile := fmt.Sprintf("%s/%s", config.LightningDir, config.RpcFile)
	remote = NewRemoteRPC(options, rpcfile)
	http.HandleFunc("/", handleRequest)

	go (func() {
		if options["remote-cert"].GetValue().(string) != " " {
			log.Fatal(http.ListenAndServeTLS(":"+remote.Port, options["remote-cert"].GetValue().(string), options["remote-key"].GetValue().(string), nil))
		} else {
			log.Fatal(http.ListenAndServe(":"+remote.Port, nil))
		}
	})()
}
