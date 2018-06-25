#!/bin/bash
set -e

WORK_DIR=`dirname $`

pushd $WORK_DIR &> /dev/null

function help ()
{
    echo "build_samples.sh"
    echo " "
    echo "build_samples.sh [options]"
    echo " "
    echo "options:"
    echo "-h			Show brief help (this)"
    echo "-l			List CMake generators"
    echo "-g <generator name>	Choose CMake generator (defaults to Makefiles)"
    echo "			  - example: ./build_samples.sh -g Xcode # Build using Xcode-based project"
    echo "-d			Build with debug symbols"
    echo
}

function list-generators ()
{
    cmake --help | awk '$0 == "Generators" {i=1};i'
}

BUILD_TYPE=RelWithDebInfo

while getopts 'dhlg:' OPT; do
    case "${OPT}" in
        g) GENERATOR="${OPTARG}" ;;
        l) list-generators
           exit 0 ;;
        h) help ; list-generators
           exit 0 ;;
        d) BUILD_TYPE=Debug
           echo "Building with debug symbols" ;;
        *) help
           exit 255 ;;
    esac
done

if [ ! -d build ]; then
    mkdir build
fi

cd build

if [ ! -z "${GENERATOR}" ]; then
    GENERATOR_ARG=-G"${GENERATOR}"
fi

cmake ${GENERATOR_ARG} -DCMAKE_BUILD_TYPE=${BUILD_TYPE} .. && cmake --build . --target clean
cmake --build . --config "${BUILD_TYPE}"

echo "Samples built into $WORK_DIR/build/bin"

popd &> /dev/null
