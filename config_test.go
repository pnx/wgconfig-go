package wgconfig

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Write(t *testing.T) {
	cfg := Config{
		Interface: Interface{
			ListenPort: 5559,
			Address:    "10.77.24.1/24",
			DNS:        []string{"8.8.8.8", "8.8.4.4"},
			MTU:        800,
			PrivateKey: "4CwbPHW85Y/xdgB/zD/P0bZdM3XVNpi85H45FMscB1A=",
			PostUp:     "iptables -A FORWARD -i %i -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE",
			PostDown:   "iptables -D FORWARD -i %i -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE",
		},
		Peers: Peers{
			{
				Comment:    "Comment1",
				PublicKey:  "gIIbPCSRw7qQnW/aS3g1PjZTEXnTBqSjo8sS9MADows=",
				AllowedIPs: "10.77.24.22/32",
				Endpoint:   "endpoint",
			},
			{
				Comment:    "Comment",
				PublicKey:  "KKJfVUC8awDEa4H7Pa5lRCvnq3cdrLMHpZVNF7YkgVA=",
				AllowedIPs: "10.77.24.24/32",
				Endpoint:   "endpoint2",
			},
			{
				PublicKey:           "IJgEGy5QPRbwuf7yY1+bbirFeHoNwdYzIfrWMNFEG30=",
				AllowedIPs:          "10.77.24.26/32",
				PersistentKeepalive: 30,
			},
			{
				PublicKey:  "NafllWlCPqa4Jhv10Rjbk38pxyWiWcpkwRYwcd47qic=",
				AllowedIPs: "10.77.24.28/32",
			},
		},
	}

	expected := `[Interface]
PrivateKey = 4CwbPHW85Y/xdgB/zD/P0bZdM3XVNpi85H45FMscB1A=
Address    = 10.77.24.1/24
ListenPort = 5559
DNS        = 8.8.8.8,8.8.4.4
MTU        = 800
PostUp     = iptables -A FORWARD -i %i -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
PostDown   = iptables -D FORWARD -i %i -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE

; Comment1
[Peer]
PublicKey  = gIIbPCSRw7qQnW/aS3g1PjZTEXnTBqSjo8sS9MADows=
AllowedIPs = 10.77.24.22/32
Endpoint   = endpoint

; Comment
[Peer]
PublicKey  = KKJfVUC8awDEa4H7Pa5lRCvnq3cdrLMHpZVNF7YkgVA=
AllowedIPs = 10.77.24.24/32
Endpoint   = endpoint2

[Peer]
PublicKey           = IJgEGy5QPRbwuf7yY1+bbirFeHoNwdYzIfrWMNFEG30=
AllowedIPs          = 10.77.24.26/32
PersistentKeepalive = 30

[Peer]
PublicKey  = NafllWlCPqa4Jhv10Rjbk38pxyWiWcpkwRYwcd47qic=
AllowedIPs = 10.77.24.28/32
`

	data := new(bytes.Buffer)
	_, err := cfg.Write(data)

	assert.NoError(t, err)
	assert.Equal(t, expected, data.String())
}

func TestConfig_Read(t *testing.T) {
	input := []byte(`[Interface]
Address = 10.8.16.1/24
ListenPort = 51820
PrivateKey = 4CwbPHW85Y/xdgB/zD/P0bZdM3XVNpi85H45FMscB1A=
DNS = 1.1.1.1,1.1.0.0
MTU = 1500
PreUp = echo "UP"
PostUp = iptables -A FORWARD -i %i -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
PreDown = echo "DOWN"
PostDown = iptables -D FORWARD -i %i -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE

# Peer 1
[Peer]
PublicKey = HHQSHN5TG6d0f3Wo0zeJM74v6073rQhc1+Yc8cwQ32Q=
AllowedIPs = 10.8.16.2/32
Endpoint = https://example.com:9800

; Peer 2
[Peer]
PublicKey = ttHzRDWUmVHWn+CXBGj04fYwdeb51wIUt0iC8ejP2wo=
AllowedIPs = 10.8.16.3/32
PersistentKeepalive = 20

[Peer]
PublicKey = 064r3zzmeaCGCEwXlfj+2tNV6tTnxbFiZalk1XIY7wI=
AllowedIPs = 10.8.16.4/32`)

	expected := Config{
		Interface: Interface{
			ListenPort: 51820,
			Address:    "10.8.16.1/24",
			PrivateKey: "4CwbPHW85Y/xdgB/zD/P0bZdM3XVNpi85H45FMscB1A=",
			DNS:        []string{"1.1.1.1", "1.1.0.0"},
			MTU:        1500,
			PreUp:      "echo \"UP\"",
			PostUp:     "iptables -A FORWARD -i %i -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE",
			PreDown:    "echo \"DOWN\"",
			PostDown:   "iptables -D FORWARD -i %i -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE",
		},
		Peers: Peers{
			{
				Comment:    "Peer 1",
				PublicKey:  "HHQSHN5TG6d0f3Wo0zeJM74v6073rQhc1+Yc8cwQ32Q=",
				AllowedIPs: "10.8.16.2/32",
				Endpoint:   "https://example.com:9800",
			},
			{
				Comment:             "Peer 2",
				PublicKey:           "ttHzRDWUmVHWn+CXBGj04fYwdeb51wIUt0iC8ejP2wo=",
				AllowedIPs:          "10.8.16.3/32",
				PersistentKeepalive: 20,
			},
			{
				PublicKey:  "064r3zzmeaCGCEwXlfj+2tNV6tTnxbFiZalk1XIY7wI=",
				AllowedIPs: "10.8.16.4/32",
			},
		},
	}

	cfg := Config{}
	err := cfg.Read(bytes.NewBuffer(input))

	assert.NoError(t, err)
	assert.Equal(t, expected, cfg)
}

