package main

import (
	"github.com/niftynei/glightning/jrpc2"
)

type ListConfigsExt struct {
	Configs []string `json:"configs"`
}

func (h *ListConfigsExt) New() interface{} {
	return &ListConfigsExt{}
}

func (h *ListConfigsExt) Name() string {
	return "listconfigsext"
}

func (h *ListConfigsExt) Call() (jrpc2.Result, error) {
	return "TBD", nil
}
