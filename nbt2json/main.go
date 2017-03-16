package main

import (
	"fmt"
	"os"

	"encoding/json"

	"github.com/midnightfreddie/nbt2json"
	"github.com/urfave/cli"
)

func main() {
	var nbtFile, jsonFile string
	app := cli.NewApp()
	app.Name = "NBT to JSON"
	app.Version = "0.0.0"
	app.Usage = "UNDER DEVELOPMENT, MOST OR ALL OPTIONS NOT IMPLEMENTED - Converts NBT-encoded data to JSON"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "reverse, json2nbt, r",
			Usage: "Convert JSON to NBT instead",
		}, cli.StringFlag{
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
		anObject := nbt2json.NewNbt2Json()
		fmt.Printf("%v\n", anObject)
		out, err := json.Marshal(anObject)
		if err != nil {
			return err
		}
		fmt.Println(string(out[:]))
		myString := "{ \"type\": 0, \"value\": \"Hi\"}"
		out2, err := nbt2json.Json2Nbt(myString)
		if err != nil {
			return err
		}
		fmt.Printf("%v\n", out2)
		return nil
	}

	app.Run(os.Args)
}
