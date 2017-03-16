package nbt2json

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
)

// NbtTag represents one NBT tag for each struct
type NbtTag struct {
	TagType byte   `json:"tagType"`
	Name    string `json:"name"`
	Value   interface{}
}

// Nbt2Json ...
func Nbt2Json(b []byte, byteOrder binary.ByteOrder) ([]byte, error) {
	var data NbtTag
	buf := bytes.NewReader(b[:])
	err := binary.Read(buf, binary.LittleEndian, &data.TagType)
	if err != nil {
		return nil, err
	}
	outJson, err := json.Marshal(data)
	return outJson, nil
}
