package netflag

import (
	"crypto/elliptic"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"

	"github.com/takumakei/go-cert4now"
	"github.com/takumakei/go-delint"
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/urfave/cli/v2"
)

// Server represents the flags related to a server to listen and accept clients.
type Server struct {
	// Name is the name of Server, may be empty string.
	Name string

	// PredeterminedFlagNetwork is true if the value of FlagNetwork is predetermined by the option.
	PredeterminedFlagNetwork bool

	// DisableTLS is true if TLS is disabled.
	DisableTLS bool

	// DisableFlagTLSGenCert is true if FlagTLSGenCert would be included in the result of Flags().
	DisableFlagTLSGenCert bool

	// FlagNetwork is the network to listen.
	FlagNetwork *cli.StringFlag

	// FlagAddress is the address to listen.
	FlagAddress *cli.StringFlag

	// FlagTLSCerts is the certificate filepath of the server.
	FlagTLSCerts *cli.StringSliceFlag

	// FlagTLSKeys is the private key filepath of the certificate.
	FlagTLSKeys *cli.StringSliceFlag

	// FlagTLSGenCert specifies whether to generate a self signed certificate.
	FlagTLSGenCert *cli.BoolFlag

	// FlagTLSCAs is the certificate filepath of the ClientCAs.
	FlagTLSCAs *cli.StringSliceFlag

	// FlagTLSMinVer is the minimum TLS version that is acceptable.
	FlagTLSMinVer *cli.GenericFlag

	// FlagTLSMaxVer is the maximum TLS version that is acceptable.
	FlagTLSMaxVer *cli.GenericFlag

	// FlagSet is clix.FlagSet.
	FlagSet clix.FlagSet
}

// NewServer returns NewServer(prefix, "", opts...).
func NewServer(prefix clix.FlagPrefix, opts ...Option) *Server {
	return NewServerName(prefix, "", opts...)
}

