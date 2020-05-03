## v0.4.0

Breaking changes!

For everyone:

- NBT Long values are now stored as valueLeast & valueMost **32-bit unsigned
integer pairs** in JSON. This is to prevent conversion issues across various
languages' JSON libaries, but it will make a little more work if you need to
modify these values.
  - Example JSON for NBT long

        {
          "nbt": [
            {
              "tagType": 4,
              "name": "LongAsUint32Pair",
              "value": {
                "valueLeast": 4294967295,
                "valueMost": 2147483647
              }
            }
          ]
        }

  - Example JSON for NBT long as string

        {
          "nbt": [
            {
              "tagType": 4,
              "name": "LongAsString",
              "value": "9223372036854775807"
            }
          ]
        }

- Optionally you can store long values as strings in JSON which may be more
convenient depending on the use case
- Converting from JSON will accept either strings or the valueLeast/valueMost
pair automatically

For utility executable users:

- Some parameters have been removed or renamed
  - Since Bedrock is default, all the little-endian parameters are gone
  - `--json2nbt`, `--big`, and `--pc` were removed, but their aliases remain
  - Command line library has been updated to the latest version
  - Added `--long-as-string` and short alias `-l` to cause NBT long values
  (64-bit integers) to be numbers-in-strings in the JSON output instead of
  valueLeast & valueMost numbers

For devs:

- byteOrder arguments are gone from all functions
- Use `UseJavaEncoding()` if you need big-endian / Java Edition encoding.
`UseBedrockEncoding()` is set by default, but you can call it to switch back if
you switched to Java previously.
- Use `UseLongAsString()` if you want NBT long int64's to be in the JSON as
strings. `UseLongAsUint32Pair()` is set by default, but you can call it to
switch back if you set the string method previously.
- Go tests are more thorough

## v0.3.4

This version has no data differences from v.0.3.3. It just has improved error
messaging and Go tests. I wanted to make this release mostly for those who would
use it as a module so they have the new improvements before I make some breaking
changes for the next version.

- Added Go tests in nbt2json_test.go
- Corrected mislabeled error messages
- Added errors for missing "nbt" in json and out-of-range numbers
- Added input value to error message output
- Refactored module code to include commmon.go for possibly interesting exports
- Executable now pulls version and url from the module values in common.go

## v0.3.3

- Merged VADemon's fix for int64s
- Int64s (NBT "long"s) will now export and import properly; they were previously getting messed up which impacted UUIDs as well
- However, other languages' JSON libraries may not handle int64s in JSON; this will be addressed in future release

## v0.3.2

- Fixed emtpy array showing as `null` in json instead of `[]` for compound tag 10
- Moved nbt2json/ main executable to cmd/nbt2json/
- Added `--gzip` / `-z` option to executable to compress output

There does not seem to be a missing data type, but users are reporting UUIDs
are not preserving due to float interpretation of json numbers.

## v0.3.1

- Added go.mod and go.sum to help manage dependency versions

There is ~~at least one new NBT data type and~~ one bug since the last update,
but these are not yet addressed. I am just adding go.mod to prevent possible
dependenccy version issues. I hope to address the issues soon.