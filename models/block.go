package models

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func (h *BlockHeader) Serialize() ([]byte, error) {
	return proto.Marshal(h)
}

func (h *BlockHeader) Deserialize(data []byte) error {
	newHeader := &BlockHeader{}
	if err := proto.Unmarshal(data, newHeader); err != nil {
		return err
	}
	h = newHeader
	return nil
}

func (h *BlockHeader) MarshalJSON() (string, error) {
	m := jsonpb.Marshaler{
		Indent: "    ",
	}
	return m.MarshalToString(h)
}

func (h *BlockHeader) UnmarshalJSON(data string) error {
	newHeader := &BlockHeader{}
	if err := jsonpb.UnmarshalString(data, newHeader); err != nil {
		return err
	}
	h = newHeader
	return nil
}

func (b *Block) Serialize() ([]byte, error) {
	return proto.Marshal(b)
}

func (b *Block) Deserialize(data []byte) error {
	newBlock := &Block{}
	if err := proto.Unmarshal(data, newBlock); err != nil {
		return err
	}
	b = newBlock
	return nil
}

func (b *Block) MarshalJSON() (string, error) {
	m := jsonpb.Marshaler{
		Indent: "    ",
	}
	return m.MarshalToString(b)
}

func (b *Block) UnmarshalJSON(data string) error {
	newBlock := &Block{}
	if err := jsonpb.UnmarshalString(data, newBlock); err != nil {
		return err
	}
	b = newBlock
	return nil
}
