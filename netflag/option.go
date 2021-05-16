package netflag

import (
	"crypto/tls"
	"strings"

	"github.com/takumakei/go-stringx"
	"github.com/urfave/cli/v2"
)

// config represents a state of options.
type config struct {
	// network
	//
	//     {}: defined, "tcp"
	//     {"udp"}: defined, "udp"
	//     {"*"}: required
	//     {"udp","*"}: optional, "udp" default
	//     {"tcp","unix"}: optional, "tcp" default
	network []string

	// address is a default address
	address string

	skipVerify    bool
	tlsMinVersion uint16
	tlsMaxVersion uint16

	genCertDisabled bool
}

func newConfig(opts ...Option) config {
	var cfg config
	// reasonable default
	cfg.apply([]Option{
		TLSMinVersion(tls.VersionTLS12),
		TLSMaxVersion(tls.VersionTLS13),
	})
	cfg.apply(opts)
	return cfg
}

// apply applies opts to c.
func (c *config) apply(opts []Option) {
	for _, v := range opts {
		v(c)
	}
}

// networkPredetermined returns true if the value of FlagNetwork can not be changed.
func (c *config) networkPredetermined() bool {
	return len(c.network) == 0 || (len(c.network) == 1 && c.network[0] != "*")
}

// networkRequired returns true if FlagNetwork is mandatory option because of
// lack of default.
func (c *config) networkRequired() bool {
	return len(c.network) == 1 && c.network[0] == "*"
}

// networkUsage returns the Usage for FlagNetwork.
func (c *config) networkUsage() string {
	if stringx.Index(c.network, "*") != -1 {
		return ""
	}
	if len(c.network) >= 2 {
		return " `[" + strings.Join(c.network, "|") + "]`"
	}
	return ""
}

// networkValue returns the initial Value of FlagNetwork.
func (c *config) networkValue() string {
	if len(c.network) == 0 {
		return "tcp"
	}
	if c.network[0] != "*" {
		return c.network[0]
	}
	return ""
}

// addressRequired returns true if FlagAddress is mandatory option because of
// lack of default.
func (c *config) addressRequired() bool {
	return len(c.address) == 0
}

// Option represents options for Client and Server.
type Option func(*config)

// Address returns the option setting default address.
// FlagAddress is set required if the default address is not set.
func Address(a string) Option {
	return func(c *config) {
		c.address = a
	}
}

// Network returns the option shows acceptable networks.
// FlagNetwork is available in command line options if more than one networks
// is set.  Otherwise FlagNetwork is not shown in the command line options.
// Default value of network is "tcp".
func Network(net ...string) Option {
	return func(c *config) {
		c.network = net
	}
}

// SkipVerify returns the option to set default value of FlagSkipVerify.
func SkipVerify(v bool) Option {
	return func(c *config) {
		c.skipVerify = v
	}
}

// TLSMinVersion returns the option to set default value of FlagTLSMinVer.
func TLSMinVersion(v uint16) Option {
	return func(c *config) {
		c.tlsMinVersion = v
	}
}

// TLSMaxVersion returns the option to set default value of FlagTLSMaxVer.
func TLSMaxVersion(v uint16) Option {
	return func(c *config) {
		c.tlsMaxVersion = v
	}
}

// GenCert returns the option whether using self-signed certificate.
func GenCert(v bool) Option {
	if v {
		return EnableGenCert
	}
	return DisableGenCert
}

// EnableGenCert is the option to use FlagTLSGenCert.
func EnableGenCert(c *config) {
	c.genCertDisabled = false
}

// DisableGenCert is the option not to use FlagTLSGenCert.
func DisableGenCert(c *config) {
	c.genCertDisabled = true
}

func useFlagIf(v bool, flag cli.Flag) []cli.Flag {
	if v {
		return []cli.Flag{flag}
	}
	return nil
}
