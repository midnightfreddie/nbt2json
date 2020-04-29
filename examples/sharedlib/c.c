#include <stdio.h>
#include "libnbt2json.h"

int main() {
    printf("Hello, C!\n");
    HelloDll();
    Json2Nbt("Hi from a parameter in C");
    printf("%s", Nbt2Json());
    return 0;
}

// gcc -Wall -L. -llibnbt2json c.c