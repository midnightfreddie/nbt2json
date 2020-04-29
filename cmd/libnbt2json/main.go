package main

import (
	"C"
	"fmt"
	"unsafe"

	"github.com/midnightfreddie/nbt2json"
)

// HelloDll is here as a test while I work out parameter passing
// Any functions or vars exposed in the shared lib must be capitalized (Go rule)
// The export comment is needed to have the "C" package make the item available
//   in the shared library. Note there must be no space between the // and
//   'export'
//export HelloDll
func HelloDll() {
	fmt.Println("Hello from the libnbt2json dll!")
}

// The NBT data must be in a byte array. Pass a pointer to the array and the
//   length of the array
// Temporarily hard-codeed for Bedrock / little-endian only
//export Nbt2Json
func Nbt2Json(byteArray unsafe.Pointer, length C.int) *C.char {
	var goByteArray = C.GoBytes(byteArray, length)
	jsonData, err := nbt2json.Nbt2Json(goByteArray, nbt2json.Bedrock, "")
	if err != nil {
		panic(err)
	}
	return C.CString(string(jsonData))
}

// The NBT data must be in a byte array. Pass a pointer to the array and the
//   length of the array
// Temporarily hard-codeed for Bedrock / little-endian only
//export Nbt2Yaml
func Nbt2Yaml(byteArray unsafe.Pointer, length C.int) *C.char {
	var goByteArray = C.GoBytes(byteArray, length)
	jsonData, err := nbt2json.Nbt2Yaml(goByteArray, nbt2json.Bedrock, "")
	if err != nil {
		panic(err)
	}
	return C.CString(string(jsonData))
}

// NOTE: Functions don't do anything yet; I'm just trying to figure out how
//   to pass C-native values to/from Go
//export Json2Nbt
func Json2Nbt(cString *C.char) {
	var s string
	s = C.GoString(cString)
	nbtData, err := nbt2json.Json2Nbt([]byte(s), nbt2json.Bedrock)
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
	fmt.Println("The first few bytes of the NBT:")
	for i := 0; i < 5; i++ {
		fmt.Printf("%d ", nbtData[i])
	}
	fmt.Println("")
	// return []byte(s)
}

// NOTE: Functions don't do anything yet; I'm just trying to figure out how
//   to pass C-native values to/from Go
//export Yaml2Nbt
func Yaml2Nbt(cString *C.char) unsafe.Pointer {
	var s string
	s = C.GoString(cString)
	nbtData, err := nbt2json.Yaml2Nbt([]byte(s), nbt2json.Bedrock)
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
	fmt.Println(len(nbtData))
	fmt.Println("The first few bytes of the NBT:")
	for i := 0; i < 5; i++ {
		fmt.Printf("%d ", nbtData[i])
	}
	fmt.Println("")
	cByteArray := C.CBytes(nbtData)
	return cByteArray
}

//export SomeGoString
func SomeGoString() string {
	return "This\x00 is a Go string"
}

func main() {}
