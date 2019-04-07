This plugin takes no additional options so to use, start `lightningd` with `--plugin` option

`lightningd --plugin=path/to/plugin`

This will add `forwardstats` cli/rpc command that takes no arguments

It will return something like the following

```json
{
    "chins": {
      "3139x13x0": {
        "msat": 7,
        "funding": 0,
        "percent_gain": 0,
        "percent_pie": 0.3333333333333333
      },
      "3139x15x0": {
        "msat": 6,
        "funding": 0,
        "percent_gain": 0,
        "percent_pie": 0.2857142857142857
      },
      "3139x17x0": {
        "msat": 8,
        "funding": 0,
        "percent_gain": 0,
        "percent_pie": 0.38095238095238093
      }
    },
    "chouts": {
      "3139x10x0": {
        "msat": 13,
        "funding": 1855572000,
        "percent_gain": 7.005925935506679e-9,
        "percent_pie": 0.6190476190476191
      },
      "3139x3x0": {
        "msat": 8,
        "funding": 7908273000,
        "percent_gain": 1.0115988661494108e-9,
        "percent_pie": 0.38095238095238093
      }
    },
    "totalfunding": 9763845000,
    "totalfees": 21,
    "total_percent_gain": 2.15079202916474e-9
  }
```

with channels in(`chins`) and channels out(`chouts`) each indicated by the short channel id.
For each channel, you see total `msat` forwarded, `funding` or how much you put into the channel,
`percent_gain` or how much you earned based on what you put in and `percent_pie` or what percent
of forwarding activity has gone through the channel.

Additionally, totals are displayed at the end.

**Suggestions welcome for additional stats that may be useful to implement.**