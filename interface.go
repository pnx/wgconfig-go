package wgconfig

type Interface struct {
	PrivateKey string   `ini:"PrivateKey,omitempty" json:"private_key,omitempty"`
	Address    string   `ini:"Address,omitempty" json:"address,omitempty"`
	ListenPort uint     `ini:"ListenPort,omitempty" json:"listen_port,omitempty"`
	DNS        []string `ini:"DNS,omitempty" json:"dns,omitempty"`
	Table      int      `ini:"Table,omitempty" json:"table,omitempty"` // Routing table.
	MTU        int      `ini:"MTU,omitempty" json:"mtu,omitempty"`
	PreUp      string   `ini:"PreUp,omitempty" json:"pre_up,omitempty"`
	PostUp     string   `ini:"PostUp,omitempty" json:"post_up,omitempty"`
	PreDown    string   `ini:"PreDown,omitempty" json:"pre_down,omitempty"`
	PostDown   string   `ini:"PostDown,omitempty" json:"post_down,omitempty"`
}