// NewServerName returns *Server.
func NewServerName(prefix clix.FlagPrefix, name string, opts ...Option) *Server {
	cfg := newConfig(opts...)

	var (
		nameNetwork    = clix.NewFlagNameAlias(prefix, name, "network", "net")
		nameAddress    = clix.NewFlagNameAlias(prefix, name, "address", "addr")
		nameTLSCert    = clix.NewFlagNameAlias(prefix, name, "tls-cert", "tlscrt")
		nameTLSKey     = clix.NewFlagNameAlias(prefix, name, "tls-cert-key", "tlskey")
		nameTLSGenCert = clix.NewFlagNameAlias(prefix, name, "tls-gen-cert", "tlsgen")
		nameTLSCAs     = clix.NewFlagNameAlias(prefix, name, "tls-ca", "tlsca")
		nameTLSMinVer  = clix.NewFlagNameAlias(prefix, name, "tls-min-version", "tlsmin")
		nameTLSMaxVer  = clix.NewFlagNameAlias(prefix, name, "tls-max-version", "tlsmax")
	)

	network := cfg.networkValue()

	flagTLSCert := &cli.StringSliceFlag{
		Name:        nameTLSCert.Name,
		Aliases:     nameTLSCert.Aliases,
		Usage:       "certificate `file`",
		EnvVars:     nameTLSCert.EnvVars,
		FilePath:    nameTLSCert.FilePath,
		TakesFile:   true,
		Destination: &cli.StringSlice{},
	}

	flagTLSKeys := &cli.StringSliceFlag{
		Name:        nameTLSKey.Name,
		Aliases:     nameTLSKey.Aliases,
		Usage:       "private key `file` of certificate",
		EnvVars:     nameTLSKey.EnvVars,
		FilePath:    nameTLSKey.FilePath,
		TakesFile:   true,
		Destination: &cli.StringSlice{},
	}

	flagTLSGenCert := &cli.BoolFlag{
		Name:        nameTLSGenCert.Name,
		Aliases:     nameTLSGenCert.Aliases,
		Usage:       "generate self-signed certificate",
		EnvVars:     nameTLSGenCert.EnvVars,
		FilePath:    nameTLSGenCert.FilePath,
		Destination: new(bool),
	}

	flagSet := clix.NewFlagSet()

	return &Server{
		Name: name,

		PredeterminedFlagNetwork: cfg.networkPredetermined(),

		DisableTLS: cfg.tlsDisabled,

		DisableFlagTLSGenCert: cfg.genCertDisabled,

		FlagNetwork: &cli.StringFlag{
			Name:        nameNetwork.Name,
			Aliases:     nameNetwork.Aliases,
			Usage:       "network to listen" + cfg.networkUsage(),
			EnvVars:     nameNetwork.EnvVars,
			FilePath:    nameNetwork.FilePath,
			Required:    cfg.networkRequired(),
			Value:       network,
			Destination: &network,
		},

		FlagAddress: &cli.StringFlag{
			Name:        nameAddress.Name,
			Aliases:     nameAddress.Aliases,
			Usage:       "address to listen",
			EnvVars:     nameAddress.EnvVars,
			FilePath:    nameAddress.FilePath,
			Required:    cfg.addressRequired(),
			Value:       cfg.address,
			Destination: new(string),
		},

		FlagTLSCerts: flagTLSCert,

		FlagTLSKeys: flagTLSKeys,

		FlagTLSGenCert: flagTLSGenCert,

		FlagTLSCAs: &cli.StringSliceFlag{
			Name:        nameTLSCAs.Name,
			Aliases:     nameTLSCAs.Aliases,
			Usage:       "root CA `file` for client auth",
			EnvVars:     nameTLSCAs.EnvVars,
			FilePath:    nameTLSCAs.FilePath,
			TakesFile:   true,
			Destination: &cli.StringSlice{},
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

		FlagSet: flagSet,
	}
}

// Before calls f.FlagSet.Init(c), validates exclusive flags.
// Before is intended to be used as cli.BeforeFunc.
func (f *Server) Before(c *cli.Context) error {
	delint.Must(f.FlagSet.Init(c))
	if f.FlagSet.IsSet(f.FlagTLSCerts) && f.TLSGenCert() {
		return fmt.Errorf(
			"%q and %q must not set at the same time",
			f.FlagTLSCerts.Name,
			f.FlagTLSGenCert.Name,
		)
	}
	if f.FlagSet.IsSet(f.FlagTLSKeys) && f.TLSGenCert() {
		return fmt.Errorf(
			"%q and %q must not set at the same time",
			f.FlagTLSKeys.Name,
			f.FlagTLSGenCert.Name,
		)
	}
	return nil
}

// Flags returns []cli.Flag.
//
// It includes the following.
//
//     f.FlagNetwork  (if not predetermined)
//     f.FlagAddress
//
// It also includes the following if TLS is enabled.
//
//     f.FlagTLSCerts
//     f.FlagTLSKeys
//     f.FlagTLSGenCert  (if not disabled)
//     f.FlagTLSCAs
//     f.FlagTLSMinVer
//     f.FlagTLSMaxVer
func (f *Server) Flags() []cli.Flag {
	return clix.Flags(
		clix.FlagIf(!f.PredeterminedFlagNetwork, f.FlagNetwork),
		f.FlagAddress,
		clix.FlagIf(!f.DisableTLS, clix.Flags(
			f.FlagTLSCerts,
			f.FlagTLSKeys,
			clix.FlagIf(!f.DisableFlagTLSGenCert, f.FlagTLSGenCert),
			f.FlagTLSCAs,
			f.FlagTLSMinVer,
			f.FlagTLSMaxVer,
		)...),
	)
}

// Network returns the value of FlagNetwork.
func (f *Server) Network() string {
	return *f.FlagNetwork.Destination
}

// Address returns the value of FlagAddress.
func (f *Server) Address() string {
	return *f.FlagAddress.Destination
}

// TLSCerts returns the value of FlagTLSCerts.
func (f *Server) TLSCerts() []string {
	return f.FlagTLSCerts.Destination.Value()
}

// TLSKeys returns the value of FlagTLSKeys.
func (f *Server) TLSKeys() []string {
	return f.FlagTLSKeys.Destination.Value()
}

// TLSGenCert returns the value of FlagTLSGenCert.
func (f *Server) TLSGenCert() bool {
	return *f.FlagTLSGenCert.Destination
}

// TLSCAs returns the value of FlagTLSCAs.
func (f *Server) TLSCAs() []string {
	return f.FlagTLSCAs.Destination.Value()
}

// TLSMinVersion returns the value of FlagTLSMinVer.
func (f *Server) TLSMinVersion() uint16 {
	return f.FlagTLSMinVer.Value.(*TLSVersion).Value()
}

// TLSMaxVersion returns the value of FlagTLSMaxVer.
func (f *Server) TLSMaxVersion() uint16 {
	return f.FlagTLSMaxVer.Value.(*TLSVersion).Value()
}

// UseTLS returns true if TLS related flags are presented.
func (f *Server) UseTLS() bool {
	list := []cli.Flag{
		f.FlagTLSCerts,
		f.FlagTLSKeys,
		f.FlagTLSCAs,
		f.FlagTLSGenCert,
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
func (f *Server) TLSConfig() (*tls.Config, error) {
	var certs []tls.Certificate
	if f.TLSGenCert() {
		names := make([]string, 0, 2)
		if v := f.Name; len(v) > 0 {
			names = append(names, v)
		}
		if v, err := os.Hostname(); err == nil && len(v) > 0 {
			names = append(names, v)
		}
		cert, err := cert4now.Generate(
			cert4now.CommonName(f.Name),
			cert4now.ECDSA(elliptic.P384()),
			cert4now.AddDate(100, 0, 0),
			cert4now.DNSNames(names...),
		)
		if err != nil {
			return nil, err
		}
		certs = []tls.Certificate{cert}
	} else {
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
	}

	clientAuth := tls.NoClientCert
	var clientCAs *x509.CertPool
	if cas := f.TLSCAs(); len(cas) > 0 {
		clientCAs = x509.NewCertPool()
		for _, ca := range f.TLSCAs() {
			p, err := os.ReadFile(ca)
			if err != nil {
				return nil, fmt.Errorf("failed to read CA %q, %w", ca, err)
			}
			if ok := clientCAs.AppendCertsFromPEM(p); !ok {
				return nil, fmt.Errorf("failed to load CA %q", ca)
			}
		}
		clientAuth = tls.RequireAndVerifyClientCert
	}

	cfg := &tls.Config{
		Certificates: certs,
		ClientAuth:   clientAuth,
		ClientCAs:    clientCAs,
		MinVersion:   f.TLSMinVersion(),
		MaxVersion:   f.TLSMaxVersion(),
	}

	return cfg, nil
}

// Listen returns the result of calling f.ListenNetwork(f.Network()).
func (f *Server) Listen() (net.Listener, error) {
	return f.ListenNetwork(f.Network())
}

// ListenNetwork returns the result of calling tls.Listen if f.UseTLS() returns
// true, otherwise returns the result of calling net.Listen.
// f.Address() and f.TLSConfig() are used.
func (f *Server) ListenNetwork(network string) (net.Listener, error) {
	a := f.Address()
	if f.UseTLS() {
		cfg, err := f.TLSConfig()
		if err != nil {
			return nil, err
		}
		return tls.Listen(network, a, cfg)
	}
	return net.Listen(network, a)
}
