package grpcflag

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/takumakei/go-delint"
	"github.com/takumakei/go-urfave-cli/clix"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Dialer struct {
	Name string

	// PredeterminedNetwork is true if the value of FlagNetwork is predetermined by the option.
	PredeterminedNetwork bool

	FlagNetwork       *cli.StringFlag
	FlagAddress       *cli.StringFlag
	FlagInsecure      *cli.BoolFlag
	FlagBlock         *cli.BoolFlag
	FlagTLSRootCAs    *cli.StringSliceFlag
	FlagTLSServerName *cli.StringFlag
	FlagTLSSkipVerify *cli.BoolFlag
	FlagTLSCerts      *cli.StringSliceFlag
	FlagTLSCertKeys   *cli.StringSliceFlag
	FlagTLSMinVersion *cli.GenericFlag
	FlagTLSMaxVersion *cli.GenericFlag

	FlagSet clix.FlagSet
}

func NewDialer(prefix clix.FlagPrefix, opt ...Option) *Dialer {
	return NewDialerName(prefix, "", opt...)
}

func NewDialerName(prefix clix.FlagPrefix, name string, opt ...Option) *Dialer {
	cfg := NewConfig(opt...)

	var (
		nameFlagNetwork       = clix.NewFlagNameAlias(prefix, name, "network", "net")
		nameFlagAddress       = clix.NewFlagNameAlias(prefix, name, "address", "addr")
		nameFlagInsecure      = clix.NewFlagNameAlias(prefix, name, "with-insecure", "insecure")
		nameFlagBlock         = clix.NewFlagNameAlias(prefix, name, "with-block", "block")
		nameFlagTLSRootCAs    = clix.NewFlagNameAlias(prefix, name, "tls-root-ca", "tlsca")
		nameFlagTLSServerName = clix.NewFlagNameAlias(prefix, name, "tls-server-name", "tlsserver")
		nameFlagTLSSkipVerify = clix.NewFlagNameAlias(prefix, name, "tls-skip-verify", "tlsinsecure")
		nameFlagTLSCerts      = clix.NewFlagNameAlias(prefix, name, "tls-cert", "tlscrt")
		nameFlagTLSCertKeys   = clix.NewFlagNameAlias(prefix, name, "tls-cert-key", "tlskey")
		nameFlagTLSMinVersion = clix.NewFlagNameAlias(prefix, name, "tls-min-version", "tlsmin")
		nameFlagTLSMaxVersion = clix.NewFlagNameAlias(prefix, name, "tls-max-version", "tlsmax")
	)

	network := cfg.NetworkValue()

	return &Dialer{
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

		FlagInsecure: &cli.BoolFlag{
			Name:        nameFlagInsecure.Name,
			Aliases:     nameFlagInsecure.Aliases,
			Usage:       "whether to use grpc.WithInsecure",
			EnvVars:     nameFlagInsecure.EnvVars,
			FilePath:    nameFlagInsecure.FilePath,
			Destination: new(bool),
		},

		FlagBlock: &cli.BoolFlag{
			Name:        nameFlagBlock.Name,
			Aliases:     nameFlagBlock.Aliases,
			Usage:       "whether to use grpc.WithBlock",
			EnvVars:     nameFlagBlock.EnvVars,
			FilePath:    nameFlagBlock.FilePath,
			Destination: new(bool),
		},

		FlagTLSRootCAs: &cli.StringSliceFlag{
			Name:        nameFlagTLSRootCAs.Name,
			Aliases:     nameFlagTLSRootCAs.Aliases,
			Usage:       "root CAs certificate `file`",
			EnvVars:     nameFlagTLSRootCAs.EnvVars,
			FilePath:    nameFlagTLSRootCAs.FilePath,
			TakesFile:   true,
			Destination: &cli.StringSlice{},
		},

		FlagTLSServerName: &cli.StringFlag{
			Name:        nameFlagTLSServerName.Name,
			Aliases:     nameFlagTLSServerName.Aliases,
			Usage:       "ServerName of tls.Config",
			EnvVars:     nameFlagTLSServerName.EnvVars,
			FilePath:    nameFlagTLSServerName.FilePath,
			Destination: new(string),
		},

		FlagTLSSkipVerify: &cli.BoolFlag{
			Name:        nameFlagTLSSkipVerify.Name,
			Aliases:     nameFlagTLSSkipVerify.Aliases,
			Usage:       "InsecureSkipVerify of tls.Config",
			EnvVars:     nameFlagTLSSkipVerify.EnvVars,
			FilePath:    nameFlagTLSSkipVerify.FilePath,
			Destination: new(bool),
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

func (f *Dialer) Flags() []cli.Flag {
	return clix.Flags(
		clix.FlagIf(!f.PredeterminedNetwork, f.FlagNetwork),
		f.FlagAddress,
		f.FlagInsecure,
		f.FlagBlock,
		f.FlagTLSRootCAs,
		f.FlagTLSServerName,
		f.FlagTLSSkipVerify,
		f.FlagTLSCerts,
		f.FlagTLSCertKeys,
		f.FlagTLSMinVersion,
		f.FlagTLSMaxVersion,
	)
}

func (f *Dialer) Before(c *cli.Context) error {
	delint.Must(f.FlagSet.Init(c))
	return nil
}

func (f *Dialer) Network() string {
	return *f.FlagNetwork.Destination
}

func (f *Dialer) Address() string {
	return *f.FlagAddress.Destination
}

func (f *Dialer) Insecure() bool {
	return *f.FlagInsecure.Destination
}

func (f *Dialer) Block() bool {
	return *f.FlagBlock.Destination
}

func (f *Dialer) TLSRootCAs() []string {
	return f.FlagTLSRootCAs.Destination.Value()
}

func (f *Dialer) TLSServerName() string {
	return *f.FlagTLSServerName.Destination
}

func (f *Dialer) TLSSkipVerify() bool {
	return *f.FlagTLSSkipVerify.Destination
}

func (f *Dialer) TLSCerts() []string {
	return f.FlagTLSCerts.Destination.Value()
}

func (f *Dialer) TLSCertKeys() []string {
	return f.FlagTLSCertKeys.Destination.Value()
}

func (f *Dialer) TLSMinVersion() uint16 {
	return f.FlagTLSMinVersion.Value.(*TLSVersion).Value()
}

func (f *Dialer) TLSMaxVersion() uint16 {
	return f.FlagTLSMaxVersion.Value.(*TLSVersion).Value()
}

func (f *Dialer) TLSConfig() (*tls.Config, error) {
	var rootCAs *x509.CertPool
	if files := f.TLSRootCAs(); len(files) > 0 {
		rootCAs = x509.NewCertPool()
		for _, v := range files {
			p, err := ioutil.ReadFile(v)
			if err != nil {
				return nil, err
			}
			if ok := rootCAs.AppendCertsFromPEM(p); !ok {
				return nil, fmt.Errorf("%s contains no certificate for RootCAs", v)
			}
		}
	}

	var certs []tls.Certificate
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

	cfg := &tls.Config{
		Certificates:       certs,
		RootCAs:            rootCAs,
		ServerName:         f.TLSServerName(),
		InsecureSkipVerify: f.TLSSkipVerify(),
		MinVersion:         f.TLSMinVersion(),
		MaxVersion:         f.TLSMaxVersion(),
	}
	return cfg, nil
}

func (f *Dialer) DialOption() ([]grpc.DialOption, error) {
	var opts []grpc.DialOption

	if f.Insecure() {
		opts = append(opts, grpc.WithInsecure())
	} else {
		cfg, err := f.TLSConfig()
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(cfg)))
	}

	if f.Block() {
		opts = append(opts, grpc.WithBlock())
	}

	opts = append(opts, grpc.WithContextDialer(func(ctx context.Context, address string) (net.Conn, error) {
		d := &net.Dialer{}
		return d.DialContext(ctx, f.Network(), address)
	}))

	return opts, nil
}

func (f *Dialer) Dial() (*grpc.ClientConn, error) {
	return f.DialContext(context.Background())
}

func (f *Dialer) DialContext(ctx context.Context) (*grpc.ClientConn, error) {
	opts, err := f.DialOption()
	if err != nil {
		return nil, err
	}
	return grpc.DialContext(ctx, f.Address(), opts...)
}
