# nbt2json

A command line utitlity and module that reads NBT data and converts it to JSON. Converting the other way is not yet implemented but intended.

## Features

- Command line utility will auto-detect and decompress gzipped files
- Can read both MCPE/Win10 and Java-based NBT data
    - Does not auto-detect which
    - Utility defaults to MCPE / little endian
- `--skip` parameter allows skipping headers. e.g. `nbt2json --skip 8 /path/to/level.dat` for MCPE's level.dat

## Why?

\<sigh\> Out of all the NBT tools out there, none of them seem to do what I want:

- Read MCPE (phones, tablets and Win10 app store edition) NBT data which use little endian encoding where the PC Java version uses big endian
- Convert back and forth to a human readable REST-API-able format
- Be includeable or as portable as my Go code

There are some Go options I could adapt for little endian, but I'll have to do enough work to make it convert to JSON and back that it's probably simpler to start from scratch, and I'm pretty sure I know how I want to do it.

## Dev notes

- I'm going to change the interface of `nbt2json.Nbt2Json` to use `[]byte` as input instead of a reader.
- The upcoming Json2Nbt function will not use the same struct; it will use an `interface{}` and encode based on the tagType fields. I had originally hoped to Marshal and Unmarshal to and from JSON and NBT, but my goal was to export to JSON, edit and then reencode. This way the struct doesn't have to match the data schema.
- My main motivation for this project is to convert to/from JSON and use any JSON editor to modify Minecraft PE data with [McpeTool](https://github.com/midnightfreddie/McpeTool), and to keep the read/write primitives in Go code while letting a client browser manage any validation to avoid having to re-release the read/write tools every time Minecraft changes formats.