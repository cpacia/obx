package main

import (
	"context"
	"github.com/cpacia/obxd/net"
	params "github.com/cpacia/obxd/params"
	"github.com/cpacia/obxd/repo"
	"github.com/ipfs/go-ds-badger"
	"github.com/libp2p/go-libp2p-core/crypto"
	"go.uber.org/zap"
	"sort"
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
		if err := repo.PutNetworkKey(ds, privKey); err != nil {
			return nil, err
		}
	}

	var seedAddrs []string
	if config.SeedAddrs != nil {
		seedAddrs = config.SeedAddrs
	} else {
		seedAddrs = netParams.SeedAddrs
	}

	var listenAddrs []string
	if config.ListenAddrs != nil {
		listenAddrs = config.ListenAddrs
	} else {
		listenAddrs = netParams.ListenAddrs
	}

	networkOpts := []net.Option{
		net.Datastore(ds),
		net.SeedAddrs(seedAddrs),
		net.ListenAddrs(listenAddrs),
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

	s.printListenAddrs()

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

func (s *Server) printListenAddrs() {
	log.Infof("PeerID: %s", s.network.Host().ID().Pretty())
	var lisAddrs []string
	ifaceAddrs := s.network.Host().Addrs()
	for _, addr := range ifaceAddrs {
		lisAddrs = append(lisAddrs, addr.String())
	}
	sort.Strings(lisAddrs)
	for _, addr := range lisAddrs {
		log.Infof("Listening on %s", addr)
	}
}
