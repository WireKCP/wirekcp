package frontend

import (
	"errors"
	"net"
	"net/netip"
	"strconv"
	"strings"

	"github.com/wirekcp/wgctrl"
	"github.com/wirekcp/wgctrl/wgtypes"
)

func ValidateDeviceExists(client *wgctrl.Client, name string) error {
	device, err := client.Device(name)
	if device == nil {
		return errors.New("device does not exist\n")
	}
	return err
}

func ValidateRequiredInt(s string) error {
	if err := ValidateRequiredString(s); err != nil {
		return err
	}
	if _, err := strconv.ParseUint(s, 10, 16); err != nil {
		return err
	}
	return nil
}

func ValidateIPAddress(ip string) error {
	_, err := netip.ParseAddr(ip)
	return err
}

func ValidateRequiredCIDRs(cidrs string) error {
	if err := ValidateRequiredString(cidrs); err != nil {
		return err
	}
	return ValidateCIDRs(cidrs)
}

func ValidateRequiredCIDR(cidr string) error {
	if err := ValidateRequiredString(cidr); err != nil {
		return err
	}
	return ValidateCIDR(cidr)
}

func ValidateCIDRs(cidrs string) error {
	cidr := strings.Split(cidrs, ",")
	for _, c := range cidr {
		c := strings.TrimSpace(c)
		if err := ValidateCIDR(c); err != nil {
			return err
		}
	}
	return nil
}

func ValidateCIDR(cidr string) error {
	_, _, err := net.ParseCIDR(cidr)
	return err
}

func ValidateRequiredUDPAddr(cidr string) error {
	if err := ValidateRequiredString(cidr); err != nil {
		return err
	}
	return ValidateUDPAddr(cidr)
}

func ValidateOptionalUDPAddr(cidr string) error {
	if cidr == "" {
		return nil
	}
	return ValidateUDPAddr(cidr)
}

func ValidateUDPAddr(cidr string) error {
	_, err := net.ResolveUDPAddr("udp", cidr)
	return err
}

func ValidateRequiredKey(s string) error {
	if err := ValidateRequiredString(s); err != nil {
		return err
	}
	if _, err := wgtypes.ParseKey(s); err != nil {
		return err
	}
	return nil
}

func ValidateOptionalKey(s string) error {
	if s == "" {
		return nil
	}
	return ValidateRequiredKey(s)
}

func ValidateRequiredString(s string) error {
	if s == "" {
		return errors.New("this field is required")
	}
	return nil
}
