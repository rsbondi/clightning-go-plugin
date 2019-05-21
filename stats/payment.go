package main

import (
	"github.com/niftynei/glightning/glightning"
	"github.com/niftynei/glightning/jrpc2"
)

type Payment struct{}

func (f *Payment) New() interface{} {
	return &Payment{}
}

func (f *Payment) Name() string {
	return "paymentstats"
}

func (z *Payment) Call() (jrpc2.Result, error) {
	return paymentSummary()
}

type ListSendPaysRequest struct{}

func (r *ListSendPaysRequest) Name() string {
	return "listsendpays"
}

func ListSendPays() ([]glightning.PaymentFields, error) {
	var result struct {
		Payments []glightning.PaymentFields `json:"payments"`
	}
	req := &ListSendPaysRequest{}
	err := lightning.Request(req, &result)
	return result.Payments, err
}

func paymentSummary() ([]glightning.PaymentFields, error) {
	payments, err := ListSendPays()
	return payments, err
}
