#include <stdlib.h>

void show_usuage(void)
{
    printf("Usuage: run APP\\n");
}

int main(int argc, char** argv) 
{
    char command[256] = {0};

    if (argc != 2) {
        show_usuage();
        exit(0);
    }

    sprintf(command, "./%s", argv[1]);
    system(command);
    sprintf(command, "./pvtrace %s", argv[1]);
    system(command);
    sprintf(command, "dot -Tpng graph.dot -o graph.png");
    system(command);
    sprintf(command, "valgrind --tool=memcheck --leak-check=full --log-file=memcheck.log ./%s", argv[1]);
    system(command);
    sprintf(command, "valgrind --tool=callgrind ./%s > callgrind.log 2>&1", argv[1]);
    system(command);
    sprintf(command, "gprof ./%s > gprof.log 2>&1", argv[1]);
    system(command);
    sprintf(command, "mtrace ./%s %s> memleak_detect.log 2>&1", argv[1], "leak_detector.log");

    /*  Need to Add gcov support for code coverage test in future */
}

