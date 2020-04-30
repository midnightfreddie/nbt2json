<#
    Powershell can use P/Invoke to load non-.NET DLL files

    Place the DLL (or .so or .dylib) in the executable search path (the local folder is the easiest
    option)

    A .NET signature must be defined. It does not use the .h file, but you can
    look at the .h file for guidance. "public" must be added in the case of
    Powershell as we want the functions to be available in the script which
    is external to the type class.
#>

$LibBase = "libnbt2json"

$SharedLib = "${LibBase}.dll"
if ($IsLinux) { $SharedLib = "${LibBase}.so" }
if ($IsMacOS) { $SharedLib = "${LibBase}.dylib" }

$Signature = @"
[DllImport("${SharedLib}", CharSet = CharSet.Ansi)]
public static extern void HelloDll();
"@
# [DllImport("libnbt2json.dll", CharSet = CharSet.Ansi)]
# public static extern void Json2Nbt(string cString);
# [DllImport("libnbt2json.dll", CharSet = CharSet.Ansi)]
# public static extern IntPtr Nbt2Json();
# '@

# Add-Type -MemberDefinition $Signature -Namespace Nbt2Json -Name Lib

# [Nbt2Json.Lib]::HelloDll()
# [Nbt2Json.Lib]::Json2Nbt("Hello from a parameter")
# [System.Runtime.InteropServices.Marshal]::PtrToStringAnsi([Nbt2Json.Lib]::Nbt2Json())


$Source = @"
using System;
using System.Runtime.InteropServices;

namespace Nbt2Json
{
    public class Lib
    {
        const string libName = "${SharedLib}";

        public struct GoSlice
        {
            public IntPtr data;
            public long len, cap;
            public GoSlice(IntPtr data, long len, long cap)
            {
                this.data = data;
                this.len = len;
                this.cap = cap;
            }
        }
        public struct GoString
        {
            public string msg;
            public long len;
            public GoString(string msg, long len)
            {
                this.msg = msg;
                this.len = len;
            }
        }

        [DllImport(libName, CharSet = CharSet.Ansi)]
        public static extern void HelloDll();

        [DllImport(libName, CharSet = CharSet.Ansi)]
        public static extern GoSlice SomeByteArray();
    }
}
"@

$Type = Add-Type -TypeDefinition $Source -PassThru

[Nbt2Json.Lib]::HelloDll()

$Foo = [Nbt2Json.Lib]::SomeByteArray()
$Foo.len
