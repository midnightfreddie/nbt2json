# nbt2json

A command line utility and Go module that reads NBT data and converts it to JSON or YAML for editing and then back to NBT.

## Features

- nbt2json executable will auto-detect and decompress gzipped files
- nbt2json executable has option to gzip output
- Can read and write both Minecraft Bedrock Edition and Java Edition NBT data
    - Does **not** auto-detect which
    - nbt2json executable defaults to Bedrock Edition / little endian
- Can import to other Go projects
- Can use either JSON or YAML
- Can include comment in JSON/YAML output (which is ignored when converting back to NBT)

## Known Issues

- "Long" NBT values which are 64-bit integers may not be properly preserved in some languages' JSON libraries. This will be fixed in a future release by breaking 64-bit integers into high/low 32-bit integer pairs in the JSON. As of nbt2json v0.3.3 they will at least export and import with the correct values when unaltered or altered manually in a text editor.

## Help screen

By defualt, the nbt2json executable waits for input from stdin, so you need to `nbt2json -h` to see the help screen.

```
NAME:
   NBT to JSON - Converts NBT-encoded data to JSON | https://github.com/midnightfreddie/nbt2json

USAGE:
   nbt2json.exe [global options] command [command options] [arguments...]

VERSION:
   0.3.4

AUTHOR:
   Jim Nelson <jim@jimnelson.us>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --reverse, --json2nbt, -r              Convert JSON to NBT instead
   --gzip, -z                             Compress output with gzip
   --comment COMMENT, -c COMMENT          Add COMMENT to json or yaml output, use quotes if contains white space
   --little-endian, --little, --mcpe, -l  For Minecraft Bedrock Edition (Pocket and Windows 10) (default)
   --big-endian, --big, --java, --pc, -b  For Minecraft Java Edition (like most other NBT tools)
   --in FILE, -i FILE                     Input FILE path (default: "-")
   --out FILE, -o FILE                    Output FILE path (default: "-")
   --yaml, --yml, -y                      Use YAML instead of JSON
   --skip NUM                             Skip NUM bytes of NBT input. For Bedrock's level.dat, use --skip 8 to bypass header (default: 0)
   --help, -h                             show help
   --version, -v                          print the version

COPYRIGHT:
   (c) 2018, 2019, 2020 Jim Nelson
```

## Dev notes

- Client Go code needs to `import "github.com/midnightfreddie/nbt2json"`
- For `binary.byteOrder` parameters, pass `nbt2json.Bedrock` for Bedrock Edition or `nbt2json.Java` for Java Edition (aliases for `binary.LittleEndian` and `binary.BigEndian`, respectively)
- The functions use byte arrays where you might expect strings. Convert as such: `var myString = someByteArray[:]` or `var myByteArray = []byte(someStringValue)`
- All errors should bubble up through the error part of the result and should describe where the problem was
- The Json2Nbt function uses an `interface{}` and encodes based on the tagType fields. I had originally hoped to Marshal and Unmarshal to and from JSON and NBT, but my goal was to export to JSON, edit and then reencode. This way the struct doesn't have to match the data schema.
- My main motivation for this project is to convert to/from JSON and use any JSON editor to modify Minecraft PE data with [McpeTool](https://github.com/midnightfreddie/McpeTool), and to keep the read/write primitives in Go code while letting a client browser manage any validation to avoid having to re-release the read/write tools every time Minecraft changes formats.

### Exported Go Functions

- **Nbt2Yaml** converts uncompressed NBT byte array to YAML byte array

		func Nbt2Yaml(b []byte, byteOrder binary.ByteOrder, comment string) ([]byte, error)

- **Nbt2Json** converts uncompressed NBT byte array to JSON byte array

		func Nbt2Json(b []byte, byteOrder binary.ByteOrder, comment string) ([]byte, error)

- **Yaml2Nbt** converts JSON byte array to uncompressed NBT byte array (Hint: You can just use this for both JSON *and* YAML if you like since JSON is a valid subeset of YAML)

		func Yaml2Nbt(b []byte, byteOrder binary.ByteOrder) ([]byte, error)

- **Json2Nbt** converts JSON byte array to uncompressed NBT byte array

		func Json2Nbt(b []byte, byteOrder binary.ByteOrder) ([]byte, error)

Other exports of possible interest are in common.go.
