package grpcflag

type Option func(*Config)

func Network(network ...string) Option {
	network = removeEmptyString(network)
	return func(c *Config) {
		c.Network = network
	}
}

func removeEmptyString(s []string) []string {
	for i, v := range s {
		if len(v) == 0 {
			t := s[:i]
			for _, v := range s[i+1:] {
				if len(v) > 0 {
					t = append(t, v)
				}
			}
			return t
		}
	}
	return s
}

func Address(address string) Option {
	return func(c *Config) {
		c.Address = address
	}
}

func TLSMinVersion(v uint16) Option {
	return func(c *Config) {
		c.TLSMinVersion = v
	}
}

func TLSMaxVersion(v uint16) Option {
	return func(c *Config) {
		c.TLSMaxVersion = v
	}
}
