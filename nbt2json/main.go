package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"time"

	"encoding/binary"

	"compress/gzip"

	"github.com/midnightfreddie/nbt2json"
	"github.com/urfave/cli"
)

func main() {
	var inFile, outFile, comment string
	var byteOrder binary.ByteOrder
	var skipBytes int
	app := cli.NewApp()
	app.Name = "NBT to JSON"
	app.Version = "0.2.0-alpha"
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
		cli.StringFlag{
			Name:        "comment, c",
			Usage:       "Add `COMMENT` to json or yaml output, use quotes if contains white space",
			Destination: &comment,
		},
		cli.BoolTFlag{
			Name:  "little-endian, little, mcpe, l",
			Usage: "For Minecraft Pocket Edition and Windows 10 Edition (default)",
		},
		cli.BoolFlag{
			Name:  "big-endian, big, java, pc, b",
			Usage: "For PC/Java-based Minecraft and most other NBT tools",
		},
		cli.StringFlag{
			Name:        "in, i",
			Value:       "-",
			Usage:       "Input `FILE` path",
			Destination: &inFile,
		},
		cli.StringFlag{
			Name:        "out, o",
			Value:       "-",
			Usage:       "Output `FILE` path",
			Destination: &outFile,
		},
		cli.BoolFlag{
			Name:  "yaml, yml, y",
			Usage: "Use YAML instead of JSON",
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

		var inData, outData []byte
		var err error

		if inFile == "-" {
			inData, err = ioutil.ReadAll(os.Stdin)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
		} else {
			inData, err = ioutil.ReadFile(inFile)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
		}

		if c.String("reverse") == "true" {
			if c.String("yaml") == "true" {
				outData, err = nbt2json.Yaml2Nbt(inData, byteOrder)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
			} else {
				outData, err = nbt2json.Json2Nbt(inData, byteOrder)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
			}
		} else {
			// is it gzipped?
			if (inData[0] == 0x1f) && (inData[1] == 0x8b) {
				var uncompressed []byte
				buf := bytes.NewReader(inData)
				zr, err := gzip.NewReader(buf)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				uncompressed, err = ioutil.ReadAll(zr)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				inData = uncompressed
			}
			if c.String("yaml") == "true" {
				outData, err = nbt2json.Nbt2Yaml(inData[skipBytes:], byteOrder, comment)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
			} else {
				outData, err = nbt2json.Nbt2Json(inData[skipBytes:], byteOrder, comment)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
			}
		}
		if outFile == "-" {
			err = binary.Write(os.Stdout, binary.LittleEndian, outData)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
		} else {
			err = ioutil.WriteFile(outFile, outData, 0644)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
		}
		return nil
	}

	app.Run(os.Args)
}
