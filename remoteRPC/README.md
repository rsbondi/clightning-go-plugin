### remoteRPC

This plugin allows you to access all RPC commands using HTTP instead of the default of unix socket.  Can use encrypted or unencrypted within local secure network.

[linux-x64 binary](https://moonbreeze.richardbondi.net/remote_plugin)

#### Usage

use defaults

`lightningd --plugin=path/to/plugin --remote-user=[username] --remote-password=[password]` 

will make RPC calls available on localhost:9222

specify alternate port

`lightningd --plugin=path/to/plugin --remote-port=1234 ...`

will listen on localhost:1234

`lightningd --plugin=path/to/plugin --remote-cert=PATH/TO/CERT --remote-key=PATH/TO/KEY ...`

will serve https

#### Limitation

uses http passthrough to unix socket, use with `wait...` commands only work within the request timeout
