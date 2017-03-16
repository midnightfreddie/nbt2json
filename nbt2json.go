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
	err := binary.Read(buf, byteOrder, &data.TagType)
	if err != nil {
		return nil, err
	}
	switch data.TagType {
	case 0:
		//do nothing
	default:
		var err error
		var nameLen int16
		err = binary.Read(buf, byteOrder, &nameLen)
		if err != nil {
			return nil, err
		}
		name := make([]byte, nameLen)
		err = binary.Read(buf, byteOrder, &name)
		if err != nil {
			return nil, err
		}
		data.Name = string(name[:])
	}
	outJson, err := json.Marshal(data)
	return outJson, nil
}
