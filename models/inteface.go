package models

type Serializable interface {
	Serialize() ([]byte, error)
	Deserialize(data []byte) error
	MarshalJSON() (string, error)
	UnmarshalJSON(data string) error
}
