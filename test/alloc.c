#include <stdlib.h>

int alloc4(int **p, int a, int b, int c, int d) {
    int *ptr = malloc(4 * sizeof(int));
    ptr[0] = a;
    ptr[1] = b;
    ptr[2] = c;
    ptr[3] = d;
    *p = ptr;
    return 0;
}
