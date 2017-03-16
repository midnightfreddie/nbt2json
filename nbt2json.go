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
func Nbt2Json(r *bytes.Reader, byteOrder binary.ByteOrder) ([]byte, error) {
	var data NbtTag
	err := binary.Read(r, byteOrder, &data.TagType)
	if err != nil {
		return nil, NbtParseError{"Reading TagType", err}
	}
	if data.TagType != 0 {
		var err error
		var nameLen int64
		nameLen, err = readInt(r, 2, byteOrder)
		if err != nil {
			return nil, NbtParseError{"Reading Name length", err}
		}
		name := make([]byte, nameLen)
		err = binary.Read(r, byteOrder, &name)
		if err != nil {
			return nil, NbtParseError{"Reading Name - is little/big endian byte order set correctly?", err}
		}
		data.Name = string(name[:])
	}
	switch data.TagType {
	case 0:
		// end tag; do nothing further
	case 1:
		data.Value, err = readInt(r, 1, byteOrder)
		if err != nil {
			return nil, NbtParseError{"Reading int8", err}
		}
	case 2:
		data.Value, err = readInt(r, 2, byteOrder)
		if err != nil {
			return nil, NbtParseError{"Reading int16", err}
		}
	case 3:
		data.Value, err = readInt(r, 4, byteOrder)
		if err != nil {
			return nil, NbtParseError{"Reading int32", err}
		}
	case 4:
		data.Value, err = readInt(r, 8, byteOrder)
		if err != nil {
			return nil, NbtParseError{"Reading int64", err}
		}
	case 10:
		// compound is currently broken by design as its recursiveness is returning JSON instead of data, but might be able to use raw json data type? Else rework data types for recursiveness
		var compound []json.RawMessage
		var tagtype int64
		for tagtype, err = readInt(r, 1, byteOrder); tagtype != 0; tagtype, err = readInt(r, 1, byteOrder) {
			if err != nil {
				return nil, NbtParseError{"compound: reading next tag type", err}
			}
			_, err = r.Seek(-1, 1)
			if err != nil {
				return nil, NbtParseError{"seeking back one", err}
			}
			tag, err := Nbt2Json(r, byteOrder)
			if err != nil {
				return nil, NbtParseError{"compound: reading a child tag", err}
			}
			compound = append(compound, json.RawMessage(string(tag)))
		}
		// var boo *json.RawMessage
		// boo = &compound
		// fmt.Printf("%v\n", boo)
		// data.Value = string(compound[:])
		// data.Value = json.RawMessage(string(compound[:]))
		// data.Value = boo
		data.Value = compound
	default:
		return nil, NbtParseError{"TagType not recognized", nil}
	}
	outJson, err := json.Marshal(data)
	return outJson, nil
}
