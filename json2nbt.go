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
	// return nbtOut.Bytes(), nil
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
				if i, ok := m["value"].(float64); ok {
					err = binary.Write(w, byteOrder, int8(i))
					if err != nil {
						return JsonParseError{"Error writing byte payload", err}
					}
				} else {
					return JsonParseError{"Tag Byte value field not a number", err}
				}
			case 2:
				if i, ok := m["value"].(float64); ok {
					err = binary.Write(w, byteOrder, int16(i))
					if err != nil {
						return JsonParseError{"Error writing short payload", err}
					}
				} else {
					return JsonParseError{"Tag Byte value field not a number", err}
				}
			case 3:
				if i, ok := m["value"].(float64); ok {
					err = binary.Write(w, byteOrder, int32(i))
					if err != nil {
						return JsonParseError{"Error writing int32 payload", err}
					}
				} else {
					return JsonParseError{"Tag Byte value field not a number", err}
				}
			case 4:
				if i, ok := m["value"].(float64); ok {
					err = binary.Write(w, byteOrder, int64(i))
					if err != nil {
						return JsonParseError{"Error writing int64 payload", err}
					}
				} else {
					return JsonParseError{"Tag Byte value field not a number", err}
				}
			case 5:
				if f, ok := m["value"].(float64); ok {
					err = binary.Write(w, byteOrder, float32(f))
					if err != nil {
						return JsonParseError{"Error writing float32 payload", err}
					}
				} else {
					return JsonParseError{"Tag Byte value field not a number", err}
				}
			case 6:
				if f, ok := m["value"].(float64); ok {
					err = binary.Write(w, byteOrder, f)
					if err != nil {
						return JsonParseError{"Error writing float64 payload", err}
					}
				} else {
					return JsonParseError{"Tag Byte value field not a number", err}
				}
			case 7:
				if values, ok := m["value"].([]interface{}); ok {
					err = binary.Write(w, byteOrder, int32(len(values)))
					if err != nil {
						return JsonParseError{"Error writing byte array length", err}
					}
					for _, value := range values {
						if i, ok := value.(float64); ok {
							err = binary.Write(w, byteOrder, int8(i))
							if err != nil {
								return JsonParseError{"Error writing element of byte array", err}
							}
						} else {
							return JsonParseError{"Tag Byte value field not a number", err}
						}
					}
				} else {
					fmt.Printf("%v\n", m["value"])
					return JsonParseError{"Tag Byte Array value field not an array", err}
				}
			case 8:
				if s, ok := m["value"].(string); ok {
					err = binary.Write(w, byteOrder, int16(len([]byte(s))))
					if err != nil {
						return JsonParseError{"Error writing string length", err}
					}
					err = binary.Write(w, byteOrder, []byte(s))
					if err != nil {
						return JsonParseError{"Error writing string payload", err}
					}
				} else {
					return JsonParseError{"Tag String value field not a number", err}
				}

			case 11:
				if values, ok := m["value"].([]interface{}); ok {
					err = binary.Write(w, byteOrder, int32(len(values)))
					if err != nil {
						return JsonParseError{"Error writing int32 array length", err}
					}
					for _, value := range values {
						if i, ok := value.(float64); ok {
							err = binary.Write(w, byteOrder, int32(i))
							if err != nil {
								return JsonParseError{"Error writing element of int32 array", err}
							}
						} else {
							return JsonParseError{"Tag Byte value field not a number", err}
						}
					}
				} else {
					fmt.Printf("%v\n", m["value"])
					return JsonParseError{"Tag Byte Array value field not an array", err}
				}
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
