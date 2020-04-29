$Signature = @'
[DllImport("libnbt2json.dll")]
public static extern void HelloDll();
'@

Add-Type -MemberDefinition $Signature -Namespace Nbt2Json -Name Lib

[Nbt2Json.Lib]::HelloDll()
