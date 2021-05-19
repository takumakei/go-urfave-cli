package grpcflag

import (
	"crypto/elliptic"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"github.com/takumakei/go-cert4now"
	"github.com/takumakei/go-delint"
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	Name string

	PredeterminedNetwork bool

	FlagNetwork *cli.StringFlag
	FlagAddress *cli.StringFlag

	FlagTLSCerts    *cli.StringSliceFlag
	FlagTLSCertKeys *cli.StringSliceFlag

	FlagTLSGenCert *cli.BoolFlag

	FlagTLSVerifyClient *cli.BoolFlag
	FlagTLSClientCAs    *cli.StringSliceFlag

	FlagTLSMinVersion *cli.GenericFlag
	FlagTLSMaxVersion *cli.GenericFlag

	FlagSet clix.FlagSet
}

func NewServer(prefix clix.FlagPrefix, opt ...Option) *Server {
	return NewServerName(prefix, "", opt...)
}

func NewServerName(prefix clix.FlagPrefix, name string, opt ...Option) *Server {
	cfg := NewConfig(opt...)

	var (
		nameFlagNetwork         = clix.NewFlagNameAlias(prefix, name, "network", "net")
		nameFlagAddress         = clix.NewFlagNameAlias(prefix, name, "address", "addr")
		nameFlagTLSCerts        = clix.NewFlagNameAlias(prefix, name, "tls-cert", "tlscrt")
		nameFlagTLSCertKeys     = clix.NewFlagNameAlias(prefix, name, "tls-cert-key", "tlskey")
		nameFlagTLSGenCert      = clix.NewFlagNameAlias(prefix, name, "tls-gen-cert", "tlsgen")
		nameFlagTLSVerifyClient = clix.NewFlagNameAlias(prefix, name, "tls-verify-client", "mtls")
		nameFlagTLSClientCAs    = clix.NewFlagNameAlias(prefix, name, "tls-client-ca", "tlsca")
		nameFlagTLSMinVersion   = clix.NewFlagNameAlias(prefix, name, "tls-min-version", "tlsmin")
		nameFlagTLSMaxVersion   = clix.NewFlagNameAlias(prefix, name, "tls-max-version", "tlsmax")
	)

	network := cfg.NetworkValue()

	return &Server{
		Name: name,

		PredeterminedNetwork: cfg.NetworkPredetermined(),

		FlagNetwork: &cli.StringFlag{
			Name:        nameFlagNetwork.Name,
			Aliases:     nameFlagNetwork.Aliases,
			Usage:       "network to connect" + cfg.NetworkUsage(),
			EnvVars:     nameFlagNetwork.EnvVars,
			FilePath:    nameFlagNetwork.FilePath,
			Required:    cfg.NetworkRequired(),
			Value:       network,
			Destination: &network,
		},

		FlagAddress: &cli.StringFlag{
			Name:        nameFlagAddress.Name,
			Aliases:     nameFlagAddress.Aliases,
			Usage:       "address to connect",
			EnvVars:     nameFlagAddress.EnvVars,
			FilePath:    nameFlagAddress.FilePath,
			Required:    len(cfg.Address) == 0,
			Value:       cfg.Address,
			Destination: new(string),
		},

		FlagTLSCerts: &cli.StringSliceFlag{
			Name:        nameFlagTLSCerts.Name,
			Aliases:     nameFlagTLSCerts.Aliases,
			Usage:       "client certificate pem `file`",
			EnvVars:     nameFlagTLSCerts.EnvVars,
			FilePath:    nameFlagTLSCerts.FilePath,
			TakesFile:   true,
			Destination: &cli.StringSlice{},
		},

		FlagTLSCertKeys: &cli.StringSliceFlag{
			Name:        nameFlagTLSCertKeys.Name,
			Aliases:     nameFlagTLSCertKeys.Aliases,
			Usage:       "private key `file` for client certificate",
			EnvVars:     nameFlagTLSCertKeys.EnvVars,
			FilePath:    nameFlagTLSCertKeys.FilePath,
			TakesFile:   true,
			Destination: &cli.StringSlice{},
		},

		FlagTLSGenCert: &cli.BoolFlag{
			Name:        nameFlagTLSGenCert.Name,
			Aliases:     nameFlagTLSGenCert.Aliases,
			Usage:       "whether to create and use self signed certificate",
			EnvVars:     nameFlagTLSGenCert.EnvVars,
			FilePath:    nameFlagTLSGenCert.FilePath,
			Destination: new(bool),
		},

		FlagTLSVerifyClient: &cli.BoolFlag{
			Name:        nameFlagTLSVerifyClient.Name,
			Aliases:     nameFlagTLSVerifyClient.Aliases,
			Usage:       "verify client certificate (mTLS)",
			EnvVars:     nameFlagTLSVerifyClient.EnvVars,
			FilePath:    nameFlagTLSVerifyClient.FilePath,
			Destination: new(bool),
		},

		FlagTLSClientCAs: &cli.StringSliceFlag{
			Name:        nameFlagTLSClientCAs.Name,
			Aliases:     nameFlagTLSClientCAs.Aliases,
			Usage:       "root CAs certificate `file`",
			EnvVars:     nameFlagTLSClientCAs.EnvVars,
			FilePath:    nameFlagTLSClientCAs.FilePath,
			TakesFile:   true,
			Destination: &cli.StringSlice{},
		},

		FlagTLSMinVersion: &cli.GenericFlag{
			Name:     nameFlagTLSMinVersion.Name,
			Aliases:  nameFlagTLSMinVersion.Aliases,
			Usage:    "TLS minimum version",
			EnvVars:  nameFlagTLSMinVersion.EnvVars,
			FilePath: nameFlagTLSMinVersion.FilePath,
			Value:    NewTLSVersion(cfg.TLSMinVersion),
		},

		FlagTLSMaxVersion: &cli.GenericFlag{
			Name:     nameFlagTLSMaxVersion.Name,
			Aliases:  nameFlagTLSMaxVersion.Aliases,
			Usage:    "TLS maximum version",
			EnvVars:  nameFlagTLSMaxVersion.EnvVars,
			FilePath: nameFlagTLSMaxVersion.FilePath,
			Value:    NewTLSVersion(cfg.TLSMaxVersion),
		},

		FlagSet: clix.NewFlagSet(),
	}
}

