package main

import (
	"github.com/niftynei/glightning/jrpc2"
)

type ListPeersExt struct {
	Ids []string `json:"ids"`
}

func (h *ListPeersExt) New() interface{} {
	return &ListPeersExt{}
}

func (h *ListPeersExt) Name() string {
	return "listpeersext"
}

func (h *ListPeersExt) Call() (jrpc2.Result, error) {
	return "TBD", nil
}
