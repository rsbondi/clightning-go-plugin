Example clightning plugins in go

### hello
basic hello world plugin

### fundsprice

show a summary of total channel funds and chain funds, additionally calling bitcoinaverage api and displaying both totals in selected currency.  Use the `--crypto` flag if not BTC, and the `--fiat` flag for other than USD when starting the daemon

TODO: add arguments for units(mbtc, bits, btc, satoshis)