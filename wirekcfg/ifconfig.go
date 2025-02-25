package wirekcfg

import (
	"github.com/wirekcp/wgctrl/wgtypes"
)

type Config struct {
	// IPv4CIDR is the CIDR of the IPv4 address.
	IPv4CIDR   string
	ListenPort int
	PrivateKey string
}

func (c Config) ToWgConfig() *wgtypes.Config {
	key, _ := wgtypes.ParseKey(c.PrivateKey)
	return &wgtypes.Config{
		PrivateKey: &key,
		ListenPort: &c.ListenPort,
	}
}
