package blockchain

import (
	"crypto/rand"
	"fmt"
	"testing"
	"time"
)

func TestBranchNode_Hash(t *testing.T) {
	data := make([][]byte, 1000000)
	for i := range data {
		data[i] = make([]byte, 5000)
		rand.Read(data[i])
	}
	start := time.Now()
	mmr := NewMerkleMountainRange()
	for _, d := range data {
		mmr.Insert(d, false)
	}
	root := mmr.Root()
	fmt.Println(root.String())
	fmt.Println(time.Since(start))
}

func BenchmarkBuildMerkleTreeStore(b *testing.B) {
	data := make([][]byte, 10000)
	for i := range data {
		data[i] = make([]byte, 32)
		rand.Read(data[i])
	}
	mmr := NewMerkleMountainRange()
	for _, d := range data {
		mmr.Insert(d, false)
	}

	b.Run("merkle", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			BuildMerkleTreeStore(append(data, make([]byte, 32)))
		}
	})

	b.Run("mmr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			mmr.Insert(make([]byte, 32), false)
			mmr.Root()
		}
	})
}
