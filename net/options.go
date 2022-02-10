package net

import (
	"errors"
	"fmt"
	"github.com/cpacia/obxd/params"
	"github.com/cpacia/obxd/repo"
	"github.com/libp2p/go-libp2p-core/crypto"
)

var ErrNetworkConfig = errors.New("network config error")

// Option is configuration option function for the Network
type Option func(cfg *config) error

func NetID(netID params.NetID) Option {
	return func(cfg *config) error {
		cfg.netID = netID
		return nil
	}
}

func PrivateKey(privKey crypto.PrivKey) Option {
	return func(cfg *config) error {
		cfg.privateKey = privKey
		return nil
	}
}

func Datastore(ds repo.Datastore) Option {
	return func(cfg *config) error {
		cfg.datastore = ds
		return nil
	}
}

func UserAgent(s string) Option {
	return func(cfg *config) error {
		cfg.userAgent = s
		return nil
	}
}

func ListenAddrs(addrs []string) Option {
	return func(cfg *config) error {
		cfg.listenAddrs = addrs
		return nil
	}
}

func SeedAddrs(addrs []string) Option {
	return func(cfg *config) error {
		cfg.seedAddrs = addrs
		return nil
	}
}

func DisableNatPortMap() Option {
	return func(cfg *config) error {
		cfg.disableNatPortMap = true
		return nil
	}
}

type config struct {
	netID             params.NetID
	userAgent         string
	seedAddrs         []string
	listenAddrs       []string
	disableNatPortMap bool
	privateKey        crypto.PrivKey
	datastore         repo.Datastore
}

func (cfg *config) validate() error {
	if cfg.privateKey == nil {
		return fmt.Errorf("%w: private key is nil", ErrNetworkConfig)
	}
	if cfg.listenAddrs == nil {
		return fmt.Errorf("%w: listen addrs is nil", ErrNetworkConfig)
	}
	if cfg.datastore == nil {
		return fmt.Errorf("%w: datastore is nil", ErrNetworkConfig)
	}
	return nil
}
