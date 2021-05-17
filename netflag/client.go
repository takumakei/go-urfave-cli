package netflag

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"

	"github.com/takumakei/go-delint"
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/urfave/cli/v2"
)

// Client represents the flags related to a client to connect a server.
type Client struct {
	// Name is the name of the Client, may be empty string.
	Name string

	// PredeterminedFlagNetwork is true if the value of FlagNetwork is predetermined by the option.
	PredeterminedFlagNetwork bool

	// FlagNetwork is the network to connect.
	FlagNetwork *cli.StringFlag

	// FlagAddress is the address to connect.
	FlagAddress *cli.StringFlag

	// FlagTLSCerts is the certificate filepath of the client.
	FlagTLSCerts *cli.StringSliceFlag

	// FlagTLSKeys is the private key filepath of the certificate.
	FlagTLSKeys *cli.StringSliceFlag

	// FlagTLSCAs is the certificate filepath of the RootCAs.
	FlagTLSCAs *cli.StringSliceFlag

	// FlagTLSServerName is the value of server name.
	FlagTLSServerName *cli.StringFlag

	// FlagTLSSkipVerify is value to InsecureSkipVerify.
	FlagTLSSkipVerify *cli.BoolFlag

	// FlagTLSMinVer is the minimum TLS version that is acceptable.
	FlagTLSMinVer *cli.GenericFlag

	// FlagTLSMaxVer is the maximum TLS version that is acceptable.
	FlagTLSMaxVer *cli.GenericFlag

	// FlagSet is clix.FlagSet.
	FlagSet clix.FlagSet
}

// NewClient returns NewClient(prefix, "", opts...).
func NewClient(prefix clix.FlagPrefix, opts ...Option) *Client {
	return NewClientName(prefix, "", opts...)
}

