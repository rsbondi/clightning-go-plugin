package main

import (
	"github.com/niftynei/glightning/jrpc2"
)

type GetLogExt struct {
	Levels []string `json:"levels"`
}

func (h *GetLogExt) New() interface{} {
	return &GetLogExt{}
}

func (h *GetLogExt) Name() string {
	return "getlogext"
}

func (h *GetLogExt) Call() (jrpc2.Result, error) {
	return "TBD", nil
}
