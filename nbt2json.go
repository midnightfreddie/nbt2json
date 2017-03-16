package nbt2json

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
)

// NbtTag represents one NBT tag for each struct
type NbtTag struct {
	TagType byte        `json:"tagType"`
	Name    string      `json:"name"`
	Value   interface{} `json:"value,omitempty"`
}

// NbtParseError is when the data does not match an expected pattern. Pass it message string and downstream error
type NbtParseError struct {
	s string
}

func (e NbtParseError) Error() string {
	return fmt.Sprintf("Error parsing NBT: %s", e.s)
}

// Nbt2Json ...
func Nbt2Json(b []byte, byteOrder binary.ByteOrder) ([]byte, error) {
	var data NbtTag
	buf := bytes.NewReader(b[:])
	err := binary.Read(buf, byteOrder, &data.TagType)
	if err != nil {
		return nil, NbtParseError{fmt.Sprintf("Reading TagType: %s", err.Error())}
	}
	if data.TagType != 0 {
		var err error
		var nameLen int16
		err = binary.Read(buf, byteOrder, &nameLen)
		if err != nil {
			return nil, NbtParseError{fmt.Sprintf("Reading Name length: %s", err.Error())}
		}
		name := make([]byte, nameLen)
		err = binary.Read(buf, byteOrder, &name)
		if err != nil {
			return nil, NbtParseError{fmt.Sprintf("Reading Name: %s", err.Error())}
		}
		data.Name = string(name[:])
	}
	switch data.TagType {
	case 0:
		// end tag; do nothing further
	case 2:
		data.Value = "test"
	default:
		return nil, NbtParseError{"TagType not recognized"}
	}
	outJson, err := json.Marshal(data)
	return outJson, nil
}
