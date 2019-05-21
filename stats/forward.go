package main

import (
	"bytes"
	"fmt"
	"html/template"
	"math"

	"github.com/niftynei/glightning/glightning"
	"github.com/niftynei/glightning/jrpc2"
)

type Forwards struct{}

type ForwardSplit struct {
	Chins        interface{} `json:"channels_in"`
	Chouts       interface{} `json:"channels_out"`
	TotalFunding uint64      `json:"totalfunding"`
	TotalFees    uint64      `json:"totalfees"`
	TotalForward uint64      `json:"totalforward"`
	PercentGain  float64     `json:"total_percent_gain"`
}

type ForwardChan struct {
	MsatFees      uint64  `json:"fee_msat"`
	MsatForward   uint64  `json:"forward_msat"`
	FailedFees    uint64  `json:"fee_fail"`
	FailedForward uint64  `json:"forward_fail"`
	Funding       uint64  `json:"funding"`
	PercentGain   float64 `json:"percent_gain"`
	PercentPie    float64 `json:"percent_pie"`
}

type ChanSegment struct {
	ForwardChan
	Color string
	Px    float64
	Py    float64
	X     float64
	Y     float64
	Flag  uint
}

type ChannelViews struct {
	Ins  map[string]ChanSegment
	Outs map[string]ChanSegment
}

func (f *Forwards) New() interface{} {
	return &Forwards{}
}

func (f *Forwards) Name() string {
	return "forwardstats"
}

func (z *Forwards) Call() (jrpc2.Result, error) {
	return forwardSummary()
}

type ForwardView struct{}

func (f *ForwardView) New() interface{} {
	return &ForwardView{}
}

func (f *ForwardView) Name() string {
	return "forwardview"
}

var palette = []string{"#f0f0f0", "#c5d5c5", "#9fa9a3", "#e3e0cc", "#eaece5", "#b2c2bf",
	"#c0ded9", "#3b3a30", "#e4d1d1", "#b9b0b0", "#d9ecd0", "#77a8a8", "#f0efef",
	"#ddeedd", "#c2d4dd", "#b0aac0"}

func processView(c map[string]*ForwardChan, index int) map[string]ChanSegment {
	view := make(map[string]ChanSegment)
	px := math.Cos(0)
	py := math.Sin(0)
	var rotation float64 = 0

	for ch, fwd := range c {
		color := palette[index%len(palette)]
		rotation += fwd.PercentPie
		x := math.Cos(2 * math.Pi * rotation)
		y := math.Sin(2 * math.Pi * rotation)
		var f uint = 0
		if fwd.PercentPie > .5 {
			f = 1
		}
		chv := &ChanSegment{
			ForwardChan: *fwd,
			Color:       color,
			Px:          px,
			Py:          py,
			X:           x,
			Y:           y,
			Flag:        f,
		}
		view[ch] = *chv
		px = x
		py = y
		index++
	}

	return view
}

