#include <stdio.h>
#include "libnbt2json.h"

int main() {
    printf("Hello, C!\n");
    
    HelloDll();
    
    Json2Nbt("Hi from a parameter in C");

    char byteArray[] = { 5, 6, 7, 8 };
    printf("%s", Nbt2Json(byteArray, sizeof(byteArray)));
    return 0;
}

// gcc -Wall -L. -llibnbt2json c.c