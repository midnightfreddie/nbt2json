package main

import (
	"C"
	"fmt"
)

//export HelloDll
func HelloDll() {
	fmt.Println("Hello from the libnbt2json dll!")
}

func main() {}
