package nbt2json

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/ghodss/yaml"
)

// Name is the json document's name:
const Name = "Named Binary Tag to JSON"

// Version is the json document's nbt2JsonVersion:
const Version = "0.3.3"

// nbt2JsonUrl is inserted in the first tag as nbt2JsonUrl
const nbt2JsonUrl = "https://github.com/midnightfreddie/nbt2json"

// NbtJson is the top-level JSON document
type NbtJson struct {
	Name           string             `json:"name"`
	Version        string             `json:"version"`
	Nbt2JsonUrl    string             `json:"nbt2JsonUrl"`
	ConversionTime string             `json:"conversionTime,omitempty"`
	Comment        string             `json:"comment,omitempty"`
	Nbt            []*json.RawMessage `json:"nbt"`
}

// NbtTag represents one NBT tag for each struct
type NbtTag struct {
	TagType byte        `json:"tagType"`
	Name    string      `json:"name"`
	Value   interface{} `json:"value,omitempty"`
}

// NbtTagList represents an NBT tag list
type NbtTagList struct {
	TagListType byte          `json:"tagListType"`
	List        []interface{} `json:"list"`
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

// Nbt2Yaml converts uncompressed NBT byte array to YAML byte array
func Nbt2Yaml(b []byte, byteOrder binary.ByteOrder, comment string) ([]byte, error) {
	jsonOut, err := Nbt2Json(b, byteOrder, comment)
	if err != nil {
		return nil, err
	}
	yamlOut, err := yaml.JSONToYAML(jsonOut)
	if err != nil {
		return yamlOut, NbtParseError{"Error converting JSON to YAML. Oops. JSON conversion succeeded, so please report this error and use JSON instead.", err}
	}
	return yamlOut, nil
}

// Nbt2Json converts uncompressed NBT byte array to JSON byte array
func Nbt2Json(b []byte, byteOrder binary.ByteOrder, comment string) ([]byte, error) {
	var nbtJson NbtJson
	nbtJson.Name = Name
	nbtJson.Version = Version
	nbtJson.Nbt2JsonUrl = nbt2JsonUrl
	nbtJson.ConversionTime = time.Now().Format(time.RFC3339)
	nbtJson.Comment = comment
	buf := bytes.NewReader(b)
	// var nbtJson.nbt []*json.RawMessage
	for buf.Len() > 0 {
		element, err := getTag(buf, byteOrder)
		if err != nil {
			return nil, err
		}
		myTemp := json.RawMessage(element)
		nbtJson.Nbt = append(nbtJson.Nbt, &myTemp)
	}
	jsonOut, err := json.MarshalIndent(nbtJson, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonOut, nil
}

// getTag broken out form Nbt2Json to allow recursion with reader but public input is []byte
func getTag(r *bytes.Reader, byteOrder binary.ByteOrder) ([]byte, error) {
	var data NbtTag
	err := binary.Read(r, byteOrder, &data.TagType)
	if err != nil {
		return nil, NbtParseError{"Reading TagType", err}
	}
	// do not try to fetch name for TagType 0 which is compound end tag
	if data.TagType != 0 {
		var err error
		var nameLen int16
		err = binary.Read(r, byteOrder, &nameLen)
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
	data.Value, err = getPayload(r, byteOrder, data.TagType)
	if err != nil {
		return nil, err
	}
	outJson, err := json.MarshalIndent(data, "", "  ")
	return outJson, err
}

// Gets the tag payload. Had to break this out from the main function to allow tag list recursion
func getPayload(r *bytes.Reader, byteOrder binary.ByteOrder, tagType byte) (interface{}, error) {
	var output interface{}
	var err error
	switch tagType {
	case 0:
		// end tag for compound; do nothing further
	case 1:
		var i int8
		err = binary.Read(r, byteOrder, &i)
		if err != nil {
			return nil, NbtParseError{"Reading int8", err}
		}
		output = i
	case 2:
		var i int16
		err = binary.Read(r, byteOrder, &i)
		if err != nil {
			return nil, NbtParseError{"Reading int16", err}
		}
		output = i
	case 3:
		var i int32
		err = binary.Read(r, byteOrder, &i)
		if err != nil {
			return nil, NbtParseError{"Reading int32", err}
		}
		output = i
	case 4:
		var i int64
		err = binary.Read(r, byteOrder, &i)
		if err != nil {
			return nil, NbtParseError{"Reading int64", err}
		}
		output = i
	case 5:
		var f float32
		err = binary.Read(r, byteOrder, &f)
		if err != nil {
			return nil, NbtParseError{"Reading float32", err}
		}
		output = f
	case 6:
		var f float64
		err = binary.Read(r, byteOrder, &f)
		if err != nil {
			return nil, NbtParseError{"Reading float64", err}
		}
		if math.IsNaN(f) {
			output = "NaN"
		} else {
			output = f
		}
	case 7:
		var byteArray []int8
		var oneByte int8
		var numRecords int32
		err := binary.Read(r, byteOrder, &numRecords)
		if err != nil {
			return nil, NbtParseError{"Reading byte array tag length", err}
		}
		for i := int32(1); i <= numRecords; i++ {
			err = binary.Read(r, byteOrder, &oneByte)
			if err != nil {
				return nil, NbtParseError{"Reading byte in byte array tag", err}
			}
			byteArray = append(byteArray, oneByte)
		}
		output = byteArray
	case 8:
		var strLen int16
		err := binary.Read(r, byteOrder, &strLen)
		if err != nil {
			return nil, NbtParseError{"Reading string tag length", err}
		}
		utf8String := make([]byte, strLen)
		err = binary.Read(r, byteOrder, &utf8String)
		if err != nil {
			return nil, NbtParseError{"Reading string tag data", err}
		}
		output = string(utf8String[:])
	case 9:
		var tagList NbtTagList
		err = binary.Read(r, byteOrder, &tagList.TagListType)
		if err != nil {
			return nil, NbtParseError{"Reading TagType", err}
		}
		var numRecords int32
		err := binary.Read(r, byteOrder, &numRecords)
		if err != nil {
			return nil, NbtParseError{"Reading list tag length", err}
		}
		for i := int32(1); i <= numRecords; i++ {
			payload, err := getPayload(r, byteOrder, tagList.TagListType)
			if err != nil {
				return nil, NbtParseError{"Reading list tag item", err}
			}
			tagList.List = append(tagList.List, payload)
		}
		output = tagList
	case 10:
		var compound []json.RawMessage
		var tagType byte
		for err = binary.Read(r, byteOrder, &tagType); tagType != 0; err = binary.Read(r, byteOrder, &tagType) {
			if err != nil {
				return nil, NbtParseError{"compound: reading next tag type", err}
			}
			_, err = r.Seek(-1, 1)
			if err != nil {
				return nil, NbtParseError{"seeking back one", err}
			}
			tag, err := getTag(r, byteOrder)
			if err != nil {
				return nil, NbtParseError{"compound: reading a child tag", err}
			}
			compound = append(compound, json.RawMessage(string(tag)))
		}
		if compound == nil {
			// Explicitly give empty array else value will be null instead of []
			output = []int{}
		} else {
			output = compound
		}
	case 11:
		var intArray []int32
		var numRecords, oneInt int32
		err := binary.Read(r, byteOrder, &numRecords)
		if err != nil {
			return nil, NbtParseError{"Reading int array tag length", err}
		}
		for i := int32(1); i <= numRecords; i++ {
			err := binary.Read(r, byteOrder, &oneInt)
			if err != nil {
				return nil, NbtParseError{"Reading int in int array tag", err}
			}
			intArray = append(intArray, oneInt)
		}
		output = intArray
	case 12:
		var longArray []int64
		var numRecords, oneInt int64
		err := binary.Read(r, byteOrder, &numRecords)
		if err != nil {
			return nil, NbtParseError{"Reading long array tag length", err}
		}
		for i := int64(1); i <= numRecords; i++ {
			err := binary.Read(r, byteOrder, &oneInt)
			if err != nil {
				return nil, NbtParseError{"Reading long in long array tag", err}
			}
			longArray = append(longArray, oneInt)
		}
		output = longArray
	default:
		return nil, NbtParseError{"TagType not recognized", nil}
	}
	return output, nil
}
