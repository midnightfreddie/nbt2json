libnbt2json is an experiment in making a shared library (.DLL, .so, or .dylib)
for use in other languages.

Build with `-buildmode=c-shared` and either rename the extensionless output with
.dll, .so, or .dylib as appropriate or use the `-o <filename>` option before the
build target.

The trick is going to be using C-native data types instead of Go data types for
passing values.