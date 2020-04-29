#include <stdio.h>
#include "libnbt2json.h"

int main() {

    // Little-endian NBT with a compound (10) tag containing one short (2) tag named "SleepTimer" with value 0
    char byteArray[] = { 0x0a, 0x00 , 0x00, 0x02, 0x0a, 0x00, 0x53, 0x6c, 0x65, 0x65, 0x70, 0x54, 0x69, 0x6d , 0x65, 0x72, 0x00, 0x00 , 0x00 };

    char yamlString[] =
        "nbt:\n"
        "- name: \"\"\n"
        "  tagType: 10\n"
        "  value:\n"
        "  - name: Test\n"
        "    tagType: 2\n"
        "    value: 256\n"
        ;

    HelloDll();
    
    // Json2Nbt("Hi from a parameter in C");
    Yaml2Nbt(yamlString);

    printf("%s\n", Nbt2Json(byteArray, sizeof(byteArray)));

    printf("%s\n", Nbt2Yaml(byteArray, sizeof(byteArray)));

    return 0;
}

// gcc -Wall -L. -llibnbt2json c.c