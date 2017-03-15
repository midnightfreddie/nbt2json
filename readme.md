# nbt2json

It's not ready yet.

## Why?

\<sigh\> Out of all the NBT tools out there, none of them seem to do what I want:

- Read MCPE (phones, tablets and Win10 app store edition) NBT data which use little endian encoding where the PC Java version uses big endian
- Convert back and forth to a human readable REST-API-able format
- Be includeable or as portable as my Go code

There are some Go options I could adapt for little endian, but I'll have to do enough work to make it convert to JSON and back that it's probably simpler to start from scratch, and I'm pretty sure I know how I want to do it.