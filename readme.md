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

## Help screen

By defualt, the nbt2json executable waits for input from stdin, so you need to `nbt2json -h` to see the help screen.

```
NAME:
   NBT to JSON - Converts NBT-encoded data to JSON | https://github.com/midnightfreddie/nbt2json

USAGE:
   nbt2json.exe [global options] command [command options] [arguments...]

VERSION:
   0.4.0-alpha

AUTHOR:
   Jim Nelson <jim@jimnelson.us>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --reverse, -r                  Convert JSON to NBT instead (default: false)
   --gzip, -z                     Compress output with gzip (default: false)
   --comment COMMENT, -c COMMENT  Add COMMENT to json or yaml output, use quotes if contains white space
   --big-endian, --java, -b       Use for Minecraft Java Edition (like most other NBT tools) (default: false)
   --in FILE, -i FILE             Input FILE path (default: "-")
   --out FILE, -o FILE            Output FILE path (default: "-")
   --yaml, --yml, -y              Use YAML instead of JSON (default: false)
   --skip NUM                     Skip NUM bytes of NBT input. For Bedrock's level.dat, use --skip 8 to bypass header (default: 0)
   --help, -h                     show help (default: false)
   --version, -v                  print the version (default: false)

COPYRIGHT:
   (c) 2018, 2019, 2020 Jim Nelson
```

## Dev notes

- Client Go code needs to `import "github.com/midnightfreddie/nbt2json"`
- Defaults to little-endian encoding for Bedrock Edition. Call `nbt2json.UseJavaEncoding()` and `nbt2json.UseBedrockEncoding()` to change encoding mode for as long as the module is open.
- The functions use byte arrays where you might expect strings. Convert as such: `var myString = someByteArray[:]` or `var myByteArray = []byte(someStringValue)`
- All errors should bubble up through the error part of the result and should describe where the problem was
- The Json2Nbt function uses an `interface{}` and encodes based on the tagType fields. I had originally hoped to Marshal and Unmarshal to and from JSON and NBT, but my goal was to export to JSON, edit and then reencode. This way the struct doesn't have to match the data schema.
- My main motivation for this project is to convert to/from JSON and use any JSON editor to modify Minecraft PE data with [McpeTool](https://github.com/midnightfreddie/McpeTool), and to keep the read/write primitives in Go code while letting a client browser manage any validation to avoid having to re-release the read/write tools every time Minecraft changes formats.

### Exported Go Functions

- **Nbt2Yaml** converts uncompressed NBT byte array to YAML byte array

		func Nbt2Yaml(b []byte, comment string) ([]byte, error)

- **Nbt2Json** converts uncompressed NBT byte array to JSON byte array

		func Nbt2Json(b []byte, comment string) ([]byte, error)

- **Yaml2Nbt** converts JSON byte array to uncompressed NBT byte array (Hint: You can just use this for both JSON *and* YAML if you like since JSON is a valid subeset of YAML)

		func Yaml2Nbt(b []byte) ([]byte, error)

- **Json2Nbt** converts JSON byte array to uncompressed NBT byte array

		func Json2Nbt(b []byte) ([]byte, error)

- **UseJavaEndoding** sets any nbt encoding/decoding to big-endian to match Minecraft Java Edition

        func UseJavaEncoding()

- **UseBedrockEncoding** sets nbt encoding/decoding to little-endian to match Minecraft Bedrock Edition

        func UseBedrockEncoding()

Other exports of possible interest are in common.go.
