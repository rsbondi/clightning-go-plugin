Example clightning plugins in go

### hello
basic hello world plugin

### fundsprice

show a summary of total channel funds and chain funds, additionally calling bitcoinaverage api and displaying both totals in selected currency.  Use the `--crypto` flag if not BTC, and the `--fiat` flag for other than USD when starting the daemon.  This only works with short list and will be updated when plugin class(below) is complete

### plugin

this directory is an attempt to create a package for creating plugins, similar to they python version provided in the clightning repo. I have migrated the fundsprice to use it.  The `hello` plugin will not use this as it better illustrates the commands and responses