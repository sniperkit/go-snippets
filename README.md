# cppinit

`cppinit` is a CLI utility written in Go that generates empty C++ projects for
my particular configuration and tastes.

It requires the following:
* [VSCode](https://code.visualstudio.com/) (and the plugin CMake Tools by vector-of-bool)
* [CMake](https://cmake.org/) (v3.11+)
* [Boost](https://www.boost.org/) (v1.67)
* [Catch](https://github.com/catchorg/Catch2) (v2)
* [Ninja](https://ninja-build.org/)

Download Boost via `brew`, `pacman` or `vcpkg`. Because Ubuntu is usually so behind, simply download and build the source.

For `cppinit` itself:

1. Download and install the Go toolkit (1.10 was used at time of writing)
2. `go get github.com/LeonineKing1199/cppinit`
3. `cd $GOPATH/src/github.com/LeonineKing1199/cppinit`
3. `go build`
4. Take the compiled binary and add it to a dedicated folder on your system's PATH

To use:

Build an empty project with:
```
mkdir my-new-project
cd my-new-project
./cppinit -name=my-new-project

=> creates
C:.
│   .gitignore
│   CMakeLists.txt
│   my-new-project.cmake
├───.vscode
│       settings.json
│
├───cmake.modules
│       ParseAndAddCatchTests.cmake
│
├───include
│   └───my-new-project
└───test
    │   main.cpp
    │
    └───include
        └───my-new-project
            └───test
```

To purge the current directory,
```
./cppinit -clean
```

**Note**: Inovking `cppinit -clean` will delete anything in the `include/` and `test/` directories which means one can lose quite a bit of work. Always make sure `cppinit -clean` is invoked with caution. And always have backups of anything you care about.

## For C++ Beginners

It's a goal of this project to make using C++ for newcomers simple and
straight-forward. `cppinit`'s main strength is that it gives users a project
with some useful include paths already configured so they can begin coding
without mucking about.

### Boost

Boost installation steps vary wildly from system to system but in most cases,
one can simply download the compressed source, decompress it and then use the
documentation to build Boost with `b2`, making sure that one uses the `prefix`
option to create an easy-to-use Boost root folder. Otherwise, the include and
library directories will both need to be specified.

### Creating Your CMake Toolchain File

One may notice in the above example that a `my-new-project.cmake` file is generated
by `cppinit`. This is intended to be an empty CMake toolchain file. Toolchain files
are easy ways to feed system-specific variables into CMake without polluting the root
`CMakeLists.txt` file which should only be used for system-agnostic requirements of
the project.

Sample toolchain file:
```
# Sample file for Windows users

set(BOOST_ROOT "/Users/cmaza/source/boosts/boost_1_67_0_b1")
set(Boost_USE_STATIC_LIBS ON)

# directory that contains our catch.hpp file
include_directories(
  "/vcpkg/installed/x64-windows/include/catch"
)
```

Here one can see that we define the `BOOST_ROOT` to be the directory where we told
`b2` to install Boost. We also wish to use only the static libraries and then
we help CMake find our installation of Catch which was done via `vcpkg`.

Building and testing your project:
```
cd my-new-project
mkdir build_Debug
cd build_Debug
cmake \
-DCMAKE_TOOLCHAIN_FILE=../my-new-project.cmake \
-DCMAKE_BUILD_TYPE=Debug \
-G Ninja ..
cmake --build .
./my-new-project_tests
```

`cppinit` ultimately gives the programmer a CMakeLists.txt file that they can
now edit and maintain themselves.

The core project is pulled into a "core" lib that is then linked to by the testing binary.
It defaults to an INTERFACE library. This needs to be changed if the
core library itself becomes non-header-only.

### Using the empty project by writing your first test

As a good example of what to do with `cppinit`, considering building the sample
demo files:
```
-------------------------------------------------
test/include/my-new-project/test/integral_add.hpp
-------------------------------------------------

#include <numeric>
#include <type_traits>

// we choose to write a constrained templated function
// that will only add integral types
//
template <
  typename T,
  typename = std::enable_if_t<std::is_integral<T>::value>
>
auto integral_add(T const& a, T const& b) -> T
{
  return a + b;
}

-------------------------------------------------
test/lifting_add_test.cpp
-------------------------------------------------

#include <vector>
#include <numeric>
#include <boost/hof/lift.hpp>
#include "my-new-project/test/integral_add.hpp"

#include <catch.hpp>

// we want to use our constrained adding function
// but we want to make sure we can also use it in
// higher-order functions
//
TEST_CASE("do u even lift?")
{
  // 0 + 1 + 2 + 3 + 4 + 5 => 3 + 3 + 4 + 5 => 6 + 9 => 15
  //
  auto const input = std::vector<int>{0, 1, 2, 3, 4, 5};
  auto const sum   = std::accumulate(
    input.begin(), input.end(),
    0,
    BOOST_HOF_LIFT(integral_add));

  REQUIRE(sum == 15);
}
```

Then in the CMakeLists.txt file, simply do:
```
add_executable(
  my-new-project_tests

  ${CMAKE_CURRENT_SOURCE_DIR}/test/main.cpp

  # all we're doing is simply adding this line to the existing file
  ${CMAKE_CURRENT_SOURCE_DIR}/test/lifting_add_test.cpp
)
```

Now, assuming we've run CMake once before with the proper toolchain file,
all we should need to do now is invoke:
```
cmake --build .
ctest .
```
from the command line and this will both build all the tests and then run
them.

Expected output:
```
Test project .../build_debug
    Start 1: my-new-project_tests:do u even lift?
1/1 Test #1: my-new-project_tests:do u even lift? .....   Passed    0.03 sec

100% tests passed, 0 tests failed out of 1

Label Time Summary:
my-new-project_tests    =   0.03 sec*proc (1 test)

Total Test time (real) =   0.07 sec

```

For Windows/msvc users, you'll have to do all this from a Developer Command Prompt.
Use `vcvarsall` to configure your specific environment.

In the above example we chose to put the file under our testing directory. This isn't
strictly necessary. `cppinit` sets default include paths that let you include anything
from Boost, `./include/` and `./test/include/`. We also add a default directory
`test/include/my-new-project/test/` so that one can also write a coupled set of tests with
their library itself. This more done out of convenience than anything else.