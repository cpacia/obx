package main

import (
	"context"
	"github.com/ipfs/go-ds-badger"
	"github.com/libp2p/go-libp2p-core/crypto"
	"go.uber.org/zap"
	"obx/net"
	params "obx/params"
	"obx/repo"
)

var log = zap.S()

type Server struct {
	cancelFunc context.CancelFunc
	config     *repo.Config
	params     *params.NetworkParams
	ds         repo.Datastore
	network    *net.Network
}

func BuildServer(config *repo.Config) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())

	if err := setupLogging(config.LogDir, config.LogLevel, config.Testnet); err != nil {
		return nil, err
	}

	var netParams *params.NetworkParams
	if config.Testnet {
		netParams = &params.Testnet1Params
	} else if config.Regest {
		netParams = &params.RegestParams
	} else {
		netParams = &params.MainnetParams
	}

	ds, err := badger.NewDatastore(config.DataDir, nil)
	if err != nil {
		return nil, err
	}

	var privKey crypto.PrivKey
	has, err := repo.HasNetworkKey(ds)
	if err != nil {
		return nil, err
	}
	if has {
		privKey, err = repo.LoadNetworkKey(ds)
		if err != nil {
			return nil, err
		}
	} else {
		privKey, _, err = repo.GenerateNetworkKeypair()
		if err != nil {
			return nil, err
		}
	}

	networkOpts := []net.Option{
		net.Datastore(ds),
		net.BootstrapAddrs(config.BoostrapAddrs),
		net.ListenAddrs(config.ListenAddrs),
		net.UserAgent(config.UserAgent),
		net.PrivateKey(privKey),
		net.NetID(netParams.NetworkID),
	}
	if config.DisableNATPortMap {
		networkOpts = append(networkOpts, net.DisableNatPortMap())
	}

	network, err := net.NewNetwork(ctx, networkOpts...)
	if err != nil {
		return nil, err
	}

	s := &Server{
		cancelFunc: cancel,
		config:     config,
		params:     netParams,
		ds:         ds,
		network:    network,
	}

	return s, nil
}

func (s *Server) Close() error {
	s.cancelFunc()
	if err := s.network.Close(); err != nil {
		return err
	}
	if err := s.ds.Close(); err != nil {
		return err
	}
	return nil
}
