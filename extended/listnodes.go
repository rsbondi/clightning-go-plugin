package main

import (
	"github.com/niftynei/glightning/jrpc2"
)

type ListNodesExt struct {
	Ids []string `json:"ids"`
}

func (h *ListNodesExt) New() interface{} {
	return &ListNodesExt{}
}

func (h *ListNodesExt) Name() string {
	return "listnodesext"
}

func (h *ListNodesExt) Call() (jrpc2.Result, error) {
	return "TBD", nil
}
