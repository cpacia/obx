package params

import (
	"crypto"
	"encoding/binary"
)

var HashFunc = crypto.BLAKE2b_256.New().Sum

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
	NetworkID NetID
}

var MainnetParams = NetworkParams{
	NetworkID: Mainnet,
}

var Testnet1Params = NetworkParams{
	NetworkID: Testnet1,
}

var RegestParams = NetworkParams{
	NetworkID: Regest,
}