func (f *Server) Flags() []cli.Flag {
	return clix.Flags(
		clix.FlagIf(!f.PredeterminedNetwork, f.FlagNetwork),
		f.FlagAddress,
		f.FlagTLSCerts,
		f.FlagTLSCertKeys,
		f.FlagTLSGenCert,
		f.FlagTLSVerifyClient,
		f.FlagTLSClientCAs,
		f.FlagTLSMinVersion,
		f.FlagTLSMaxVersion,
	)
}

func (f *Server) Before(c *cli.Context) error {
	delint.Must(f.FlagSet.Init(c))
	if f.FlagSet.IsSet(f.FlagTLSCerts) && f.TLSGenCert() {
		return fmt.Errorf(
			"%q and %q must not set at the same time",
			f.FlagTLSCerts.Name,
			f.FlagTLSGenCert.Name,
		)
	}
	if f.FlagSet.IsSet(f.FlagTLSCertKeys) && f.TLSGenCert() {
		return fmt.Errorf(
			"%q and %q must not set at the same time",
			f.FlagTLSCertKeys.Name,
			f.FlagTLSGenCert.Name,
		)
	}
	return nil
}

func (f *Server) Network() string {
	return *f.FlagNetwork.Destination
}

func (f *Server) Address() string {
	return *f.FlagAddress.Destination
}

func (f *Server) TLSCerts() []string {
	return f.FlagTLSCerts.Destination.Value()
}

