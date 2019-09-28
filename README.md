A collection of plugins for use with clightning written in go using [glightning](https://github.com/niftynei/glightning)

## Plugins

### remoteRPC

This plugin allows you to access all RPC commands using HTTP instead of the default of unix socket

[more](remoteRPC/README.md)

### stats

Shows some additional stats, currently shows forwarding info by channel, amount earned, percent gain etc..

[more](stats/README.md)

### dumpkeys

Export xpriv/xpub

[more](dump_keys/README.md)

### setban

ban peers for specified time

[more](setban/README.md)

### multifund

fund multiple channels with single transaction

[more](multifund/README.md)

## General Plugin Installation

[download and install go](https://golang.org/dl/)

```
git clone https://github.com/rsbondi/clightning-go-plugin.git

cd clightning-go-plugin

# the DIRECTORY below is for the plugin you want to build
# example if you want to build remoteRPC plugin, use that directory
# this will build for your current architecture, for others see go documentation
go build -o path/to/where/you/want DIRECTORY
```

Once the plugin is built, there are several ways to use it.

Add plugin info to the [config](https://github.com/ElementsProject/lightning/blob/fd63b8bf53b9a14f29701d1a8cc57b6bee6b1096/doc/lightningd-config.5.txt#L325) file `plugin=path/to/plugin`

Also in the above link, you can specify `plugin-dir`option and put the plugin there, or note from the ling you can also use the defualt plugin directory.

In addition to the config file,you can also add these options when launching the daemon

`lightningd --plugin=path/to/plugin [other options]` 

or `--plugin-dir` as command line option

Additionally with version 0.7.2 forward, you can use the [cli/rpc](https://github.com/ElementsProject/lightning/blob/fd63b8bf53/doc/lightning-plugin.7.txt) to launch


##### _note_

_The `old` directory was my initial attempt at creating a plugin module for clightning but glightning is a much more complete effort of what I intended to achieve, I suggest using it and see no need for me to continue with a duplicate effort, but keeping for reference examples_