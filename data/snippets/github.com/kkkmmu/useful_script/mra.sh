#!/bin/sh

if [ $# -ne 1 ]; then
    echo "Usuage: gnew.sh ProjectName"
    exit -1;
fi

project=$1;
package=`echo $project | tr '[:upper:]' '[:lower:]'`

mkdir $1
cd $1

# Create Project class
mkdir -p src
mkdir -p build
mkdir -p lib
mkdir -p app
mkdir -p tools
mkdir -p test

touch CMakeLists.txt
echo "CMAKE_MINIMUM_REQUIRED(VERSION 2.8)"                                      > CMakeLists.txt
echo "PROJECT($1)"                                                              >> CMakeLists.txt
echo "ADD_SUBDIRECTORY(src .src)"                                               >> CMakeLists.txt
echo "ADD_SUBDIRECTORY(app .app)"                                               >> CMakeLists.txt
echo "ADD_SUBDIRECTORY(lib .lib)"                                               >> CMakeLists.txt
echo "ADD_SUBDIRECTORY(tools .tools)"                                           >> CMakeLists.txt
echo "SET(CMAKE_VERBOSE_MAKEFILE on)"                                           >> CMakeLists.txt
echo "ADD_SUBDIRECTORY(test .test)"                                             >> CMakeLists.txt
echo "SET(CFLAGS \"-Wall -Werror\")"                                            >> CMakeLists.txt

touch lib/CMakeLists.txt
echo "ADD_SUBDIRECTORY(unity)"                                                  > lib/CMakeLists.txt
# echo "ADD_SUBDIRECTORY(mysql)"                                                >> lib/CMakeLists.txt
#echo "ADD_SUBDIRECTORY(libds)"                                                  >> lib/CMakeLists.txt
# echo "ADD_SUBDIRECTORY(collection)"                                             >> lib/CMakeLists.txt
# echo "ADD_SUBDIRECTORY(onion)"                                                >> lib/CMakeLists.txt
echo "ADD_SUBDIRECTORY(llog)"                                                   >> lib/CMakeLists.txt
echo "ADD_SUBDIRECTORY(libfault)"                                               >> lib/CMakeLists.txt
# echo "ADD_SUBDIRECTORY(libmill)"                                                >> lib/CMakeLists.txt
echo "ADD_SUBDIRECTORY(memwatch)"                                               >> lib/CMakeLists.txt
# echo "ADD_SUBDIRECTORY(libapue)"                                                >> lib/CMakeLists.txt
# echo "ADD_SUBDIRECTORY(libunp)"                                                 >> lib/CMakeLists.txt
echo "ADD_SUBDIRECTORY(libran)"                                                 >> lib/CMakeLists.txt
echo "SET(CMAKE_VERBOSE_MAKEFILE on)"                                           >> lib/CMakeLists.txt

touch src/$1.c
echo "#include <stdio.h>"                                                       > src/$1.c
echo "#include <random.h>"                                                      > src/$1.c
echo "#include <llog.h>"                                                        > src/$1.c
echo "void $1(void)"                                                            >> src/$1.c
echo "{"                                                                        >> src/$1.c
echo "      printf(\"Hello world!\\\n\");"                                      >> src/$1.c
echo "}"                                                                        >> src/$1.c
touch src/$1.h
echo "#ifndef __$1_H__"                                                         > src/$1.h
echo "#define __$1_H__"                                                         >> src/$1.h
echo "void $1(void);"                                                           >> src/$1.h
echo "#endif /* __$1_H__ */"                                                    >> src/$1.h
touch src/CMakeLists.txt

echo "AUX_SOURCE_DIRECTORY(\${CMAKE_CURRENT_SOURCE_DIR} LIB_SRC)"               > src/CMakeLists.txt
echo "ADD_LIBRARY(kkkmmu SHARED \${LIB_SRC})"                                   >> src/CMakeLists.txt
echo "ADD_LIBRARY(kkkmmu_static STATIC \${LIB_SRC})"                            >> src/CMakeLists.txt
echo "SET(CMAKE_C_FLAGS \"-Wall -Werror -g -pg -finstrument-functions -ftest-coverage -fprofile-arcs -funwind-tables -rdynamic\")"                                                                                >> src/CMakeLists.txt
echo "SET(LIBRARY_OUTPUT_PATH \${PROJECT_BINARY_DIR}/lib)"                      >> src/CMakeLists.txt
echo "SET_TARGET_PROPERTIES (kkkmmu_static PROPERTIES OUTPUT_NAME kkkmmu)"      >> src/CMakeLists.txt
echo "SET_TARGET_PROPERTIES (kkkmmu PROPERTIES CLEAN_DIRECT_OUTPUT 1)"          >> src/CMakeLists.txt
echo "SET_TARGET_PROPERTIES (kkkmmu PROPERTIES CLEAN_DIRECT_OUTPUT 1)"          >> src/CMakeLists.txt
echo "SET_TARGET_PROPERTIES (kkkmmu PROPERTIES VERSION 1.1 SOVERSION 1)"        >> src/CMakeLists.txt
echo "INCLUDE_DIRECTORIES(\${PROJECT_SOURCE_DIR}/src \${PROJECT_SOURCE_DIR}/lib/llog/src \${PROJECT_SOURCE_DIR}/lib/libfault \${PROJECT_SOURCE_DIR}/lib/libran/src)"                                              >> src/CMakeLists.txt
echo "# install (TARGETS $1 $1_static "                                         >> src/CMakeLists.txt
echo "#     LIBRARY DESTINATION \${PROJECT_BINARY_DIR}/lib"                     >> src/CMakeLists.txt
echo "#     ARCHIVE DESTINATION \${PROJECT_BINARY_DIR}/lib)"                    >> src/CMakeLists.txt
echo "INSTALL (FILES $1.h DESTINATION \${PROJECT_BINARY_DIR}/include/$1)"       >> src/CMakeLists.txt

touch app/$1Main.c
echo "#include \"$1.h\""                                                        > app/$1Main.c
echo "#include \"fault.h\""                                                     >> app/$1Main.c
echo "#include \"leak_detector.h\""                                             >> app/$1Main.c
echo ""                                                                         >> app/$1Main.c
echo "int main(int argc, char** argv)"                                          >> app/$1Main.c
echo "{"                                                                        >> app/$1Main.c
echo "    register_fault_handlers();"                                           >> app/$1Main.c
echo "    memleak_detect_start();"                                              >> app/$1Main.c
echo "    $1();"                                                                >> app/$1Main.c
echo "    memleak_detect_stop();"                                               >> app/$1Main.c
echo "    return 0;"                                                            >> app/$1Main.c
echo "}"                                                                        >> app/$1Main.c

touch app/CMakeLists.txt
echo "AUX_SOURCE_DIRECTORY(\${CMAKE_CURRENT_SOURCE_DIR} APP_SRC)"               > app/CMakeLists.txt
echo "ADD_EXECUTABLE($1 \${APP_SRC})"                                           >> app/CMakeLists.txt
echo "SET(EXECUTABLE_OUTPUT_PATH \${PROJECT_BINARY_DIR}/bin)"                   >> app/CMakeLists.txt
echo "LINK_DIRECTORIES (\${PROJECT_BINARY_DIR}/lib)"                            >> app/CMakeLists.txt
echo "TARGET_LINK_LIBRARIES($1 kkkmmu fault llog libran)"                       >> app/CMakeLists.txt
echo "SET(CMAKE_C_FLAGS \"-Wall -Werror -g -pg -finstrument-functions -ftest-coverage -fprofile-arcs -funwind-tables -rdynamic\")"                                      >> app/CMakeLists.txt
echo "INCLUDE_DIRECTORIES(\${PROJECT_SOURCE_DIR}/src \${CMAKE_CURRENT_SOURCE_DIR} \${PROJECT_SOURCE_DIR}/lib/libfault \${PROJECT_SOURCE_DIR}/lib/llog/src \${PROJECT_SOURCE_DIR}/lib/libran/src)"                >> app/CMakeLists.txt

touch test/$1Test.c
echo "#include \"$1.h\""                                                        > test/$1Test.c
echo "#include <llog.h>"                                                        >> test/$1Test.c
echo "#include <random.h>"                                                      >> test/$1Test.c
echo "#include <unity.h>"                                                       >> test/$1Test.c
echo "#include \"fault.h\""                                                     >> test/$1Test.c
echo "#include \"unity_internals.h\""                                           >> test/$1Test.c
echo "#include \"leak_detector.h\""                                             >> test/$1Test.c
echo ""                                                                         >> test/$1Test.c
echo ""                                                                         >> test/$1Test.c
echo "void test_failure_func(void)"                                             >> test/$1Test.c
echo "{"                                                                        >> test/$1Test.c
echo "      TEST_ASSERT_EQUAL_STRING(\"ello\", \"Hello\");"                     >> test/$1Test.c
echo "}"                                                                        >> test/$1Test.c
echo ""                                                                         >> test/$1Test.c
echo "void test_success_func(void)"                                             >> test/$1Test.c
echo "{"                                                                        >> test/$1Test.c
echo "      TEST_ASSERT_EQUAL_STRING(\"Hello\", \"Hello\");"                    >> test/$1Test.c
echo "}"                                                                        >> test/$1Test.c
echo ""                                                                         >> test/$1Test.c
echo "void test_ignore_func(void)"                                              >> test/$1Test.c
echo "{"                                                                        >> test/$1Test.c
echo "      TEST_IGNORE_MESSAGE(\"This is a ignored test!\");"                  >> test/$1Test.c
echo "}"                                                                        >> test/$1Test.c
echo ""                                                                         >> test/$1Test.c
echo "int main(int argc, char** argv)"                                          >> test/$1Test.c
echo "{"                                                                        >> test/$1Test.c
echo "    register_fault_handlers();"                                           >> test/$1Test.c
echo "    memleak_detect_start();"                                              >> test/$1Test.c
echo "    llog(INFO, \"%s\", \"Start the test for Project tt!\");"              >> test/$1Test.c
echo "    UNITY_BEGIN();"                                                       >> test/$1Test.c
echo ""                                                                         >> test/$1Test.c
echo "    RUN_TEST(test_failure_func);"                                         >> test/$1Test.c
echo "    RUN_TEST(test_success_func);"                                         >> test/$1Test.c
echo "    RUN_TEST(test_ignore_func);"                                          >> test/$1Test.c
echo ""                                                                         >> test/$1Test.c
echo "    UNITY_END();"                                                         >> test/$1Test.c
echo "    llog(INFO, \"%s\", \"Test For Project tt is Finished!\");"            >> test/$1Test.c
echo "    memleak_detect_stop();"                                               >> test/$1Test.c
echo "    return 0;"                                                            >> test/$1Test.c
echo " }"                                                                       >> test/$1Test.c

touch test/CMakeLists.txt
echo "AUX_SOURCE_DIRECTORY(\${CMAKE_CURRENT_SOURCE_DIR} APP_SRC)"               > test/CMakeLists.txt
echo "ADD_EXECUTABLE($1Test \${APP_SRC})"                                       >> test/CMakeLists.txt
echo "SET(EXECUTABLE_OUTPUT_PATH \${PROJECT_BINARY_DIR}/bin)"                   >> test/CMakeLists.txt
echo "LINK_DIRECTORIES (\${PROJECT_BINARY_DIR}/lib)"                            >> test/CMakeLists.txt
echo "TARGET_LINK_LIBRARIES($1Test kkkmmu unity llog fault libran)"             >> test/CMakeLists.txt
echo "SET(CMAKE_C_FLAGS \"-Wall -Werror -g -pg -finstrument-functions -ftest-coverage -fprofile-arcs -funwind-tables -rdynamic\")"                                      >> test/CMakeLists.txt
echo "INCLUDE_DIRECTORIES(\${PROJECT_SOURCE_DIR}/src \${PROJECT_SOURCE_DIR}/lib/unity \${PROJECT_SOURCE_DIR}/lib/llog/src \${PROJECT_SOURCE_DIR}/lib/libfault \${PROJECT_SOURCE_DIR}/lib/libran/src)"             >> test/CMakeLists.txt

echo "CMAKE_MINIMUM_REQUIRED(VERSION 2.8)"                                      >> tools/CMakeLists.txt
echo "SET(CMAKE_VERBOSE_MAKEFILE on)"                                           >> tools/CMakeLists.txt
echo "ADD_SUBDIRECTORY(run)"                                                    >> tools/CMakeLists.txt
echo "ADD_SUBDIRECTORY(pvtrace)"                                                >> tools/CMakeLists.txt

mkdir -p lib/
# Tools for Generate Call Graph
echo "CMAKE_MINIMUM_REQUIRED(VERSION 2.8)"                                        > tools/CMakeLists.txt
echo "SET(CMAKE_VERBOSE_MAKEFILE on)"                                             >> tools/CMakeLists.txt
echo "ADD_SUBDIRECTORY(run)"                                                      >> tools/CMakeLists.txt
echo "ADD_SUBDIRECTORY(pvtrace)"                                                  >> tools/CMakeLists.txt

# instrument.c should not be build in pvtrace daemon
git clone https://github.com/lamproae/pvtrace.git                                 tools/pvtrace 
cp -a tools/pvtrace/instrument.c test/ && cp -a tools/pvtrace/instrument.c app/ && rm -rf tools/pvtrace/instrument.c && rm -rf tools/pvtrace/Makefile
echo "AUX_SOURCE_DIRECTORY(\${CMAKE_CURRENT_SOURCE_DIR} APP_SRC)"                 > tools/pvtrace/CMakeLists.txt
echo "ADD_EXECUTABLE(pvtrace \${APP_SRC})"                                        >> tools/pvtrace/CMakeLists.txt
echo "INCLUDE_DIRECTORIES(\${CMAKE_CURRENT_SOURCE_DIR})"                          >> tools/pvtrace/CMakeLists.txt
echo "SET(EXECUTABLE_OUTPUT_PATH \${PROJECT_BINARY_DIR}/bin)"                     >> tools/pvtrace/CMakeLists.txt
echo "SET(CMAKE_C_FLAGS "-Wall")"                                                 >> tools/pvtrace/CMakeLists.txt

mkdir -p tools/run
echo "AUX_SOURCE_DIRECTORY(\${CMAKE_CURRENT_SOURCE_DIR} APP_SRC)"                 > tools/run/CMakeLists.txt
echo "ADD_EXECUTABLE(run \${APP_SRC})"                                            >> tools/run/CMakeLists.txt
echo "INCLUDE_DIRECTORIES(\${CMAKE_CURRENT_SOURCE_DIR})"                          >> tools/run/CMakeLists.txt
echo "SET(EXECUTABLE_OUTPUT_PATH \${PROJECT_BINARY_DIR}/bin)"                     >> tools/run/CMakeLists.txt
echo "SET(CMAKE_C_FLAGS "-Wall")"                                                 >> tools/run/CMakeLists.txt

echo "#include <stdlib.h>"                                                        > tools/run/run.c
echo ""                                                                           >> tools/run/run.c
echo "void show_usuage(void)"                                                     >> tools/run/run.c
echo "{"                                                                          >> tools/run/run.c
echo "    printf(\"Usuage: run APP\\\n\");"                                       >> tools/run/run.c
echo "}"                                                                          >> tools/run/run.c
echo ""                                                                           >> tools/run/run.c
echo "int main(int argc, char** argv) "                                           >> tools/run/run.c
echo "{"                                                                          >> tools/run/run.c
echo "    char command[256] = {0};"                                               >> tools/run/run.c
echo ""                                                                           >> tools/run/run.c
echo "    if (argc != 2) {"                                                       >> tools/run/run.c
echo "        show_usuage();"                                                     >> tools/run/run.c
echo "        exit(0);"                                                           >> tools/run/run.c
echo "    }"                                                                      >> tools/run/run.c
echo ""                                                                           >> tools/run/run.c
echo "    sprintf(command, \"./%s\", argv[1]);"                                   >> tools/run/run.c
echo "    system(command);"                                                       >> tools/run/run.c
echo "    sprintf(command, \"./pvtrace %s\", argv[1]);"                           >> tools/run/run.c
echo "    system(command);"                                                       >> tools/run/run.c
echo "    sprintf(command, \"dot -Tpng graph.dot -o graph.png\");"                >> tools/run/run.c
echo "    system(command);"                                                       >> tools/run/run.c
echo "    sprintf(command, \"valgrind --tool=memcheck --leak-check=full --log-file=memcheck.log ./%s\", argv[1]);">> tools/run/run.c
echo "    system(command);"                                                       >> tools/run/run.c
echo "    sprintf(command, \"valgrind --tool=callgrind ./%s > callgrind.log 2>&1\", argv[1]);">> tools/run/run.c
echo "    system(command);"                                                       >> tools/run/run.c
echo "    sprintf(command, \"gprof ./%s > gprof.log 2>&1\", argv[1]);"            >> tools/run/run.c
echo "    system(command);"                                                       >> tools/run/run.c
echo "    sprintf(command, \"mtrace ./%s %s> memleak_detect.log 2>&1\", argv[1], \"leak_detector.log\");" >> tools/run/run.c
echo ""                                                                           >> tools/run/run.c
echo "    /*  Need to Add gcov support for code coverage test in future */"       >> tools/run/run.c
echo "}"                                                                          >> tools/run/run.c
echo ""                                                                           >> tools/run/run.c

# Unity C unit test framwrok.
git clone https://github.com/lamproae/unity.git                                 lib/unity
# For mysql operation
# git clone https://github.com/lamproae/libmysql.git                              lib/mysql
# git clone https://github.com/lamproae/c-algorithms.git                          lib/libds
# git clone https://github.com/lamproae/Collections-C.git                         lib/collection
# Onion can be enabled if necessary
# git clone https://github.com/lamproae/onion.git                                 lib/onion
# libnet
# git clone https://github.com/lamproae/libnet.git                                lib/net
git clone https://github.com/lamproae/llog.git                                  lib/llog
git clone https://github.com/lamproae/libfault.git                              lib/libfault
git clone https://github.com/lamproae/memwatch.git                              lib/memwatch
git clone https://github.com/lamproae/etrace.git                                lib/etrace

#Common library for practise
# git clone https://github.com/lamproae/liblw.git                                 lib/liblw

# Generic Container library

# git clone https://github.com/lamproae/klib.git                                  lib/klib.git
# git clone https://github.com/lamproae/qlibc.git                                 lib/qlibc.git

# git clone https://github.com/lamproae/gperftools.git                            lib/gperftools
# git clone https://github.com/lamproae/libev.git                                 lib/libev
# git clone https://github.com/lamproae/Libevent.git                              lib/Libevent
# git clone https://github.com/lamproae/unpv13e.git                               lib/unpv13e
# git clone https://github.com/lamproae/apue.3e.git                               lib/apue

# A random number generte lib from ovs
git clone https://github.com/lamproae/libran.git                                lib/libran
# Unix network programming and Advanced programming in unix environment
# git clone https://github.com/lamproae/libunp.git                                lib/libunp
# git clone https://github.com/lamproae/libapue.git                               lib/libapue
# git clone https://github.com/lamproae/jemalloc.git                              lib/jemalloc
# git clone https://github.com/lamproae/jemalloc-cmake.git                        lib/jemalloc-cmake


# git clone https://github.com/lamproae/pvtrace.git                               tools/pvtrace
# git clone https://github.com/lamproae/gprof2dot.git                             tools/gprof2dot
# git clone https://github.com/lamproae/xdot.py.git                               tools/xdot
# git clone https://github.com/lamproae/codeviz.git                               tools/codeviz
# git clone https://github.com/lamproae/graphviz.git                              tools/graphviz

# A Modern and Easy-to-use crypto library
# git clone https://github.com/lamproae/libsodium.git                             lib/libsodium

# A Portable UI library for C                                                     
# git clone https://github.com/lamproae/libui.git                                 lib/libui

# Go-style concurrency in C
# git clone https://github.com/lamproae/libmill.git                                 lib/libmill

# Structured concurrency in C
# git clone https://github.com/lamproae/libdill.git                               lib/libdill

# An eventing framework for building high performance and high scalability systems in C
# git clone https://github.com/lamproae/libphenom.git                             lib/libphenom

# An library for parsin/normalizing street addresses around the world
# git clone https://github.com/lamproae/libpostal.git                             lib/libpostal

# An ultra low-latency trading connectivity library for C and C++
# git clone https://github.com/lamproae/libtrading.git                            lib/libtrading

# An Basic libraries all written in c. Include network, event, config, hash, ipc, rpc, mem and so on
# git clone https://github.com/lamproae/libraries.git                             lib/libraries

# An Mail Framwork for C language
# git clone https://github.com/lamproae/libetpan.git                              lib/libetpan

# An Adaptive radix tree implementation
# git clone https://github.com/lamproae/libart.git                                lib/libart.git

# A Common C libarary
# git clone https://github.com/lamproae/libfastcommon.git                         lib/libfastcommon

# A multiplatfrom utility libarary written in C 
# git clone https://github.com/koanlogic/libu.git                                 lib/libu

# A simple, easy embeddable cross-platform C library
# git clone https://github.com/lamproae/libcork.git                               lib/libcork

# C library for thread safe hash-table linked list...
# git clone https://github.com/lamproae/libhl.git                                 lib/libhl

# A JSON parser and printer library in C
# git clone https://github.com/lamproae/libjson.git                               lib/libjson

# Cross-platform library supports iOS/Android development in C language
# git clone https://github.com/lamproae/libphone.git                              lib/libphone

# Cross-platform public domain foundation library in C 
# git clone https://github.com/lamproae/foundation_lib.git                        lib/libfoundation

# C library for easy WebSockets server
# git clone https://github.com/lamproae/libwebsock.git                            lib/libwebsock

# A c unit test framewor
# git clone https://github.com/libcheck/check.git                                 lib/check

# A simple http-parser
# git clone https://github.com/nodejs/http-parser.git                             lib/http-parser

# A command line tool and library to transfer data in URL format
# git clone https://github.com/lamproae/curl.git                                  lib/curl
