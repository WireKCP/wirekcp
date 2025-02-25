package wirektypes

import "github.com/wirekcp/wgctrl/wgtypes"

func GeneratePrivateKey() string {
	key, _ := wgtypes.GeneratePrivateKey()
	return key.String()
}

func ParseKey(key string) (wgtypes.Key, error) {
	return wgtypes.ParseKey(key)
}