func TestConfig_ReadFile(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "wgconfig-go-ini")
	assert.NoError(t, err)
	defer os.Remove(file.Name())

	_, err = file.Write([]byte(`[Interface]
	Address = 10.8.16.1/24
	ListenPort = 51820
	PrivateKey = 4CwbPHW85Y/xdgB/zD/P0bZdM3XVNpi85H45FMscB1A=`))
	assert.NoError(t, err)

	cfg := Config{}
	cfg.ReadFile(file.Name())

	expected := Config{
		Interface: Interface{
			ListenPort: 51820,
			Address:    "10.8.16.1/24",
			PrivateKey: "4CwbPHW85Y/xdgB/zD/P0bZdM3XVNpi85H45FMscB1A=",
		},
	}

	assert.Equal(t, expected, cfg)
}

func TestConfig_ReadInlineComment(t *testing.T) {
	input := []byte(`[Interface]
PostUp = iptables -A FORWARD -i %i -j ACCEPT ; comment
`)

	expected := Config{
		Interface: Interface{
			PostUp: "iptables -A FORWARD -i %i -j ACCEPT",
		},
	}
	cfg := Config{}
	err := cfg.Read(bytes.NewBuffer(input))

	assert.NoError(t, err)
	assert.Equal(t, expected, cfg)
}

func TestConfig_ExportJSON(t *testing.T) {
	var cfg Config
	cfg.Interface = Interface{
		Address:    "10.10.10.4/32",
		ListenPort: 5670,
		MTU:        1420,
		DNS:        []string{"8.8.4.4", "8.8.8.8"},
		PrivateKey: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX=",
		Table:      25,
		PreUp:      "pre-up",
		PostUp:     "post-up",
		PreDown:    "pre-down",
		PostDown:   "post-down",
	}

	cfg.Peers.Add(&Peer{
		Comment:             "comment",
		PublicKey:           "YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY=",
		PresharedKey:        "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ=",
		AllowedIPs:          "10.10.10.8/32",
		Endpoint:            "endpoint",
		PersistentKeepalive: 40,
	})

	data, err := json.MarshalIndent(cfg, "", "  ")
	assert.NoError(t, err)

	expected := []byte(`{
  "interface": {
    "private_key": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX=",
    "address": "10.10.10.4/32",
    "listen_port": 5670,
    "dns": [
      "8.8.4.4",
      "8.8.8.8"
    ],
    "table": 25,
    "mtu": 1420,
    "pre_up": "pre-up",
    "post_up": "post-up",
    "pre_down": "pre-down",
    "post_down": "post-down"
  },
  "peers": [
    {
      "comment": "comment",
      "public_key": "YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY=",
      "preshared_key": "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ=",
      "allowed_ips": "10.10.10.8/32",
      "endpoint": "endpoint",
      "persistent_keepalive": 40
    }
  ]
}`)
	assert.Equal(t, expected, data)
}

func TestConfig_ImportJSON(t *testing.T) {
	data := []byte(`{
"interface": {
	"private_key": "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX=",
	"address": "10.85.22.1/24",
	"listen_port": 5670,
	"dns": [ "1.1.1.1", "1.1.0.0" ],
	"table": 9842,
	"mtu": 2400,
	"pre_up": "pre-up",
	"post_up": "post-up",
	"pre_down": "pre-down",
	"post_down": "post-down"
},
"peers": [
	{
		"comment": "comment",
		"public_key": "YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY=",
		"preshared_key": "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ=",
		"allowed_ips": "10.85.22.40/32",
		"endpoint": "endpoint",
		"persistent_keepalive": 40
	},
	{
		"public_key": "EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE=",
		"preshared_key": "TTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTT=",
		"allowed_ips": "10.85.22.45/32",
		"endpoint": "endpoint2",
		"persistent_keepalive": 35
	}
]
}`)

	var cfg Config

	err := json.Unmarshal(data, &cfg)
	assert.NoError(t, err)

	expected := Config{
		Interface: Interface{
			PrivateKey: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX=",
			Address:    "10.85.22.1/24",
			ListenPort: 5670,
			DNS:        []string{"1.1.1.1", "1.1.0.0"},
			Table:      9842,
			MTU:        2400,
			PreUp:      "pre-up",
			PreDown:    "pre-down",
			PostUp:     "post-up",
			PostDown:   "post-down",
		},

		Peers: Peers{
			Peer{
				Comment:             "comment",
				PublicKey:           "YYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYYY=",
				PresharedKey:        "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ=",
				AllowedIPs:          "10.85.22.40/32",
				Endpoint:            "endpoint",
				PersistentKeepalive: 40,
			},
			Peer{
				PublicKey:           "EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE=",
				PresharedKey:        "TTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTTT=",
				AllowedIPs:          "10.85.22.45/32",
				Endpoint:            "endpoint2",
				PersistentKeepalive: 35,
			},
		},
	}

	assert.Equal(t, expected, cfg)
}
