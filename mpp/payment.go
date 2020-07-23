package main

import (
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
		pay, ok := resultMap[p.PaymentHash]
		if !ok {
			resultMap[p.PaymentHash] = &SendPayFieldsMpp{}
			pay = resultMap[p.PaymentHash]
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

		if multi {
			pay.Parts++
		} else {
			pay.Id = p.Id
		}
		pay.MilliSatoshiSentRaw += p.MilliSatoshiSentRaw
		pay.MilliSatoshiSent = p.MilliSatoshiSent
		pay.Status = p.Status
		pay.CreatedAt = p.CreatedAt

	}

	for _, value := range resultMap {
		result = append(result, value)
	}

	return result, err
}
