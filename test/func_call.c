#include <stdio.h>

int func0() {
    printf("%s called\n", __func__);
    return 1;
}

int func1(int a) {
    printf("%s(%d) called\n", __func__, a);
    return a + 1;
}

int func2(int a, int b) {
    printf("%s(%d, %d) called\n", __func__, a, b);
    return a + b + 1;
}

int func3(int a, int b, int c) {
    printf("%s(%d, %d, %d) called\n", __func__, a, b, c);
    return a + b + c + 1;
}

int func4(int a, int b, int c, int d) {
    printf("%s(%d, %d, %d, %d) called\n", __func__, a, b, c, d);
    return a + b + c + d + 1;
}

int func5(int a, int b, int c, int d, int e) {
    printf("%s(%d, %d, %d, %d, %d) called\n", __func__, a, b, c, d, e);
    return a + b + c + d + e + 1;
}

int func6(int a, int b, int c, int d, int e, int f) {
    printf("%s(%d, %d, %d, %d, %d, %d) called\n", __func__, a, b, c, d, e, f);
    return a + b + c + d + e + f + 1;
}
