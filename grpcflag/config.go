package grpcflag

import (
	"crypto/tls"
	"strings"
)

type Config struct {
	Network       []string
	Address       string
	TLSMinVersion uint16
	TLSMaxVersion uint16
}

func NewConfig(opt ...Option) Config {
	cfg := Config{
		TLSMinVersion: tls.VersionTLS12,
		TLSMaxVersion: tls.VersionTLS13,
	}
	cfg.Apply(opt)
	return cfg
}

func (c *Config) Apply(opts []Option) {
	for _, opt := range opts {
		opt(c)
	}
}

func (c *Config) NetworkPredetermined() bool {
	return len(c.Network) == 0 || (len(c.Network) == 1 && c.Network[0] != "*")
}

func (c *Config) NetworkRequired() bool {
	return len(c.Network) == 1 && c.Network[0] == "*"
}

func (c *Config) NetworkUsage() string {
	switch len(c.Network) {
	case 0:
		return ""
	case 1:
		v := c.Network[0]
		if v == "*" {
			return ""
		}
		return v
	}
	return " [" + strings.Join(c.Network, "|") + "]"
}

func (c *Config) NetworkValue() string {
	if len(c.Network) > 0 {
		v := c.Network[0]
		if v == "*" {
			return ""
		}
		return v
	}
	return "tcp"
}
