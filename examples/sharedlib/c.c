#include <stdio.h>
#include "libnbt2json.h"

int main() {
    printf("Hello, C!\n");
    
    HelloDll();
    
    Json2Nbt("Hi from a parameter in C");

    // Little-endian NBT with a compound (10) tag containing one short (2) tag named "SleepTimer" with value 0
    char byteArray[] = { 0x0a, 0x00 , 0x00, 0x02, 0x0a, 0x00, 0x53, 0x6c, 0x65, 0x65, 0x70, 0x54, 0x69, 0x6d , 0x65, 0x72, 0x00, 0x00 , 0x00 };
    printf("%s", Nbt2Json(byteArray, sizeof(byteArray)));
    return 0;
}

// gcc -Wall -L. -llibnbt2json c.c