package main

import (
	"fmt"
	"os"

	"encoding/binary"

	"github.com/midnightfreddie/nbt2json"
	"github.com/urfave/cli"
)

func main() {
	var nbtFile, jsonFile string
	var byteOrder binary.ByteOrder
	app := cli.NewApp()
	app.Name = "NBT to JSON"
	app.Version = "0.0.0"
	app.Usage = "UNDER DEVELOPMENT, MOST OR ALL OPTIONS NOT IMPLEMENTED - Converts NBT-encoded data to JSON"
	app.Flags = []cli.Flag{
		// cli.BoolFlag{
		// 	Name:  "reverse, json2nbt, r",
		// 	Usage: "Convert JSON to NBT instead",
		// },
		cli.BoolTFlag{
			Name:  "little-endian, little, mcpe, l",
			Usage: "Number format for Minecraft Pocket Edition and Windows 10 Edition",
		},
		cli.BoolFlag{
			Name:  "big-endian, big, java, pc, b",
			Usage: "Number format for Minecraft Pocket Edition and Windows 10 Edition",
		},
		cli.StringFlag{
			Name:        "nbt-file, n",
			Value:       "-",
			Usage:       "NBT `FILE` path",
			Destination: &nbtFile,
		},
		cli.StringFlag{
			Name:        "json-file, j",
			Value:       "-",
			Usage:       "JSON `FILE` path",
			Destination: &jsonFile,
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.String("big-endian") == "true" {
			byteOrder = binary.BigEndian
		} else {
			byteOrder = binary.LittleEndian
		}
		myNbt := []byte{2, 2, 0, 'h', 'i', 0, 1}
		out, err := nbt2json.Nbt2Json(myNbt, byteOrder)
		if err != nil {
			return err
		}
		// fmt.Printf("%v\n", out)
		fmt.Println(string(out[:]))
		return nil
	}

	app.Run(os.Args)
}
