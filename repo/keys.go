package repo

import (
	"context"
	"crypto/rand"
	"github.com/ipfs/go-datastore"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
)

func HasNetworkKey(ds Datastore) (bool, error) {
	return ds.Has(context.Background(), datastore.NewKey(Libp2pDatastoreKey))
}

func LoadNetworkKey(ds Datastore) (crypto.PrivKey, error) {
	keyBytes, err := ds.Get(context.Background(), datastore.NewKey(Libp2pDatastoreKey))
	if err != nil {
		return nil, err
	}
	return crypto.UnmarshalPrivateKey(keyBytes)
}

func GenerateNetworkKeypair() (crypto.PrivKey, crypto.PubKey, error) {
	privkey, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return privkey, nil, nil
}