func (f *Server) TLSCertKeys() []string {
	return f.FlagTLSCertKeys.Destination.Value()
}

func (f *Server) TLSGenCert() bool {
	return *f.FlagTLSGenCert.Destination
}

func (f *Server) TLSVerifyClient() bool {
	return *f.FlagTLSVerifyClient.Destination
}

func (f *Server) LookupTLSVerifyClient() (val bool, ok bool) {
	return f.TLSVerifyClient(), f.FlagSet.IsSet(f.FlagTLSVerifyClient)
}

func (f *Server) TLSClientCAs() []string {
	return f.FlagTLSClientCAs.Destination.Value()
}

func (f *Server) TLSMinVersion() uint16 {
	return f.FlagTLSMinVersion.Value.(*TLSVersion).Value()
}

func (f *Server) TLSMaxVersion() uint16 {
	return f.FlagTLSMaxVersion.Value.(*TLSVersion).Value()
}

func (f *Server) UseTLS() bool {
	list := []cli.Flag{
		f.FlagTLSCerts,
		f.FlagTLSCertKeys,
		f.FlagTLSGenCert,
		f.FlagTLSVerifyClient,
		f.FlagTLSClientCAs,
		f.FlagTLSMinVersion,
		f.FlagTLSMaxVersion,
	}
	for _, flag := range list {
		if f.FlagSet.IsSet(flag) {
			return true
		}
	}
	return false
}

func (f *Server) TLSConfig() (*tls.Config, error) {
	var certs []tls.Certificate
	if f.TLSGenCert() {
		hostname, _ := os.Hostname()
		cert, err := cert4now.Generate(
			cert4now.CommonName(f.Name),
			cert4now.ECDSA(elliptic.P384()),
			cert4now.AddDate(100, 0, 0),
			cert4now.DNSNames(f.Name, hostname),
		)
		if err != nil {
			return nil, err
		}
		certs = []tls.Certificate{cert}
	} else {
		certFiles := f.TLSCerts()
		keyFiles := f.TLSCertKeys()
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

	var clientCAs *x509.CertPool
	if files := f.TLSClientCAs(); len(files) > 0 {
		clientCAs = x509.NewCertPool()
		for _, v := range files {
			p, err := ioutil.ReadFile(v)
			if err != nil {
				return nil, fmt.Errorf("failed to read CA %q, %w", v, err)
			}
			if ok := clientCAs.AppendCertsFromPEM(p); !ok {
				return nil, fmt.Errorf("%q contains no certificate for client CAs", v)
			}
		}
	}

	clientAuth := tls.NoClientCert
	if v, ok := f.LookupTLSVerifyClient(); (ok && v) || (!ok && len(f.TLSClientCAs()) > 0) {
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

func (f *Server) ServerOptions(opt ...grpc.ServerOption) ([]grpc.ServerOption, error) {
	var opts []grpc.ServerOption
	if f.UseTLS() {
		cfg, err := f.TLSConfig()
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.Creds(credentials.NewTLS(cfg)))
	}
	return append(opts, opt...), nil
}

func (f *Server) NewServer(opt ...grpc.ServerOption) (*grpc.Server, error) {
	opts, err := f.ServerOptions(opt...)
	if err != nil {
		return nil, err
	}
	return grpc.NewServer(opts...), nil
}

func (f *Server) Listen() (net.Listener, error) {
	return f.ListenNetwork(f.Network())
}

func (f *Server) ListenNetwork(network string) (net.Listener, error) {
	return net.Listen(network, f.Address())
}

func (f *Server) ListenAndServe(s *grpc.Server) error {
	return f.ListenNetworkAndServe(f.Network(), s)
}

func (f *Server) ListenNetworkAndServe(network string, s *grpc.Server) error {
	lis, err := f.ListenNetwork(network)
	if err != nil {
		return err
	}
	return s.Serve(lis)
}
