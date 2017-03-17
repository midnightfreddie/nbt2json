package nbt2json

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
)

// JsonParseError is when the data does not match an expected pattern. Pass it message string and downstream error
type JsonParseError struct {
	s string
	e error
}

func (e JsonParseError) Error() string {
	var s string
	if e.e != nil {
		s = fmt.Sprintf(": %s", e.e.Error())
	}
	return fmt.Sprintf("Error parsing json2nbt: %s%s", e.s, s)
}

// Json2Nbt converts JSON byte array to uncompressed NBT byte array
// During development, returning a hex dump instead of raw data
func Json2Nbt(b []byte, byteOrder binary.ByteOrder) ([]byte, error) {
	nbtOut := new(bytes.Buffer)
	var jsonData interface{}
	var err error
	err = json.Unmarshal(b, &jsonData)
	if err != nil {
		return nil, err
	}
	err = writeTag(nbtOut, byteOrder, jsonData)
	if err != nil {
		return nil, err
	}

	return []byte(hex.Dump(nbtOut.Bytes())), nil
}

func writeTag(w io.Writer, byteOrder binary.ByteOrder, myMap interface{}) error {
	var err error
	if m, ok := myMap.(map[string]interface{}); ok {
		if tagType, ok := m["tagType"].(float64); ok {
			if tagType == 0 {
				// not expecting a 0 tag, but if it occurs just ignore it
				return nil
			}
			err = binary.Write(w, byteOrder, byte(tagType))
			if err != nil {
				return JsonParseError{"Error writing tagType" + string(byte(tagType)), err}
			}
			if name, ok := m["name"].(string); ok {
				err = binary.Write(w, byteOrder, int16(len(name)))
				if err != nil {
					return JsonParseError{"Error writing name length", err}
				}
				err = binary.Write(w, byteOrder, []byte(name))
				if err != nil {
					return JsonParseError{"Error converting name", err}
				}
			} else {
				return JsonParseError{"name field not a string", err}
			}
			switch int(tagType) {
			case 1:
			case 2:
			default:
				return JsonParseError{"tagType " + string(int(tagType)) + " is not recognized", err}
			}
		} else {
			return JsonParseError{"tagType is not numeric", err}
		}
	} else {
		return JsonParseError{"writeTag: myMap is not map[string]interface{}", err}
	}
	return err
}
