package netflag

import (
	"crypto/tls"
	"fmt"

	"github.com/urfave/cli/v2"
)

// TLSVersion wraps a uint16 as tls.Version* to satisfy flag.Value.
type TLSVersion struct {
	ver uint16
}

// Note: cli.Generic is equivalent to flag.Value.
var _ cli.Generic = (*TLSVersion)(nil)

// NewTLSVersion creates a *TLSVersion with a default value.
func NewTLSVersion(value uint16) *TLSVersion {
	return &TLSVersion{ver: value}
}

// Set parses value as TLS version string, sets it.
func (tv *TLSVersion) Set(value string) error {
	switch value {
	case "1.0":
		tv.ver = tls.VersionTLS10
	case "1.1":
		tv.ver = tls.VersionTLS11
	case "1.2":
		tv.ver = tls.VersionTLS12
	case "1.3":
		tv.ver = tls.VersionTLS13
	default:
		return fmt.Errorf("%s is not a TLS version", value)
	}
	return nil
}

// String returns a readable representation of this value (for usage defaults)
func (tv *TLSVersion) String() string {
	switch tv.ver {
	case tls.VersionTLS10:
		return "1.0"
	case tls.VersionTLS11:
		return "1.1"
	case tls.VersionTLS12:
		return "1.2"
	case tls.VersionTLS13:
		return "1.3"
	}
	return ""
}

// Value returns an uint16 as TLS version set by this flag.
func (tv *TLSVersion) Value() uint16 {
	return tv.ver
}
