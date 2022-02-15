package params

import (
	"github.com/cpacia/obxd/models/blocks"
	"github.com/cpacia/obxd/models/transactions"
	"math"
	"time"
)

var MainnetGenesisBlock = blocks.Block{
	Header: &blocks.BlockHeader{
		Producer_ID:   []byte{0x00}, //TODO
		Height:        0,
		Timestamp:     time.Unix(0, 0).Unix(), //TODO
		Parent:        make([]byte, 32),
		ValidatorRoot: make([]byte, 32),
		Version:       1,
		TxRoot:        []byte{0x00}, //TODO
		UtxoRoot:      []byte{0x00}, //TODO
		NullifierRoot: []byte{0x00}, //TODO
		Signature:     []byte{0x00}, //TODO
	},
	Transactions: []*transactions.Transaction{
		{
			Data: &transactions.Transaction_CoinbaseTransaction{
				CoinbaseTransaction: &transactions.CoinbaseTransaction{
					Validator_ID: []byte{0x00}, //TODO
					NewCoins:     math.MaxUint64 / 10,
					Outputs: []*transactions.Output{
						{
							Commitment:      []byte{0x00}, //TODO
							EphemeralPubkey: []byte{0x00}, //TODO
							Ciphertext:      []byte{0x00}, //TODO
						},
					},
					Signature: []byte{0x00}, //TODO
					Proof:     []byte{0x00}, //TODO
				},
			},
		},
	},
}
