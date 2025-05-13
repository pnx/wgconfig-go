# wgconfig-go

[![Test](https://github.com/pnx/wgconfig-go/actions/workflows/test.yml/badge.svg)](https://github.com/pnx/wgconfig-go/actions/workflows/test.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/pnx/wgconfig-go.svg)](https://pkg.go.dev/github.com/pnx/wgconfig-go)

Simple go module that can manipulate [Wireguard](https://www.wireguard.com) config files.

## Install

```sh
go get github.com/pnx/wgconfig-go@latest
```

## Example

```go
package main

import (
    "encoding/json"
    "log"
    "os"
    "strings"

    "github.com/pnx/wgconfig-go"
)

func main() {
    var cfg wgconfig.Config
    r := strings.NewReader(`
        [Interface]
        PrivateKey = XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX=
        Address = 10.10.10.4/32
        DNS = 8.8.4.4
        MTU = 1420

        [Peer]
        PublicKey = YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY=
        PresharedKey = ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ=
        AllowedIPs = 0.0.0.0/0
        Endpoint = example.com:51820
        `)

    err := cfg.Read(r)
    if err != nil {
        log.Fatal(err)
    }

    err = json.NewEncoder(os.Stdout).Encode(cfg)
    if err != nil {
        panic(err)
    }
}
```

Or reading from from a file directly:

```go
err := cfg.ReadFile("/path/to/file.ini")
```

## Documentation

Documentation can be found over at [pkg.go.dev](https://pkg.go.dev/github.com/pnx/wgconfig-go)

## Author

Henrik Hautakoski - [henrik.hautakoski@gmail.com](henrik.hautakoski@gmail.com)
