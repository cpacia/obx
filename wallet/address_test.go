package wallet

import (
	"crypto/rand"
	"github.com/cpacia/obxd/params"
	"github.com/libp2p/go-libp2p-core/crypto"
	"testing"
)

func TestBasicAddress(t *testing.T) {
	_, pubkey, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	_, viewKey, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	ss := SpendNote{
		threshold: 1,
		pubkeys:   []crypto.PubKey{pubkey},
	}

	addr, err := NewBasicAddress(ss, viewKey, &params.MainnetParams)
	if err != nil {
		t.Fatal(err)
	}

	addr2, err := DecodeAddress(addr.String(), &params.MainnetParams)
	if err != nil {
		t.Fatal(err)
	}

	if addr2.String() != addr.String() {
		t.Error("Decoded address does not match encoded")
	}
}
