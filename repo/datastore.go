package repo

import (
	"github.com/ipfs/go-datastore"
)

const Libp2pDatastoreKey = "/obx/libp2pkey/"

type Datastore interface {
	datastore.Datastore
	datastore.Batching
	datastore.PersistentDatastore
	datastore.TxnDatastore
}
