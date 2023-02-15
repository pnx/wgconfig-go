package wgconfig

import (
	"io"

	"gopkg.in/ini.v1"
)

type Config struct {
	Interface Interface `json:"interface,omitempty"`
	Peers     Peers     `json:"peers,omitempty"`
}

func (cfg *Config) AddPeer(p *Peer) {
	cfg.Peers.Add(p)
}

func (cfg *Config) load(source interface{}) error {
	opt := ini.LoadOptions{
		SkipUnrecognizableLines:  true,
		AllowNonUniqueSections:   true,
		SpaceBeforeInlineComment: true,
	}

	f, err := ini.LoadSources(opt, source)
	if err != nil {
		return err
	}

	err = f.Section("Interface").MapTo(&cfg.Interface)
	if err != nil {
		return err
	}

	peers, err := f.SectionsByName("Peer")
	if err != nil {
		// No "Peer" sections is not an error. so just return here.
		return nil
	}

	for _, p := range peers {
		peer := Peer{}
		err = p.MapTo(&peer)
		if err != nil {
			return err
		}
		peer.Comment = trimComment(p.Comment)
		cfg.AddPeer(&peer)
	}

	return nil
}

func (cfg *Config) ReadFile(file string) error {
	return cfg.load(file)
}

func (cfg *Config) Read(r io.Reader) error {
	return cfg.load(r)
}

// Write config to file
func (cfg *Config) Write(w io.Writer) (int64, error) {
	file := ini.Empty(ini.LoadOptions{
		AllowNonUniqueSections: true,
		IgnoreInlineComment:    true,
	})

	err := file.ReflectFrom(&cfg)
	if err != nil {
		return -1, err
	}

	return file.WriteTo(w)
}
