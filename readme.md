# nbt2json

A command line utility and module that reads NBT data and converts it to JSON or YAML for editing and then back to NBT.

## Features

- Command line utility will auto-detect and decompress gzipped files
- Can read both MCPE/Win10 and Java-based NBT data
    - Does not auto-detect which
    - Utility defaults to MCPE / little endian
- `--skip` parameter allows skipping headers. e.g. `nbt2json --skip 8 /path/to/level.dat` for MCPE's level.dat
- `--yaml` parameter converts to/from YAML instead of JSON
- `--comment "My comment here"` parameter allows adding a comment field to the JSON/YAML output
- `nbt2json -h` for usage info

## Help screen

	NAME:
	   NBT to JSON - Converts NBT-encoded data to JSON | https://github.com/midnightfreddie/nbt2json

	USAGE:
	   nbt2json.exe [global options] command [command options] [arguments...]

	VERSION:
	   0.2.0

	AUTHOR:
	   Jim Nelson <jim@jimnelson.us>

	COMMANDS:
		 help, h  Shows a list of commands or help for one command

	GLOBAL OPTIONS:
	   --reverse, --json2nbt, -r              Convert JSON to NBT instead
	   --comment COMMENT, -c COMMENT          Add COMMENT to json or yaml output, use quotes if contains white space
	   --little-endian, --little, --mcpe, -l  For Minecraft Pocket Edition and Windows 10 Edition (default)
	   --big-endian, --big, --java, --pc, -b  For PC/Java-based Minecraft and most other NBT tools
	   --in FILE, -i FILE                     Input FILE path (default: "-")
	   --out FILE, -o FILE                    Output FILE path (default: "-")
	   --yaml, --yml, -y                      Use YAML instead of JSON
	   --skip NUM                             Skip NUM bytes of NBT input. For MCPE level.dat, use --skip 8 to bypass header (default: 0)
	   --help, -h                             show help
	   --version, -v                          print the version

	COPYRIGHT:
	   (c) 2018 Jim Nelson

## Why?

\<sigh\> Out of all the NBT tools out there, none of them seem to do what I want:

- Read MCPE (phones, tablets and Win10 app store edition) NBT data which use little endian encoding where the PC Java version uses big endian
- Convert back and forth to a human readable REST-API-able format
- Be includeable or as portable as my Go code

There are some Go options I could adapt for little endian, but I'll have to do enough work to make it convert to JSON and back that it's probably simpler to start from scratch, and I'm pretty sure I know how I want to do it.

## Dev notes

- The Json2Nbt function uses an `interface{}` and encodes based on the tagType fields. I had originally hoped to Marshal and Unmarshal to and from JSON and NBT, but my goal was to export to JSON, edit and then reencode. This way the struct doesn't have to match the data schema.
- My main motivation for this project is to convert to/from JSON and use any JSON editor to modify Minecraft PE data with [McpeTool](https://github.com/midnightfreddie/McpeTool), and to keep the read/write primitives in Go code while letting a client browser manage any validation to avoid having to re-release the read/write tools every time Minecraft changes formats.