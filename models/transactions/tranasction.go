package transactions

import (
	"bufio"
	"bytes"
	"github.com/cpacia/obxd/models"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func (tx *Transaction) ID() models.ID {
	if tx.GetStandardTransaction() != nil {
		ser, _ := tx.Serialize()
		return models.NewIDFromData(ser)
	}
	if tx.GetCoinbaseTransaction() != nil {
		ser, _ := tx.Serialize()
		return models.NewIDFromData(ser)
	}
	if tx.GetStakeTransaction() != nil {
		ser, _ := tx.Serialize()
		return models.NewIDFromData(ser)
	}
	return models.ID{}
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

func (tx *Transaction) MarshalJSON() ([]byte, error) {
	m := jsonpb.Marshaler{
		Indent: "    ",
	}
	var buf bytes.Buffer
	err := m.Marshal(bufio.NewWriter(&buf), tx)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (tx *Transaction) UnmarshalJSON(data []byte) error {
	newTx := &Transaction{}
	if err := jsonpb.Unmarshal(bytes.NewReader(data), newTx); err != nil {
		return err
	}
	tx = newTx
	return nil
}

func (tx *StandardTransaction) ID() models.ID {
	ser, _ := tx.Serialize()
	return models.NewIDFromData(ser)
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

func (tx *StandardTransaction) MarshalJSON() ([]byte, error) {
	m := jsonpb.Marshaler{
		Indent: "    ",
	}
	var buf bytes.Buffer
	err := m.Marshal(bufio.NewWriter(&buf), tx)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (tx *StandardTransaction) UnmarshalJSON(data []byte) error {
	newTx := &StandardTransaction{}
	if err := jsonpb.Unmarshal(bytes.NewReader(data), newTx); err != nil {
		return err
	}
	tx = newTx
	return nil
}

func (tx *CoinbaseTransaction) ID() models.ID {
	ser, _ := tx.Serialize()
	return models.NewIDFromData(ser)
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

func (tx *CoinbaseTransaction) MarshalJSON() ([]byte, error) {
	m := jsonpb.Marshaler{
		Indent: "    ",
	}
	var buf bytes.Buffer
	err := m.Marshal(bufio.NewWriter(&buf), tx)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (tx *CoinbaseTransaction) UnmarshalJSON(data []byte) error {
	newTx := &CoinbaseTransaction{}
	if err := jsonpb.Unmarshal(bytes.NewReader(data), newTx); err != nil {
		return err
	}
	tx = newTx
	return nil
}

func (tx *StakeTransaction) ID() models.ID {
	ser, _ := tx.Serialize()
	return models.NewIDFromData(ser)
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

func (tx *StakeTransaction) MarshalJSON() ([]byte, error) {
	m := jsonpb.Marshaler{
		Indent: "    ",
	}
	var buf bytes.Buffer
	err := m.Marshal(bufio.NewWriter(&buf), tx)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (tx *StakeTransaction) UnmarshalJSON(data []byte) error {
	newTx := &StakeTransaction{}
	if err := jsonpb.Unmarshal(bytes.NewReader(data), newTx); err != nil {
		return err
	}
	tx = newTx
	return nil
}
