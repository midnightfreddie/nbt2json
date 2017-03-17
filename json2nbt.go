package nbt2json

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
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

	return nbtOut.Bytes(), nil
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
			err = writePayload(w, byteOrder, m, tagType)
			if err != nil {
				return err
			}

		} else {
			return JsonParseError{"tagType is not numeric", err}
		}
	} else {
		return JsonParseError{"writeTag: myMap is not map[string]interface{}", err}
	}
	return err
}

func writePayload(w io.Writer, byteOrder binary.ByteOrder, m map[string]interface{}, tagType float64) error {
	var err error

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
	case 9:
		// important: tagListType needs to be in scope to be passed to writePayload
		// := were keeping it in a lower scope and zeroing it out.
		var tagListType float64
		if listMap, ok := m["value"].(map[string]interface{}); ok {
			if tagListType, ok = listMap["tagListType"].(float64); ok {
				err = binary.Write(w, byteOrder, byte(tagListType))
				if err != nil {
					return JsonParseError{"While writing tag list type", err}
				}
			}
			if values, ok := listMap["list"].([]interface{}); ok {
				err = binary.Write(w, byteOrder, int32(len(values)))
				if err != nil {
					return JsonParseError{"While writing tag list size", err}
				}
				for _, value := range values {
					fakeTag := make(map[string]interface{})
					fakeTag["value"] = value
					err = writePayload(w, byteOrder, fakeTag, tagListType)
					if err != nil {
						return JsonParseError{"While writing tag list of type " + strconv.Itoa(int(tagListType)), err}
					}
				}
			} else if listMap["list"] == nil {
				// NBT lists can be null / nil and therefore aren't represented as an array in JSON
				err = binary.Write(w, byteOrder, int32(0))
				if err != nil {
					return JsonParseError{"While writing tag list null size", err}
				}
				return nil
			} else {
				return JsonParseError{"Tag List's List value field not an array", err}
			}

		} else {
			return JsonParseError{"Tag List value field not an object", err}
		}
	case 10:
		if values, ok := m["value"].([]interface{}); ok {
			for _, value := range values {
				err = writeTag(w, byteOrder, value)
				if err != nil {
					return JsonParseError{"While writing Compound tags", err}
				}
			}
			// write the end tag which is just a single byte 0
			err = binary.Write(w, byteOrder, byte(0))
			if err != nil {
				return JsonParseError{"Writing End tag", err}
			}
		} else {
			return JsonParseError{"Tag Compound value field not an array", err}
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
					return JsonParseError{"Tag Int value field not a number", err}
				}
			}
		} else {
			return JsonParseError{"Tag Int Array value field not an array", err}
		}
	default:
		return JsonParseError{"tagType " + strconv.Itoa(int(tagType)) + " is not recognized", err}
	}
	return err
}
