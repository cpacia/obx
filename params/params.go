package params

import (
	"encoding/binary"
	"github.com/cpacia/obxd/models/blocks"
)

type NetID uint32

func (n *NetID) Bytes() []byte {
	ret := make([]byte, 4)
	binary.BigEndian.PutUint32(ret, uint32(*n))
	return ret
}

const (
	Mainnet  NetID = 0x1c9db355
	Testnet1 NetID = 0xbb798738
	Regest   NetID = 0xde2d41e6
)

type NetworkParams struct {
	NetworkID     NetID
	GenesisBlock  blocks.Block
	DefaultPort   string
	SeedAddrs     []string
	ListenAddrs   []string
	AddressPrefix string
}

var MainnetParams = NetworkParams{
	NetworkID:    Mainnet,
	GenesisBlock: MainnetGenesisBlock,
	SeedAddrs: []string{
		"/ip4/167.172.126.176/tcp/4001/p2p/12D3KooWHnpVyu9XDeFoAVayqr9hvc9xPqSSHtCSFLEkKgcz5Wro",
	},
	ListenAddrs: []string{
		"/ip4/0.0.0.0/tcp/9001",
		"/ip6/::/tcp/9001",
		"/ip4/0.0.0.0/udp/9001/quic",
		"/ip6/::/udp/9001/quic",
	},
	AddressPrefix: "ob",
}

var Testnet1Params = NetworkParams{
	NetworkID: Testnet1,
	SeedAddrs: []string{
		"/ip4/167.172.126.176/tcp/4001/p2p/12D3KooWHnpVyu9XDeFoAVayqr9hvc9xPqSSHtCSFLEkKgcz5Wro",
	},
	ListenAddrs: []string{
		"/ip4/0.0.0.0/tcp/9002",
		"/ip6/::/tcp/9002",
		"/ip4/0.0.0.0/udp/9002/quic",
		"/ip6/::/udp/9002/quic",
	},
	AddressPrefix: "tn1",
}

var RegestParams = NetworkParams{
	NetworkID: Regest,
	ListenAddrs: []string{
		"/ip4/0.0.0.0/tcp/9003",
		"/ip6/::/tcp/9003",
		"/ip4/0.0.0.0/udp/9003/quic",
		"/ip6/::/udp/9003/quic",
	},
	AddressPrefix: "reg",
}
