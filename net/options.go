package net

import (
	"errors"
	"fmt"
	"github.com/libp2p/go-libp2p-core/crypto"
)

var NetworkConfigErr = errors.New("network config error")

// Option is configuration option function for the Network
type Option func(cfg *config) error

func PrivateKey(privKey crypto.PrivKey) Option {
	return func(cfg *config) error {
		cfg.privateKey = privKey
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

func BootstrapAddrs(addrs []string) Option {
	return func(cfg *config) error {
		cfg.bootstrapAddrs = addrs
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
	userAgent         string
	bootstrapAddrs    []string
	listenAddrs       []string
	disableNatPortMap bool
	privateKey        crypto.PrivKey
}

func (cfg *config) validate() error {
	if cfg.privateKey == nil {
		return fmt.Errorf("%w: private key is nil", NetworkConfigErr)
	}
	if cfg.bootstrapAddrs == nil {
		return fmt.Errorf("%w: bootstrap addrs is nil", NetworkConfigErr)
	}
	if cfg.listenAddrs == nil {
		return fmt.Errorf("%w: listen addrs is nil", NetworkConfigErr)
	}
	return nil
}
