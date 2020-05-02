package nbt2json

import (
	"encoding/binary"
	"fmt"
)

// Version is the json document's nbt2JsonVersion:
const Version = "0.3.4"

// Nbt2JsonUrl is inserted in the json document as nbt2JsonUrl
const Nbt2JsonUrl = "https://github.com/midnightfreddie/nbt2json"

// Name is the json document's name:
var Name = "Named Binary Tag to JSON"

// Bedrock is for Bedrock Edition (little endian NBT encoding); unable to make const, but **do not alter**
var Bedrock = binary.LittleEndian

// Java is for Java Edition (big endian NBT encoding); unable to make const, but **do not alter**
var Java = binary.BigEndian

// NbtParseError is when the nbt data does not match an expected pattern. Pass it message string and downstream error
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

// NOTE: Although this file is named "common.go", the above values are only
// used in nbt2json.go and the below in json2nbt.go. I mainly wanted to
// separate Version to a sensible place and threw in the other potentially-
// interesting values and exported structs. The other non-function exports are
// for marshalling via reflect and not because they need to be used by client code.

// JsonParseError is when the json data does not match an expected pattern. Pass it message string and downstream error
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
