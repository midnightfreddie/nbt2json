package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"encoding/binary"

	"compress/gzip"

	"github.com/midnightfreddie/nbt2json"
	"github.com/urfave/cli"
)

func main() {
	var nbtFile string
	var byteOrder binary.ByteOrder
	var skipBytes int
	app := cli.NewApp()
	app.Name = "NBT to JSON"
	app.Version = "0.0.0"
	app.Usage = "Converts NBT-encoded data to JSON"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "reverse, json2nbt, r",
			// Usage: "Convert JSON to NBT instead",
			Usage:  "FLAG UNDER DEVELOPMENT, NOT YET WORKING",
			Hidden: true,
		},
		cli.BoolTFlag{
			Name:  "little-endian, little, mcpe, l",
			Usage: "Number format for Minecraft Pocket Edition and Windows 10 Edition (default)",
		},
		cli.BoolFlag{
			Name:  "big-endian, big, java, pc, b",
			Usage: "Number format for PC/Java-based Minecraft and most other NBT tools",
		},
		cli.StringFlag{
			Name:        "nbt-file, n",
			Value:       "-",
			Usage:       "NBT `FILE` path",
			Destination: &nbtFile,
		},
		// cli.StringFlag{
		// 	Name:        "json-file, j",
		// 	Value:       "-",
		// 	Usage:       "JSON `FILE` path",
		// 	Destination: &jsonFile,
		// },
		cli.IntFlag{
			Name:        "skip",
			Value:       0,
			Usage:       "Skip `NUM` bytes of NBT input. For MCPE level.dat, use --skip 8 to bypass header",
			Destination: &skipBytes,
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.String("big-endian") == "true" {
			byteOrder = binary.BigEndian
		} else {
			byteOrder = binary.LittleEndian
		}

		var myNbt, myJson, out []byte
		var err error

		if c.String("reverse") == "true" {
			myJson = []byte(`{"tagType": 1,"name": "Difficulty","value": 2}`)

			// fmt.Println(myJson)
			myNbt, err = nbt2json.Json2Nbt(myJson, binary.LittleEndian)
			if err != nil {
				return err
			}
			fmt.Println(string(myNbt[:]))
		} else {

			if c.String("nbt-file") == "-" {
				myNbt, err = ioutil.ReadAll(os.Stdin)
				if err != nil {
					return err
				}
			} else {
				f, err := os.Open(c.String("nbt-file"))
				if err != nil {
					return err
				}
				defer f.Close()
				myNbt, err = ioutil.ReadAll(f)
				if err != nil {
					return err
				}
			}

			// is it gzipped?
			if (myNbt[0] == 0x1f) && (myNbt[1] == 0x8b) {
				var uncompressed []byte
				buf := bytes.NewReader(myNbt)
				zr, err := gzip.NewReader(buf)
				if err != nil {
					return err
				}
				uncompressed, err = ioutil.ReadAll(zr)
				myNbt = uncompressed
			}
			out, err = nbt2json.Nbt2Json(myNbt[skipBytes:], byteOrder)
			if err != nil {
				return err
			}
			fmt.Println(string(out[:]))
		}
		return nil
	}

	app.Run(os.Args)
}
