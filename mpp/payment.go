package main

import (
	"fmt"

	"github.com/niftynei/glightning/glightning"
	"github.com/niftynei/glightning/jrpc2"
)

type PaymentMpp struct{}

func (f *PaymentMpp) New() interface{} {
	return &PaymentMpp{}
}

func (f *PaymentMpp) Name() string {
	return "mpp_payments"
}

func (z *PaymentMpp) Call() (jrpc2.Result, error) {
	return paymentSummary()
}

type SendPayFieldsMpp struct {
	Id                 uint64 `json:"id,omitempty"`
	AmountMilliSatoshi string `json:"amount_msat,omitempty"`
	glightning.SendPayFields
	Parts int `json:"parts,omitempty"`
}

func paymentSummary() ([]*SendPayFieldsMpp, error) {
	payments, err := lightning.ListSendPaysAll()
	if err != nil {
		r := []*SendPayFieldsMpp{}
		return r, nil
	}

	if len(payments) == 0 {
		r := []*SendPayFieldsMpp{}
		return r, nil
	}

	var result []*SendPayFieldsMpp
	resultMap := make(map[string]*SendPayFieldsMpp)

	for _, p := range payments {
		identifier := fmt.Sprintf("%s-%s", p.PaymentHash, p.Status)
		pay, ok := resultMap[identifier]
		if !ok {
			resultMap[identifier] = &SendPayFieldsMpp{}
			pay = resultMap[identifier]
			pay.PaymentHash = p.PaymentHash
		}

		var multi bool
		if p.PartId == 0 {
			multi = false
		} else {
			multi = true
		}

		if p.PaymentPreimage != "" {
			pay.PaymentPreimage = p.PaymentPreimage
		}

		pay.MilliSatoshiSentRaw += p.MilliSatoshiSentRaw
		if multi {
			pay.Parts++
			pay.MilliSatoshiSent = fmt.Sprintf("%dmsat", pay.MilliSatoshiSentRaw)
		} else {
			pay.Id = p.Id
			pay.MilliSatoshiSent = p.MilliSatoshiSent
		}
		pay.Status = p.Status
		pay.CreatedAt = p.CreatedAt

	}

	for _, value := range resultMap {
		result = append(result, value)
	}

	return result, err
}
