//go:build windows

package wirektun

import (
	"fmt"
	"net"
	"os/exec"

	"github.com/wirekcp/wireguard-go/tun"
	"golang.org/x/sys/windows"
	"golang.zx2c4.com/wireguard/windows/tunnel/winipcfg"
)

func init() {
	tun.WintunTunnelType = "WireKCP"
	guid, err := windows.GUIDFromString("{42b30cad-f96c-4369-9094-47a0d68cd40f}")
	if err != nil {
		panic(err)
	}
	tun.WintunStaticRequestedGUID = &guid
	// ipc.UAPISecurityDescriptor, err = windows.SecurityDescriptorFromString("O:SYD:P(A;;GA;;;SY)(A;;GA;;;BA)S:(ML;;NWNRNX;;;HI)")
	// if err != nil {
	// 	panic(err)
	// }
}

func interfaceName(dev tun.Device) (string, error) {
	guid, err := winipcfg.LUID(dev.(*tun.NativeTun).LUID()).GUID()
	if err != nil {
		return "", err
	}
	return guid.String(), nil
}

func SetIPwithoutTun(cidr string) error {
	ip, network, _ := net.ParseCIDR(cidr)
	mask := network.Mask
	subnet := fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3])
	cmd := exec.Command("netsh", "interface", "ip", "set", "address", "name="+DefaultTunName(), "static", ip.String(), subnet)
	return cmd.Run()
}
