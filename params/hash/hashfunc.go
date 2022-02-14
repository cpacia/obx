package hash

import (
	"golang.org/x/crypto/blake2b"
)

const HashSize = 32

func HashFunc(data []byte) []byte {
	h := blake2b.Sum256(data)
	return h[:]
}
