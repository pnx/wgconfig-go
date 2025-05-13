package wgconfig

import (
	"gopkg.in/ini.v1"
)

type Peer struct {
	Comment             string `ini:"-" json:"comment,omitempty"`
	PublicKey           string `ini:"PublicKey,omitempty" json:"public_key,omitempty"`
	PresharedKey        string `ini:"PresharedKey,omitempty" json:"preshared_key,omitempty"`
	AllowedIPs          string `ini:"AllowedIPs,omitempty" json:"allowed_ips,omitempty"`
	Endpoint            string `ini:"Endpoint,omitempty" json:"endpoint,omitempty"`
	PersistentKeepalive int    `ini:"PersistentKeepalive,omitempty" json:"persistent_keepalive,omitempty"`
}

type Peers []Peer

// Add peers to the list.
func (p *Peers) Add(peer *Peer) {
	*p = append(*p, *peer)
}

// There is a hidden (e.g. not so well documented) "marshal" interface
// called StructReflector (https://pkg.go.dev/github.com/go-ini/ini#StructReflector)
// So we can define one for "Peers" (slice of Peer) type here.
// This way. the ini library is able to marshal our struct.
//
// NOTE: This should really exists in the library already as all this function
// do is call ReflectINIStruct() on it's elements.
func (peers Peers) ReflectINIStruct(f *ini.File) error {
	for _, p := range peers {
		err := p.ReflectINIStruct(f)
		if err != nil {
			return err
		}
	}
	return nil
}

// Bit of a hack here.
// The ini library overwrites sections with the same name when ReflectFrom is used
// and there is no other way around it (that is known atleast).
// So, We need to be abit cleaver here and first create a new empty ini.File object and perform ReflectFrom there.
//
// See: https://github.com/go-ini/ini/blob/842b9a946dd3bb45b95001ce6d9ae86f3b3ce921/struct_test.go#L663
func (p Peer) ReflectINIStruct(f *ini.File) error {
	s, err := f.NewSection("Peer")
	if err != nil {
		return err
	}

	// Perform reflectFrom on a temporary section object.
	tmp := ini.Empty().Section("")
	err = tmp.ReflectFrom(&p)
	if err != nil {
		return err
	}

	// Then we can copy the data to the real section object.
	*s = *tmp

	// Dont forget the comment :)
	s.Comment = p.Comment
	return nil
}
