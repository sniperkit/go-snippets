// comment initial
// #include <testcommentinclude.h>
#define ZERO 0
#define SOMEMACRO(x) x*x
#define POLOMPOS 0
#include <stdio.h>
#include "missing.h"

int main(int i, char *argv[]) {
    // comment before
    return ZERO;
    // comment after
}

namespace my_namespace {

bool valid_function() { return true; /* comment at the end of the functionline */ } /* comment after the closing of the functionline */

int another(int param) {
    auto j = param + 50; // comment after semicolon
    auto x = SOMEMACRO(j);
    // comment inside "another" function
    int z = SOMEMACRO(x);
    int l = imported_function(z, j, f);
    /* comment at start of line */ auto validstuff = j + 10;
    THISMUSTFAIL;

    if (valid_function()) {
        auto insideValidIf /* comment inside a line */ = 10;
#if POLOMPOS
        auto definedInsidePreprocIf = 50;
#else
        auto definedInsidePreprocElse = 50;
#endif
    }

    if (imported_function(z)) {
        /*
         * block comment inside if imported_function
         */
        char insideInvalidIf1 = 'a';
        long insideInvalidIf2 = 98765;
    }

    return j;
}
// comment at the end of the namespace
}
// comment at the end of the file
