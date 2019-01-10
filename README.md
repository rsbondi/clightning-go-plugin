Example clightning plugins in go

### hello
basic hello world plugin

### fundsprice

show a summary of total channel funds and chain funds, additionally calling bitcoinaverage api and displaying both totals in selected currency.  Use the `--crypto` flag if not BTC, and the `--fiat` flag for other than USD when starting the daemon.  This only works with short list and will be updated when plugin class(below) is complete

### plugin

this directory(WIP) is an attempt to create a "class" for creating plugins, it will migrate the others to use it, similar to they python version.