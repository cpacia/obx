package blockchain

import (
	"encoding/binary"
	"github.com/cpacia/obxd/models"
	"github.com/cpacia/obxd/params/hash"
)

type node interface {
	Hash() []byte
	Height() uint32
	Prunable() bool
}

type branchNode struct {
	hash     []byte
	height   uint32
	prunable bool
}

func (b *branchNode) Hash() []byte {
	return b.hash
}

func (b *branchNode) Height() uint32 {
	return b.height
}

func (b *branchNode) Prunable() bool {
	return b.prunable
}

type leafNode struct {
	*branchNode
}

type MerkleMountainRange struct {
	peaks     [][]node
	nElements uint64
}

func NewMerkleMountainRange() *MerkleMountainRange {
	return &MerkleMountainRange{
		peaks:     make([][]node, 1),
		nElements: 0,
	}
}

func (mmr *MerkleMountainRange) Insert(data []byte, protect bool) {
	defer func() {
		mmr.nElements++
	}()

	d := make([]byte, len(data)+8)
	copy(d[:8], nElementsToBytes(mmr.nElements))
	copy(d[len(data):], data)
	h := hash.HashFunc(d)
	l := &leafNode{
		&branchNode{
			hash:     h,
			prunable: !protect,
			height:   0,
		},
	}

	i := 0
	for {
		if len(mmr.peaks[i]) == 0 {
			mmr.peaks[i] = append(mmr.peaks[i], l)
			return
		}
		last := mmr.peaks[i][len(mmr.peaks[i])-1]
		_, ok := last.(*leafNode)
		if ok {
			mmr.peaks[i] = append(mmr.peaks[i], l)
			catHash := hashMerkleBranches(last.Hash(), h)
			b := &branchNode{
				hash:     catHash,
				prunable: true,
				height:   1,
			}
			mmr.peaks[i] = append(mmr.peaks[i], b)

			for s := len(mmr.peaks) - 1; s > 0; s-- {
				if peakRoot(mmr.peaks[s-1]).Height() == peakRoot(mmr.peaks[s]).Height() {
					mmr.mergePeaks(s-1, s)
				}
			}
			return
		}
		if len(mmr.peaks) == i+1 {
			mmr.peaks = append(mmr.peaks, []node{})
		}
		i++
	}
}

func (mmr *MerkleMountainRange) Root() models.ID {
	combined := make([]byte, 0, hash.HashSize*len(mmr.peaks))
	for _, peak := range mmr.peaks {
		combined = append(combined, peakRoot(peak).Hash()...)
	}
	root := hash.HashFunc(combined)
	id, _ := models.NewID(root)
	return id
}

func (mmr *MerkleMountainRange) Prune() {

}

func (mmr *MerkleMountainRange) mergePeaks(a, b int) {
	height := mmr.peaks[a][len(mmr.peaks[a])-1].Height()
	rootHash := hashMerkleBranches(peakRoot(mmr.peaks[a]).Hash(), peakRoot(mmr.peaks[b]).Hash())

	mmr.peaks[a] = append(mmr.peaks[a], mmr.peaks[b]...)
	mmr.peaks[a] = append(mmr.peaks[a], &branchNode{hash: rootHash, prunable: true, height: height + 1})
	mmr.peaks = mmr.peaks[:len(mmr.peaks)-1]
}

func peakRoot(rng []node) node {
	return rng[len(rng)-1]
}

func nElementsToBytes(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return b
}
