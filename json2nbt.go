package nbt2json

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"io"
)

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

	// testing
	// fmt.Printf("%v\n", jsonData)
	// m := jsonData.(map[string]interface{})
	// fmt.Printf("%v\n", m["name"])
	// i := m["tagType"].(float64)
	// fmt.Println(i)

	err = writeTag(nbtOut, byteOrder, jsonData)
	if err != nil {
		return nil, err
	}
	// nbtOut.Write([]byte("yomamma"))
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Printf("%v\n", nbtOut.Bytes())
	// fmt.Printf("%v\n", hex.Dump(nbtOut.Bytes()))
	// fmt.Printf("%v\n", []byte(hex.Dump(nbtOut.Bytes())))
	return []byte(hex.Dump(nbtOut.Bytes())), nil
}

func writeTag(w io.Writer, byteOrder binary.ByteOrder, myMap interface{}) error {
	var err error
	_, err = w.Write([]byte("Hello"))
	return err
}
