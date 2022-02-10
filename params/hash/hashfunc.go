package hash

import (
	"golang.org/x/crypto/blake2b"
)

func HashFunc(data []byte) []byte {
	h := blake2b.Sum256(data)
	return h[:]
}
