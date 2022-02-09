package net

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	coreconmgr "github.com/libp2p/go-libp2p-core/connmgr"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	pb "github.com/libp2p/go-libp2p-pubsub/pb"
	quic "github.com/libp2p/go-libp2p-quic-transport"
	"github.com/libp2p/go-tcp-transport"
	"github.com/multiformats/go-multiaddr"
)

type Network struct {
	host        host.Host
	connManager coreconmgr.ConnManager
	dht         *dht.IpfsDHT
	pubsub      *pubsub.PubSub
}

func NewNetwork(ctx context.Context, opts ...Option) (*Network, error) {
	var cfg config
	for _, opt := range opts {
		opt(&cfg)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	var idht *dht.IpfsDHT

	cmgr, err := connmgr.NewConnManager(
		100,                                  // Lowwater
		400,                                  // HighWater,
		connmgr.WithGracePeriod(time.Minute), // GracePeriod
	)
	if err != nil {
		return nil, err
	}

	hostOpts := libp2p.ChainOptions(
		// Use the keypair we generated
		libp2p.Identity(cfg.privateKey),
		// Multiple listen addresses
		libp2p.ListenAddrStrings(cfg.listenAddrs...),
		// Noise and TLS
		libp2p.DefaultSecurity,

		// QUIC and TCP
		libp2p.Transport(
			libp2p.ChainOptions(
				libp2p.Transport(tcp.NewTCPTransport),
				libp2p.Transport(quic.NewTransport),
			),
		),

		libp2p.DefaultMuxers,

		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(cmgr),

		// Let this host use the DHT to find other hosts
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht, err = dht.New(ctx, h)
			return idht, err
		}),

		// Let this host use relays and advertise itself on relays if
		// it finds it is behind NAT. Use libp2p.Relay(options...) to
		// enable active relays and more.
		libp2p.EnableAutoRelay(),
		// If you want to help other peers to figure out if they are behind
		// NATs, you can launch the server-side of AutoNAT too (AutoRelay
		// already runs the client)
		//
		// This service is highly rate-limited and should not cause any
		// performance issues.
		libp2p.EnableNATService(),

		libp2p.UserAgent(cfg.userAgent),

		libp2p.Ping(true),
	)

	if !cfg.disableNatPortMap {
		hostOpts = libp2p.ChainOptions(libp2p.NATPortMap(), hostOpts)
	}

	host, err := libp2p.New(
		ctx,
		hostOpts,
	)

	// The last step to get fully up and running would be to connect to
	// bootstrap peers (or any other peers). We leave this commented as
	// this is an example and the peer will die as soon as it finishes, so
	// it is unnecessary to put strain on the network.
	for _, addr := range cfg.bootstrapAddrs {
		ma, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			return nil, fmt.Errorf("%w: malformatted bootstrap peer", NetworkConfigErr)
		}

		pi, err := peer.AddrInfoFromP2pAddr(ma)
		if err != nil {
			return nil, err
		}
		// We ignore errors as some bootstrap peers may be down
		// and that is fine.
		host.Connect(ctx, *pi)
	}

	// create a new PubSub service using the GossipSub router
	ps, err := pubsub.NewGossipSub(
		ctx,
		host,
		pubsub.WithNoAuthor(),
		pubsub.WithMessageIdFn(func(pmsg *pb.Message) string {
			h := sha256.Sum256(pmsg.Data)
			return hex.EncodeToString(h[:])
		}),
	)
	if err != nil {
		return nil, err
	}

	// TODO: the pubsub object must have a validator set for the block and transaction
	// topics so that invalid blocks and transactions will not be relayed.

	net := &Network{
		host:        host,
		connManager: cmgr,
		dht:         idht,
		pubsub:      ps,
	}

	go func() {
		<-ctx.Done()
		net.Close()
	}()

	return net, nil
}

func (n *Network) Close() error {
	if err := n.host.Close(); err != nil {
		return err
	}
	if err := n.connManager.Close(); err != nil {
		return err
	}
	if err := n.dht.Close(); err != nil {
		return err
	}
	return nil
}

func (n *Network) Host() host.Host {
	return n.host
}

func (n *Network) ConnManager() coreconmgr.ConnManager {
	return n.connManager
}

func (n *Network) DHT() *dht.IpfsDHT {
	return n.dht
}

func (n *Network) Pubsub() *pubsub.PubSub {
	return n.pubsub
}

func (n *Network) SubscribeBlocks() {

}

func (n *Network) SubscribeTransactions() {

}

func (n *Network) BroadcastBlock() {

}

func (n *Network) BroadcastTransaction() {

}
