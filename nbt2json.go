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
	e error
}

func (e NbtParseError) Error() string {
	var s string
	if e.e != nil {
		s = fmt.Sprintf(": %s", e.e.Error())
	}
	return fmt.Sprintf("Error parsing NBT: %s%s", e.s, s)
}

// Reads 0-8 bytes and returns an int64 value
func readInt(r *bytes.Reader, numBytes int, byteOrder binary.ByteOrder) (i int64, err error) {
	var myInt64 []byte
	temp := make([]byte, numBytes)
	err = binary.Read(r, byteOrder, &temp)
	if err != nil {
		return i, NbtParseError{fmt.Sprintf("Reading %v bytes for intxx", numBytes), err}
	}
	padding := make([]byte, 8-numBytes)
	if byteOrder == binary.BigEndian {
		myInt64 = append(padding, temp...)
	} else if byteOrder == binary.LittleEndian {
		myInt64 = append(temp, padding...)
	} else {
		_ = myInt64
		return i, NbtParseError{"byteOrder not recognized", nil}
	}
	buf := bytes.NewReader(myInt64)
	err = binary.Read(buf, byteOrder, &i)
	return i, nil
}

// Nbt2Json ...
func Nbt2Json(b []byte, byteOrder binary.ByteOrder) ([]byte, error) {
	var data NbtTag
	buf := bytes.NewReader(b[:])
	err := binary.Read(buf, byteOrder, &data.TagType)
	if err != nil {
		return nil, NbtParseError{"Reading TagType", err}
	}
	if data.TagType != 0 {
		var err error
		var nameLen int64
		nameLen, err = readInt(buf, 2, byteOrder)
		if err != nil {
			return nil, NbtParseError{"Reading Name length", err}
		}
		name := make([]byte, nameLen)
		err = binary.Read(buf, byteOrder, &name)
		if err != nil {
			return nil, NbtParseError{"Reading Name - is little/big endian byte order set correctly?", err}
		}
		data.Name = string(name[:])
	}
	switch data.TagType {
	case 0:
		// end tag; do nothing further
	case 1:
		data.Value, err = readInt(buf, 1, byteOrder)
		if err != nil {
			return nil, NbtParseError{"Reading int8", err}
		}
	case 2:
		data.Value, err = readInt(buf, 2, byteOrder)
		if err != nil {
			return nil, NbtParseError{"Reading int16", err}
		}
	case 3:
		data.Value, err = readInt(buf, 4, byteOrder)
		if err != nil {
			return nil, NbtParseError{"Reading int32", err}
		}
	case 4:
		data.Value, err = readInt(buf, 8, byteOrder)
		if err != nil {
			return nil, NbtParseError{"Reading int64", err}
		}
	default:
		return nil, NbtParseError{"TagType not recognized", nil}
	}
	outJson, err := json.Marshal(data)
	return outJson, nil
}
