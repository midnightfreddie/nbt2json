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