// NewClientName returns *Client.
func NewClientName(prefix clix.FlagPrefix, name string, opts ...Option) *Client {
	cfg := newConfig(opts...)

	var (
		nameNetwork       = clix.NewFlagNameAlias(prefix, name, "network", "net")
		nameAddress       = clix.NewFlagNameAlias(prefix, name, "address", "addr")
		nameTLSCert       = clix.NewFlagNameAlias(prefix, name, "tls-cert", "tlscrt")
		nameTLSKey        = clix.NewFlagNameAlias(prefix, name, "tls-cert-key", "tlskey")
		nameTLSCAs        = clix.NewFlagNameAlias(prefix, name, "tls-ca", "tlsca")
		nameTLSServerName = clix.NewFlagNameAlias(prefix, name, "tls-server-name", "tlssrv")
		nameTLSSkipVerify = clix.NewFlagNameAlias(prefix, name, "tls-skip-verify", "tlsinsecure")
		nameTLSMinVer     = clix.NewFlagNameAlias(prefix, name, "tls-min-version", "tlsmin")
		nameTLSMaxVer     = clix.NewFlagNameAlias(prefix, name, "tls-max-version", "tlsmax")
	)

	network := cfg.networkValue()

	return &Client{
		Name: name,

		PredeterminedFlagNetwork: cfg.networkPredetermined(),

		FlagNetwork: &cli.StringFlag{
			Name:        nameNetwork.Name,
			Aliases:     nameNetwork.Aliases,
			Usage:       "network to connect" + cfg.networkUsage(),
			EnvVars:     nameNetwork.EnvVars,
			FilePath:    nameNetwork.FilePath,
			Required:    cfg.networkRequired(),
			Value:       network,
			Destination: &network,
		},

		FlagAddress: &cli.StringFlag{
			Name:        nameAddress.Name,
			Aliases:     nameAddress.Aliases,
			Usage:       "address to connect",
			EnvVars:     nameAddress.EnvVars,
			FilePath:    nameAddress.FilePath,
			Required:    cfg.addressRequired(),
			Value:       cfg.address,
			Destination: new(string),
		},

		FlagTLSCerts: &cli.StringSliceFlag{
			Name:        nameTLSCert.Name,
			Aliases:     nameTLSCert.Aliases,
			Usage:       "certificate `file`",
			EnvVars:     nameTLSCert.EnvVars,
			FilePath:    nameTLSCert.FilePath,
			TakesFile:   true,
			Destination: &cli.StringSlice{},
		},

		FlagTLSKeys: &cli.StringSliceFlag{
			Name:        nameTLSKey.Name,
			Aliases:     nameTLSKey.Aliases,
			Usage:       "private key `file` of certificate",
			EnvVars:     nameTLSKey.EnvVars,
			FilePath:    nameTLSKey.FilePath,
			TakesFile:   true,
			Destination: &cli.StringSlice{},
		},

		FlagTLSCAs: &cli.StringSliceFlag{
			Name:        nameTLSCAs.Name,
			Aliases:     nameTLSCAs.Aliases,
			Usage:       "root CA `file` of server",
			EnvVars:     nameTLSCAs.EnvVars,
			FilePath:    nameTLSCAs.FilePath,
			TakesFile:   true,
			Destination: &cli.StringSlice{},
		},

		FlagTLSServerName: &cli.StringFlag{
			Name:        nameTLSServerName.Name,
			Aliases:     nameTLSServerName.Aliases,
			Usage:       "server name for verification",
			EnvVars:     nameTLSServerName.EnvVars,
			FilePath:    nameTLSServerName.FilePath,
			Destination: new(string),
		},

		FlagTLSSkipVerify: &cli.BoolFlag{
			Name:        nameTLSSkipVerify.Name,
			Aliases:     nameTLSSkipVerify.Aliases,
			Usage:       "TLS insecure skip verify",
			EnvVars:     nameTLSSkipVerify.EnvVars,
			FilePath:    nameTLSSkipVerify.FilePath,
			Value:       cfg.skipVerify,
			Destination: new(bool),
		},

		FlagTLSMinVer: &cli.GenericFlag{
			Name:     nameTLSMinVer.Name,
			Aliases:  nameTLSMinVer.Aliases,
			Usage:    "TLS minimum version",
			EnvVars:  nameTLSMinVer.EnvVars,
			FilePath: nameTLSMinVer.FilePath,
			Value:    NewTLSVersion(cfg.tlsMinVersion),
		},

		FlagTLSMaxVer: &cli.GenericFlag{
			Name:     nameTLSMaxVer.Name,
			Aliases:  nameTLSMaxVer.Aliases,
			Usage:    "TLS maximum version",
			EnvVars:  nameTLSMaxVer.EnvVars,
			FilePath: nameTLSMaxVer.FilePath,
			Value:    NewTLSVersion(cfg.tlsMaxVersion),
		},

		FlagSet: clix.NewFlagSet(),
	}
}

// Before calls f.FlagSet.Init(c).
// Before is intended to be used as cli.BeforeFunc.
func (f *Client) Before(c *cli.Context) error {
	delint.Must(f.FlagSet.Init(c))
	return nil
}

// Flags returns []cli.Flag.
//
//     f.FlagNetwork  (conditional)
//     f.FlagAddress
//     f.FlagTLSCerts
//     f.FlagTLSKeys
//     f.FlagTLSCAs
//     f.FlagTLSServerName
//     f.FlagTLSSkipVerify
//     f.FlagTLSMinVer
//     f.FlagTLSMaxVer
func (f *Client) Flags() []cli.Flag {
	return clix.Flags(
		clix.FlagIf(!f.PredeterminedFlagNetwork, f.FlagNetwork),
		f.FlagAddress,
		f.FlagTLSCerts,
		f.FlagTLSKeys,
		f.FlagTLSCAs,
		f.FlagTLSServerName,
		f.FlagTLSSkipVerify,
		f.FlagTLSMinVer,
		f.FlagTLSMaxVer,
	)
}

// Network returns the value of FlagNetwork.
func (f *Client) Network() string {
	return *f.FlagNetwork.Destination
}

