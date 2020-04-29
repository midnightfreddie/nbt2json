<#
    Powershell can use P/Invoke to load non-.NET DLL files

    Place the DLL in the executable search path (the local folder is the easiest
    option)

    A .NET signature must be defined. It does not use the .h file, but you can
    look at the .h file for guidance. "public" must be added in the case of
    Powershell as we want the functions to be available in the script which
    is external to the type class.
#>

$Signature = @'
[DllImport("libnbt2json.dll")]
public static extern void HelloDll();
'@

Add-Type -MemberDefinition $Signature -Namespace Nbt2Json -Name Lib

[Nbt2Json.Lib]::HelloDll()
