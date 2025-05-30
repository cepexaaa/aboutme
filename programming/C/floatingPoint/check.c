#include "return_codes.h"
#include <inttypes.h>
#include <stdio.h>
#include <alg.h>
int main() {
    uint8_t a = 10;
    uint8_t b = ~a;
    printf("%i\n", (uint16_t) b);
    printf("%i", b);
}
