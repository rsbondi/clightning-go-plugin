This plugin takes no additional options so to use, start `lightningd` with `--plugin` option

`lightningd --plugin=path/to/plugin`

## Forwards

This will add `forwardstats` and `forwardview` cli/rpc commands that takes no arguments

`forwardview` returns a simple, inflexible html view with pie charts and fee info, mainly
geared toward [a mobile app](https://github.com/rsbondi/clightning-mobile) customization

`forwardstats` result

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

## Payments

This will also add `paymentstats` and `paymentview` cli/rpc commands that takes no arguments

`paymentstats` result

```json
{
    "complete": {
      "average": 5587785,
      "median": 5941000,
      "count": 14,
      "total": 78229000,
      "rate": 25,
      "min": 1624000,
      "max": 9845000
    },
    "failed": {
      "average": 2098273,
      "median": 572534,
      "count": 42,
      "total": 88127476,
      "rate": 75,
      "min": 48608,
      "max": 9831000
    }
  }
```

## Channel activity

`channel_activity` command will show how active a channel is, a channel may recieve, send or forward payments.  This provides the aggregate of all successful events. 

`channel_activity` results

```json
[
    {
      "short_channel_id": "560x1x0",
      "msatoshi": 3642026,
      "direction": "receive"
    },
    {
      "short_channel_id": "560x1x0",
      "msatoshi": 63174198,
      "direction": "send"
    },
    {
      "short_channel_id": "560x3x0",
      "msatoshi": 7332147,
      "direction": "receive"
    },
    {
      "short_channel_id": "560x3x0",
      "msatoshi": 52883940,
      "direction": "send"
    },
    {
      "short_channel_id": "560x4x0",
      "msatoshi": 102357847,
      "direction": "receive"
    },
    {
      "short_channel_id": "560x4x0",
      "msatoshi": 47073072,
      "direction": "send"
    }
  ]
```