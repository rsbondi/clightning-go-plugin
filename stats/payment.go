package main

import (
	"github.com/niftynei/glightning/glightning"
	"github.com/niftynei/glightning/jrpc2"
	"sort"
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

type PaymentDetail struct {
	Average uint64  `json:"average"`
	Median  uint64  `json:"median"`
	Count   int     `json:"count"`
	Total   uint64  `json:"total"`
	Rate    float32 `json:"rate"`
	Min     uint64  `json:"min"`
	Max     uint64  `json:"max"`
}

type PaymentResult struct {
	Complete PaymentDetail `json:"complete"`
	Failed   PaymentDetail `json:"failed"`
}

type ByMsat []uint64

func (a ByMsat) Len() int           { return len(a) }
func (a ByMsat) Less(i, j int) bool { return a[i] < a[j] }
func (a ByMsat) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func avg(msats []uint64) uint64 {
	if len(msats) == 0 {
		return 0
	}
	var total uint64 = 0
	for _, msat := range msats {
		total += msat
	}
	return total / uint64(len(msats))
}

func statloop(msats []uint64) (uint64, uint64, uint64) {
	var total uint64
	var min uint64
	var max uint64
	for i, msat := range msats {
		total += msat
		if i > 0 {
			if msat < min {
				min = msat
			}
			if msat > max {
				max = msat
			}
		} else {
			min = msat
			max = msat
		}
	}
	return total, min, max
}

func med(msats []uint64) uint64 {
	l := len(msats)
	if l == 0 {
		return 0
	}
	if l%2 == 1 {
		return msats[(len(msats)-1)/2]
	}
	return avg(msats[l/2-1 : l/2+1])
}

func paymentSummary() (PaymentResult, error) {
	payments, err := ListSendPays()
	if err != nil {
		return PaymentResult{}, err
	}

	if len(payments) == 0 {
		return PaymentResult{}, nil
	}

	amounts := make(map[string][]uint64)
	amounts["complete"] = []uint64{}
	amounts["failed"] = []uint64{}

	for _, p := range payments {
		amounts[p.Status] = append(amounts[p.Status], p.MilliSatoshi)
	}

	sort.Sort(ByMsat(amounts["complete"]))
	sort.Sort(ByMsat(amounts["failed"]))

	sumc, minc, maxc := statloop(amounts["complete"])
	sumf, minf, maxf := statloop(amounts["failed"])
	result := PaymentResult{
		Complete: PaymentDetail{
			Average: avg(amounts["complete"]),
			Median:  med(amounts["complete"]),
			Count:   len(amounts["complete"]),
			Total:   sumc,
			Rate:    float32(len(amounts["complete"])) / float32(len(payments)),
			Min:     minc,
			Max:     maxc,
		},
		Failed: PaymentDetail{
			Average: avg(amounts["failed"]),
			Median:  med(amounts["failed"]),
			Count:   len(amounts["failed"]),
			Total:   sumf,
			Rate:    float32(len(amounts["failed"])) / float32(len(payments)),
			Min:     minf,
			Max:     maxf,
		},
	}
	return result, err
}
