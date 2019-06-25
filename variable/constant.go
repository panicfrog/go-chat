package variable

import "strings"

type Platform string

func StoP(s string) (p Platform) {
	if s == strings.ToLower(string(PlatformIOS)) {
		return PlatformIOS
	} else if s == strings.ToLower(string((PlatformAndroid))) {
		return PlatformAndroid
	} else if s == strings.ToLower(string(PlatformWeb)) {
		return PlatformWeb
	} else if s == strings.ToLower(string(PlatformDesktop)) {
		return PlatformDesktop
	} else {
		return PlatformUnknow
	}
}

const (
	PlatformIOS     Platform = "iOS"
	PlatformAndroid Platform = "android"
	PlatformWeb     Platform = "web"
	PlatformDesktop Platform = "desktop"
	PlatformUnknow  Platform = "unknow"
)

const (
	AuthKey = "user"
	PlatformKey = "platform"
)