package models

type Serializable interface {
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(data []byte) error
}