func (z *ForwardView) Call() (jrpc2.Result, error) {
	html := `<body>
    <div>
        <h3>Incoming Channels Fees</h3>
        {{range $ch, $data := .Ins}}
        <div style="clear: both; width: 100%; text-align: left;">
            <div style="width: 20px; height: 20px; float: left; background-color: {{$data.Color}};"></div>
            <div style="float: left; margin-left: 10px; width: 150px;">{{$ch}}</div>
            <div style="float: left;">{{$data.MsatFees}}</div>
        </div>
        {{end}}
        <svg viewBox="-1 -1 2 2" style="transform: rotate(-0.25turn)">
            {{range .Ins}}
            <path d="M {{.Px}} {{.Py}} A 1 1 0 {{.Flag}} 1 {{.X}} {{.Y}} L 0 0" fill="{{.Color}}"></path>
            {{end}}
        </svg>
    </div>
    <div>
        <h3>Outgoing Channels Fees</h3>
       {{range $ch, $data := .Outs}}        
        <div style="clear: both; width: 100%; text-align: left;">
            <div style="width: 20px; height: 20px; float: left; background-color: {{$data.Color}};"></div>
            <div style="float: left; margin-left: 10px; width: 150px;">{{$ch}}</div>
            <div style="float: left;">{{$data.MsatFees}}</div>
        </div>
        {{end}}
        <svg viewBox="-1 -1 2 2" style="transform: rotate(-0.25turn)">
            {{range .Outs}}
            <path d="M {{.Px}} {{.Py}} A 1 1 0 {{.Flag}} 1 {{.X}} {{.Y}} L 0 0" fill="{{.Color}}"></path>
            {{end}}
        </svg>
    </div>
</body>`
	sum, err := forwardSummary()
	if err != nil {
		return fmt.Sprintf("forward: %s\n", err.Error()), nil
	}

	view := &ChannelViews{}
	s := sum.(ForwardSplit)

	c := s.Chins.(map[string]*ForwardChan)
	index := 0
	view.Ins = processView(c, index)

	c = s.Chouts.(map[string]*ForwardChan)
	index = len(c) % len(palette)
	view.Outs = processView(c, index)

	tmpl, err := template.New("view").Parse(html)
	var b bytes.Buffer
	err = tmpl.Execute(&b, view)

	return b.String(), nil
}

func forwardSummary() (interface{}, error) {
	forwards, err := lightning.ListForwards()

	if err != nil {
		return fmt.Sprintf("forward: %s\n", err.Error()), nil
	}

	if len(forwards) == 0 {
		return "no forwarding information available", nil
	}

	peers, err := lightning.ListPeers()
	if err != nil {
		return fmt.Sprintf("forward: %s\n", err.Error()), nil
	}

	funds := make(map[string]uint64, 0)
	var totalfunding uint64
	for _, p := range peers {
		for _, c := range p.Channels {
			funds[c.ShortChannelId] = c.FundingAllocations[myid]
			totalfunding += c.FundingAllocations[myid]
		}
	}

	chins := make(map[string][]glightning.Forwarding, 0)
	chouts := make(map[string][]glightning.Forwarding, 0)
	var totalfees uint64
	var totalforwards uint64
	for _, f := range forwards {
		chins[f.InChannel] = append(chins[f.InChannel], f)
		chouts[f.OutChannel] = append(chouts[f.OutChannel], f)
		if f.Status == "settled" {
			totalfees += f.Fee
			totalforwards += f.MilliSatoshiOut
		}
	}

	chinsfinal := make(map[string]*ForwardChan, 0)
	choutsfinal := make(map[string]*ForwardChan, 0)

	processChannels(chins, chinsfinal, funds, totalfees)
	processChannels(chouts, choutsfinal, funds, totalfees)

	c := ForwardSplit{
		Chins:        chinsfinal,
		Chouts:       choutsfinal,
		TotalFunding: totalfunding,
		TotalFees:    totalfees,
		TotalForward: totalforwards,
		PercentGain:  float64(totalfees) / float64(totalfunding),
	}

	return c, nil
}

func processChannels(src map[string][]glightning.Forwarding,
	dest map[string]*ForwardChan,
	funds map[string]uint64,
	totalfees uint64) {
	for k, _ := range src {
		fees := uint64(0)
		forwarded := uint64(0)
		feesfail := uint64(0)
		forwardedfail := uint64(0)
		for _, f := range src[k] {
			if f.Status == "settled" {
				fees += f.Fee
				forwarded += f.MilliSatoshiOut
			} else {
				feesfail += f.Fee
				forwardedfail += f.MilliSatoshiOut
			}
		}
		var gain float64
		if funds[k] > 0 {
			gain = float64(fees) / float64(funds[k])

		}

		dest[k] = &ForwardChan{
			MsatFees:      fees,
			MsatForward:   forwarded,
			FailedFees:    feesfail,
			FailedForward: forwardedfail,
			Funding:       funds[k],
			PercentGain:   gain,
			PercentPie:    float64(fees) / float64(totalfees),
		}
	}
}
