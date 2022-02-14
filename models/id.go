package models

import (
	"encoding/hex"
	"fmt"
	"github.com/cpacia/obxd/params/hash"
)

var ErrIDStrSize = fmt.Errorf("max ID string length is %v bytes", hash.HashSize*2)

type ID [hash.HashSize]byte

func (id *ID) String() string {
	return hex.EncodeToString(id[:])
}

func (id *ID) SetBytes(data []byte) {
	copy(id[:], data)
}

func (id *ID) MarshalJSON() ([]byte, error) {
	return []byte(hex.EncodeToString(id[:])), nil
}

func (id *ID) UnmarshalJSON(data []byte) error {
	i, err := NewIDFromString(string(data))
	if err != nil {
		return err
	}
	id = &i
	return nil
}

func NewID(digest []byte) ID {
	var sh ID
	sh.SetBytes(digest)
	return sh
}

func NewIDFromString(id string) (ID, error) {
	// Return error if hash string is too long.
	if len(id) > hash.HashSize*2 {
		return ID{}, ErrIDStrSize
	}
	ret, err := hex.DecodeString(id)
	if err != nil {
		return ID{}, err
	}
	var newID ID
	newID.SetBytes(ret)
	return newID, nil
}

func NewIDFromData(data []byte) ID {
	var id ID
	hash := hash.HashFunc(data)
	id.SetBytes(hash)
	return id
}
