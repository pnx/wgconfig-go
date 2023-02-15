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
