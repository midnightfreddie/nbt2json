package nbt2json

import (
	"encoding/json"
)

// Nbt2Json represents one NBT tag for each struct
type Nbt2Json struct {
	TagType byte `json:"tagType"`
	Name string `json:"name"`
	Value *json.RawMessage
}

func NewNbt2Json() *Nbt2Json {
	return &Nbt2Json{TagType: 0, Name: "Hi this isn't valid"}
}