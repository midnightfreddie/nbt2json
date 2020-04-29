package main

import (
	"C"
	"fmt"
	"unsafe"
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

// NOTE: Functions don't do anything yet; I'm just trying to figure out how
//   to pass C-native values from Go
// Oh cool, these comments are in the .h file when no blank lines separate them
//export Nbt2Json
func Nbt2Json(byteArray unsafe.Pointer, length C.int) *C.char {
	var goByteArray = C.GoBytes(byteArray, length)
	fmt.Print("The first byte in the byte array is ")
	fmt.Println(goByteArray[0])
	var tempString = "Hello from a Go string"
	return C.CString(tempString)
}

// NOTE: Functions don't do anything yet; I'm just trying to figure out how
//   to pass C-native values from Go
//export Json2Nbt
func Json2Nbt(cString *C.char) {
	var s string
	s = C.GoString(cString)
	fmt.Println(s)
	// return []byte(s)
}

func main() {}
