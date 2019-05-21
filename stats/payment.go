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

type AvgMedian struct {
	Average uint64 `json:"average"`
	Median  uint64 `json:"median"`
}

type PaymentResult struct {
	Complete AvgMedian `json:"complete"`
	Failed   AvgMedian `json:"failed"`
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

	amounts := make(map[string][]uint64)
	amounts["complete"] = []uint64{}
	amounts["failed"] = []uint64{}

	for _, p := range payments {
		amounts[p.Status] = append(amounts[p.Status], p.MilliSatoshi)
	}

	sort.Sort(ByMsat(amounts["complete"]))
	sort.Sort(ByMsat(amounts["failed"]))
	result := PaymentResult{
		Complete: AvgMedian{
			Average: avg(amounts["complete"]),
			Median:  med(amounts["complete"]),
		},
		Failed: AvgMedian{
			Average: avg(amounts["failed"]),
			Median:  med(amounts["failed"]),
		},
	}
	return result, err
}
