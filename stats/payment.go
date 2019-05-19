package main

import (
	"fmt"

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

func paymentSummary() (interface{}, error) {
	payments, err := lightning.ListPayments("")
	if err != nil {
		return fmt.Sprintf("forward: %s\n", err.Error()), nil
	}

	if len(payments) == 0 {
		return "no forwarding information available", nil
	}

	return &Payment{}, nil
}
