package main

import (
	"bytes"
	"fmt"
	"github.com/niftynei/glightning/jrpc2"
	"html/template"
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

type PaymentView struct{}

func (f *PaymentView) New() interface{} {
	return &PaymentView{}
}

func (f *PaymentView) Name() string {
	return "paymentview"
}

func (z *PaymentView) Call() (jrpc2.Result, error) {
	sum, err := paymentSummary()
	if err != nil {
		return fmt.Sprintf("payment: %s\n", err.Error()), nil
	}

	html := `<body>
	<h2>Payment Summary</h2>
	<table>
	  <thead>
		<tr>
		  <th></th>
		  <th>Completed</th>
		  <th>Failed</th>
		</tr>
	  </thead>
	  <tbody>
		<tr>
		  <th>%</th>
		  <td>{{with .Complete}} {{printf "%2.1f" .Rate}} {{end}}</td>
		  <td>{{with .Failed}} {{printf "%2.1f" .Rate}} {{end}}</td>
		</tr>
		<tr>
		  <th>Count</th>
		  <td>{{.Complete.Count}}</td>
		  <td>{{.Failed.Count}}</td>
		</tr>
		<tr>
		  <th>Min</th>
		  <td>{{.Complete.Min}}</td>
		  <td>{{.Failed.Min}}</td>
		</tr>
		<tr>
		  <th>Max</th>
		  <td>{{.Complete.Max}}</td>
		  <td>{{.Failed.Max}}</td>
		</tr>
		<tr>
		  <th>Average</th>
		  <td>{{.Complete.Average}}</td>
		  <td>{{.Failed.Average}}</td>
		</tr>
		<tr>
		  <th>Median</th>
		  <td>{{.Complete.Median}}</td>
		  <td>{{.Failed.Median}}</td>
		</tr>
	  </tbody>
	</table>
	<style>
	  body {
		font-family: Arial, Helvetica, sans-serif;
	  }
	  table {
	  }
	  td, th {
		text-align: right;
	  }
	</style>
  </body>`

	tmpl, err := template.New("view").Parse(html)
	var b bytes.Buffer
	err = tmpl.Execute(&b, sum)

	return b.String(), nil

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
	payments, err := lightning.ListSendPaysAll() //ListSendPays()
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
		amounts[p.Status] = append(amounts[p.Status], p.AmountMilliSatoshiRaw)
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
			Rate:    100 * float32(len(amounts["complete"])) / float32(len(payments)),
			Min:     minc,
			Max:     maxc,
		},
		Failed: PaymentDetail{
			Average: avg(amounts["failed"]),
			Median:  med(amounts["failed"]),
			Count:   len(amounts["failed"]),
			Total:   sumf,
			Rate:    100 * float32(len(amounts["failed"])) / float32(len(payments)),
			Min:     minf,
			Max:     maxf,
		},
	}
	return result, err
}
