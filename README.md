Example clightning plugins in go

### hello
basic hello world plugin

### fundsprice

show a summary of total channel funds and chain funds, additionally calling bitcoinaverage api and displaying both totals in selected currency.  Use the `--crypto` flag if not BTC, and the `--fiat` flag for other than USD when starting the daemon.  This only works with short list and will be updated when plugin package(below) is complete

### clplugin

this directory is an attempt to create a package for creating plugins, similar to they python version provided in the clightning repo. I have migrated the fundsprice to use it.  The `hello` plugin will not use this as it better illustrates the commands and responses

Usage

```golang
    p = clplugin.NewPlugin()
    p.AddMethod("methodname", "description of method", methodImplementation)
    p.AddOption("option1", "default value", "option description")
    // ... add as many options as you like

    p.AddInit(func(msg json.RawMessage) {
           // optional if you need to do any additional actions on call to `init`
           // for example reading option values passed on command line to daemon
    })
    p.Run()
```

[Here](https://github.com/niftynei/glightning) is a much more complete effort of what I intended to achieve, I suggest using it and see no need for me to continue with a duplicate effort