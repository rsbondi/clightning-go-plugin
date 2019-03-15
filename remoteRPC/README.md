### remoteRPC

This plugin allows you to access all RPC commands using HTTP instead of the default of unix socket.  Not encrypted use within local secure network only.

#### Usage

use defaults

`lightningd --plugin=path/to/plugin --remote-user=[username] --remote-password=[password]` 

will make RPC calls available on localhost:9222

specify alternate port

`lightningd --plugin=path/to/plugin --remote-port=1234 ...`

will listen on localhost:1234

#### Limitation

uses http passthrough to unix socket, use with `wait...` commands only work within the request timeout
