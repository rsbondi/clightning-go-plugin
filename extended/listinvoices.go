package main

import (
	"github.com/niftynei/glightning/jrpc2"
)

type ListInvoicesExt struct {
	Labels []string `json:"labels"`
}

func (h *ListInvoicesExt) New() interface{} {
	return &ListInvoicesExt{}
}

func (h *ListInvoicesExt) Name() string {
	return "listinvoicesext"
}

func (h *ListInvoicesExt) Call() (jrpc2.Result, error) {
	return "TBD", nil
}
