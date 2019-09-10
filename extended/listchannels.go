package main

import (
	"github.com/niftynei/glightning/jrpc2"
)

type ListChannelsExt struct {
	ShortChannelIds []string `json:"short_channel_ids"`
}

func (h *ListChannelsExt) New() interface{} {
	return &ListChannelsExt{}
}

func (h *ListChannelsExt) Name() string {
	return "listchannelsext"
}

func (h *ListChannelsExt) Call() (jrpc2.Result, error) {
	return "TBD", nil
}
