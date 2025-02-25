package wirekcfg

import (
	"github.com/wirekcp/wgctrl"
	"github.com/wirekcp/wireguard-go/conn"
	"github.com/wirekcp/wireguard-go/device"
	"github.com/wirekcp/wireguard-go/tun"
)

// NewDevice returns a wireguard-go Device configured for WireKCP use.
func NewDevice(tunDev tun.Device, bind conn.Bind, logger *device.Logger) *device.Device {
	ret := device.NewDevice(tunDev, bind, logger)
	ret.DisableSomeRoamingForBrokenMobileSemantics()
	return ret
}

func ConfigureDevice(tun string, config Config) error {
	client, err := wgctrl.New()
	if err != nil {
		return err
	}
	return client.ConfigureDevice(tun, *config.ToWgConfig())
}
