package models

import (
	"encoding/hex"
	"fmt"
	"obx/params"
)

const (
	HashSize          = 32
	MaxHashStringSize = 64
)

var ErrIDStrSize = fmt.Errorf("max ID string length is %v bytes", MaxHashStringSize)

type ID [HashSize]byte

func (id *ID) String() string {
	return hex.EncodeToString(id[:])
}

func (id *ID) SetBytes(data []byte) error {
	nhlen := len(data)
	if nhlen != HashSize {
		return fmt.Errorf("invalid hash length of %v, want %v", nhlen,
			HashSize)
	}
	copy(id[:], data)

	return nil
}

func NewID(data []byte) (ID, error) {
	var sh ID
	err := sh.SetBytes(data)
	if err != nil {
		return ID{}, err
	}
	return sh, err
}

func NewIDFromString(id string) (ID, error) {
	// Return error if hash string is too long.
	if len(id) > MaxHashStringSize {
		return ID{}, ErrIDStrSize
	}
	ret, err := hex.DecodeString(id)
	if err != nil {
		return ID{}, err
	}
	var newID ID
	if err := newID.SetBytes(ret); err != nil {
		return ID{}, err
	}
	return newID, nil
}

func NewIDFromData(data []byte) ID {
	var id ID
	hash := params.HashFunc(data)
	id.SetBytes(hash)
	return id
}
