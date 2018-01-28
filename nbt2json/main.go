package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"encoding/binary"

	"compress/gzip"

	"github.com/midnightfreddie/nbt2json"
	"github.com/urfave/cli"
)

func main() {
	var nbtFile, jsonFile string
	var byteOrder binary.ByteOrder
	var skipBytes int
	app := cli.NewApp()
	app.Name = "NBT to JSON"
	app.Version = "0.1.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Jim Nelson",
			Email: "jim@jimnelson.us",
		},
	}
	app.Copyright = "(c) 2018 Jim Nelson"
	app.Usage = "Converts NBT-encoded data to JSON | https://github.com/midnightfreddie/nbt2json"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "reverse, json2nbt, r",
			Usage: "Convert JSON to NBT instead",
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
		cli.StringFlag{
			Name:        "json-file, j",
			Value:       "-",
			Usage:       "JSON `FILE` path",
			Destination: &jsonFile,
		},
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
			if c.String("json-file") == "-" {
				myJson, err = ioutil.ReadAll(os.Stdin)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
			} else {
				f, err := os.Open(c.String("json-file"))
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				defer f.Close()
				myJson, err = ioutil.ReadAll(f)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
			}
			myNbt, err = nbt2json.Json2Nbt(myJson, byteOrder)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			if c.String("nbt-file") == "-" {
				err = binary.Write(os.Stdout, binary.LittleEndian, myNbt)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
			} else {
				err = ioutil.WriteFile(c.String("nbt-file"), myNbt, 0644)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
			}
		} else {

			if c.String("nbt-file") == "-" {
				myNbt, err = ioutil.ReadAll(os.Stdin)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
			} else {
				f, err := os.Open(c.String("nbt-file"))
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				defer f.Close()
				myNbt, err = ioutil.ReadAll(f)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
			}

			// is it gzipped?
			if (myNbt[0] == 0x1f) && (myNbt[1] == 0x8b) {
				var uncompressed []byte
				buf := bytes.NewReader(myNbt)
				zr, err := gzip.NewReader(buf)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				uncompressed, err = ioutil.ReadAll(zr)
				myNbt = uncompressed
			}
			out, err = nbt2json.Nbt2Json(myNbt[skipBytes:], byteOrder)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			fmt.Println(string(out[:]))
		}
		return nil
	}

	app.Run(os.Args)
}
