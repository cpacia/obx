package wallet

import (
	"fmt"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/cpacia/obxd/params"
	"github.com/cpacia/obxd/params/hash"
	"github.com/libp2p/go-libp2p-core/crypto"
)

type Address interface {
	EncodeAddress() string
	String() string
}

type SpendNote struct {
	threshold uint8
	pubkeys   []crypto.PubKey
}

func (s *SpendNote) Serialize() ([]byte, error) {
	ser := []byte{s.threshold}
	for _, pub := range s.pubkeys {
		keyBytes, err := pub.Raw()
		if err != nil {
			return nil, err
		}
		ser = append(ser, keyBytes...)
	}
	return ser, nil
}

type BasicAddress struct {
	hash    [32]byte
	viewKey crypto.PubKey
	version byte
	params  *params.NetworkParams
}

func NewBasicAddress(note SpendNote, viewKey crypto.PubKey, params *params.NetworkParams) (*BasicAddress, error) {
	ser, err := note.Serialize()
	if err != nil {
		return nil, err
	}
	h := hash.HashFunc(ser)
	var h2 [32]byte
	copy(h2[:], h)

	return &BasicAddress{
		hash:    h2,
		viewKey: viewKey,
		version: 1,
		params:  params,
	}, nil
}

func (a *BasicAddress) EncodeAddress() string {
	keyBytes, err := crypto.MarshalPublicKey(a.viewKey)
	if err != nil {
		return ""
	}
	converted, err := bech32.ConvertBits(append(a.hash[:], keyBytes...), 8, 5, true)
	if err != nil {
		return ""
	}
	combined := make([]byte, len(converted)+1)
	combined[0] = a.version
	copy(combined[1:], converted)
	ret, err := bech32.EncodeM(a.params.AddressPrefix, combined)
	if err != nil {
		return ""
	}
	return ret
}

func (a *BasicAddress) String() string {
	return a.EncodeAddress()
}

func DecodeAddress(addr string, params *params.NetworkParams) (Address, error) {
	// Decode the bech32 encoded address.
	_, data, err := bech32.DecodeNoLimit(addr)
	if err != nil {
		return nil, err
	}

	// The first byte of the decoded address is the version, it must exist.
	if len(data) < 1 {
		return nil, fmt.Errorf("no version")
	}

	// The remaining characters of the address returned are grouped into
	// words of 5 bits. In order to restore the original address bytes,
	// we'll need to regroup into 8 bit words.
	regrouped, err := bech32.ConvertBits(data[1:], 5, 8, false)
	if err != nil {
		return nil, err
	}

	var h2 [32]byte
	copy(h2[:], regrouped[:32])

	pub, err := crypto.UnmarshalPublicKey(regrouped[32:])
	if err != nil {
		return nil, err
	}

	ba := BasicAddress{
		params:  params,
		version: data[0],
		hash:    h2,
		viewKey: pub,
	}

	return &ba, nil
}
