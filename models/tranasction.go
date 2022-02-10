package models

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func (tx *Transaction) ID() ID {
	if tx.GetStandardTransaction() != nil {
		ser, _ := tx.Serialize()
		return NewIDFromData(ser)
	}
	if tx.GetCoinbaseTransaction() != nil {
		ser, _ := tx.Serialize()
		return NewIDFromData(ser)
	}
	if tx.GetStakeTransaction() != nil {
		ser, _ := tx.Serialize()
		return NewIDFromData(ser)
	}
	return ID{}
}

func (tx *Transaction) Serialize() ([]byte, error) {
	return proto.Marshal(tx)
}

func (tx *Transaction) Deserialize(data []byte) error {
	newTx := &Transaction{}
	if err := proto.Unmarshal(data, newTx); err != nil {
		return err
	}
	tx = newTx
	return nil
}

func (tx *Transaction) MarshalJSON() (string, error) {
	m := jsonpb.Marshaler{
		Indent: "    ",
	}
	return m.MarshalToString(tx)
}

func (tx *Transaction) UnmarshalJSON(data string) error {
	newTx := &Transaction{}
	if err := jsonpb.UnmarshalString(data, newTx); err != nil {
		return err
	}
	tx = newTx
	return nil
}

func (tx *StandardTransaction) ID() ID {
	ser, _ := tx.Serialize()
	return NewIDFromData(ser)
}

func (tx *StandardTransaction) Serialize() ([]byte, error) {
	return proto.Marshal(tx)
}

func (tx *StandardTransaction) Deserialize(data []byte) error {
	newTx := &StandardTransaction{}
	if err := proto.Unmarshal(data, newTx); err != nil {
		return err
	}
	tx = newTx
	return nil
}

func (tx *StandardTransaction) MarshalJSON() (string, error) {
	m := jsonpb.Marshaler{
		Indent: "    ",
	}
	return m.MarshalToString(tx)
}

func (tx *StandardTransaction) UnmarshalJSON(data string) error {
	newTx := &StandardTransaction{}
	if err := jsonpb.UnmarshalString(data, newTx); err != nil {
		return err
	}
	tx = newTx
	return nil
}

func (tx *CoinbaseTransaction) ID() ID {
	ser, _ := tx.Serialize()
	return NewIDFromData(ser)
}

func (tx *CoinbaseTransaction) Serialize() ([]byte, error) {
	return proto.Marshal(tx)
}

func (tx *CoinbaseTransaction) Deserialize(data []byte) error {
	newTx := &CoinbaseTransaction{}
	if err := proto.Unmarshal(data, newTx); err != nil {
		return err
	}
	tx = newTx
	return nil
}

func (tx *CoinbaseTransaction) MarshalJSON() (string, error) {
	m := jsonpb.Marshaler{
		Indent: "    ",
	}
	return m.MarshalToString(tx)
}

func (tx *CoinbaseTransaction) UnmarshalJSON(data string) error {
	newTx := &CoinbaseTransaction{}
	if err := jsonpb.UnmarshalString(data, newTx); err != nil {
		return err
	}
	tx = newTx
	return nil
}

func (tx *StakeTransaction) ID() ID {
	ser, _ := tx.Serialize()
	return NewIDFromData(ser)
}

func (tx *StakeTransaction) Serialize() ([]byte, error) {
	return proto.Marshal(tx)
}

func (tx *StakeTransaction) Deserialize(data []byte) error {
	newTx := &StakeTransaction{}
	if err := proto.Unmarshal(data, newTx); err != nil {
		return err
	}
	tx = newTx
	return nil
}

func (tx *StakeTransaction) MarshalJSON() (string, error) {
	m := jsonpb.Marshaler{
		Indent: "    ",
	}
	return m.MarshalToString(tx)
}

func (tx *StakeTransaction) UnmarshalJSON(data string) error {
	newTx := &StakeTransaction{}
	if err := jsonpb.UnmarshalString(data, newTx); err != nil {
		return err
	}
	tx = newTx
	return nil
}
