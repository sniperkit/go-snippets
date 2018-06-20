# genericPointer
Is interface{} a generic pointer?

# What is the function point in GO ?

In C:
```
#include <stdio.h>

typedef int (*FP)(int, int);

int add(int a, int b) {
    return a+ b;
}

int main() {
    FP fp1 = add;
    printf("1 + 2 = %d", (*fp1)(1, 2));
    printf("\n");
    printf("1 + 2 = %d", fp1(1, 2));
    printf("\n");

    // Another way
    FP fp2 = &add;
    printf("1 + 2 = %d", (*fp2)(1, 2));
    printf("\n");
    printf("1 + 2 = %d", fp2(1, 2));
    printf("\n");

    return 0;
}
```

What is the **function point** in GO?
