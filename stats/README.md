This plugin takes no additional options so to use, start `lightningd` with `--plugin` option

`lightningd --plugin=path/to/plugin`

This will add `forwardstats` and `forwardview` cli/rpc commands that takes no arguments

`forwardview` returns a simple, inflexible html view with pie charts and fee info, mainly
geared toward [a mobile app](https://github.com/rsbondi/clightning-mobile) customization

`forwardstats` will return something like the following

```json
{
  "channels_in": {
    "3265x3x0": {
      "fee_msat": 30,
      "forward_msat": 1546027,
      "funding": 0,
      "percent_gain": 0,
      "percent_pie": 1
    }
  },
  "channels_out": {
    "3265x10x0": {
      "fee_msat": 10,
      "forward_msat": 501000,
      "funding": 3971688000,
      "percent_gain": 2.517821138014869e-9,
      "percent_pie": 0.3333333333333333
    },
    "3265x13x0": {
      "fee_msat": 20,
      "forward_msat": 1045027,
      "funding": 7909922000,
      "percent_gain": 2.528469939400161e-9,
      "percent_pie": 0.6666666666666666
    }
  },
  "totalfunding": 15891391000,
  "totalfees": 30,
  "totalforward": 1546027,
  "total_percent_gain": 1.887814603517087e-9
}
```

with `channels_in` and `channels_out` each indicated by the short channel id.
For each channel, you see total fees collected, `fee_msat` total forwarded, `forward_msat`, 
`funding` or how much you put into the channel,
`percent_gain` or how much you earned based on what you put in and `percent_pie` or what percent
of forwarding activity has gone through the channel.

Additionally, totals are displayed at the end.

**Suggestions welcome for additional stats that may be useful to implement.**