// Address returns the value of FlagAddress.
func (f *Client) Address() string {
	return *f.FlagAddress.Destination
}

// TLSCerts returns the value of FlagTLSCerts.
func (f *Client) TLSCerts() []string {
	return f.FlagTLSCerts.Destination.Value()
}

// TLSKeys returns the value of FlagTLSKeys.
func (f *Client) TLSKeys() []string {
	return f.FlagTLSKeys.Destination.Value()
}

// TLSCAs returns the value of FlagTLSCAs.
func (f *Client) TLSCAs() []string {
	return f.FlagTLSCAs.Destination.Value()
}

func (f *Client) TLSServerName() string {
	return *f.FlagTLSServerName.Destination
}

// TLSSkipVerify returns the value of FlagTLSSkipVerify.
func (f *Client) TLSSkipVerify() bool {
	return *f.FlagTLSSkipVerify.Destination
}

// TLSMinVersion returns the value of FlagTLSMinVer.
func (f *Client) TLSMinVersion() uint16 {
	return f.FlagTLSMinVer.Value.(*TLSVersion).Value()
}

// TLSMaxVersion returns the value of FlagTLSMaxVer.
func (f *Client) TLSMaxVersion() uint16 {
	return f.FlagTLSMaxVer.Value.(*TLSVersion).Value()
}

// UseTLS returns true if TLS related flags are presented.
func (f *Client) UseTLS() bool {
	list := []cli.Flag{
		f.FlagTLSCerts,
		f.FlagTLSKeys,
		f.FlagTLSCAs,
		f.FlagTLSServerName,
		f.FlagTLSSkipVerify,
		f.FlagTLSMinVer,
		f.FlagTLSMaxVer,
	}
	for _, flag := range list {
		if f.FlagSet.IsSet(flag) {
			return true
		}
	}
	return false
}

// TLSConfig returns *tls.Config.
func (f *Client) TLSConfig() (*tls.Config, error) {
	var certs []tls.Certificate
	certFiles := f.TLSCerts()
	keyFiles := f.TLSKeys()
	if len(certFiles) != len(keyFiles) {
		if len(certFiles) < len(keyFiles) {
			return nil, fmt.Errorf("no certificate file for private key")
		}
		return nil, fmt.Errorf("no key file for certificate")
	}
	for i := range certFiles {
		cert, err := tls.LoadX509KeyPair(certFiles[i], keyFiles[i])
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)
	}

	var rootCAs *x509.CertPool
	if cas := f.TLSCAs(); len(cas) > 0 {
		rootCAs = x509.NewCertPool()
		for _, ca := range f.TLSCAs() {
			p, err := os.ReadFile(ca)
			if err != nil {
				return nil, fmt.Errorf("failed to read CA %q, %w", ca, err)
			}
			if ok := rootCAs.AppendCertsFromPEM(p); !ok {
				return nil, fmt.Errorf("failed to load CA %q", ca)
			}
		}
	}

	cfg := &tls.Config{
		Certificates:       certs,
		RootCAs:            rootCAs,
		InsecureSkipVerify: f.TLSSkipVerify(),
		MinVersion:         f.TLSMinVersion(),
		MaxVersion:         f.TLSMaxVersion(),
		ServerName:         f.TLSServerName(),
	}

	return cfg, nil
}

// Dial returns the result of calling f.DialNetwork(f.Network()).
func (f *Client) Dial() (net.Conn, error) {
	return f.DialNetwork(f.Network())
}

// DialNetwork returns the result of calling tls.Dial if f.UseTLS() returns
// true, otherwise returns the result of calling net.Dial.
// f.Address() and f.TLSConfig() are used.
func (f *Client) DialNetwork(network string) (net.Conn, error) {
	a := f.Address()
	if f.UseTLS() {
		cfg, err := f.TLSConfig()
		if err != nil {
			return nil, err
		}
		return tls.Dial(network, a, cfg)
	}
	return net.Dial(network, a)
}
