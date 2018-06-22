---
layout: post
title: "CMake: Some notes as I discover how to deal with it."
tagline: "CMake... oh my!"
category : notes
published: true
tags : [C++, CMake, howto]
---

_Some notes about how to setup CMake projects._

* **2017-12 :** I'm rediscovering all this, again; please consider the
  information here with skepticism.
* __TODO:__ Have a look at [GYP](https://gyp.gsrc.io) &ndash; “meta-build
  system” used for building Chrome/Chromium.
* [github.com/ttroy50/cmake-examples](https://github.com/ttroy50/cmake-examples)
  : good tutorial.

## Running CMake

The simplest usage would be something like :

```bash
$ mkdir build/ &&
  cd build/ &&
  cmake ../ &&
  cmake --build build/ --target simplest-example
```

Big one-liner where we also wipe out the `build/` sub-directory :

```bash
$ ( [ -d build ] && rm -r build || true ) &&
( mkdir build && cd build &&
  CC=clang CXX=clang++ \
    cmake -G Ninja \
      -DCMAKE_BUILD_TYPE=RelWithDebInfo \
      -DCMAKE_INSTALL_PREFIX=~/local ../
) &&
     cmake --build build --target simplest-example
```

## CMake language, syntax, commands

* [list()](https://cmake.org/cmake/help/latest/command/list.html)

Test empty string variable :

```cmake
if (NOT LIBUNW_LDFLAGS)
  ...
endif()
```

## CMake variables

### Those few ones you'll be using all the time

* [cmake-properties](https://cmake.org/cmake/help/latest/manual/cmake-properties.7.html)

A list of variables and properties we often have to use :

  <table class="phi phi-custom">
<thead>
  <tr>
    <th>CMake variable</th>
    <th>Description</th>
  </tr>
</thead>
<tbody>
  <tr>
    <th>CMAKE_SOURCE_DIR</th>
    <td>The source directory (i.e. the directory of the top level CMakeLists.txt file)</td>
  </tr>
  <tr>
    <th>CMAKE_BINARY_DIR</th>
    <td>The (top level) build directory</td>
  </tr>
  <tr>
    <th>CMAKE_CURRENT_SOURCE_DIR</th>
    <td>The current source directory, i.e. location of the currently processed CMakeLists.txt file (top level or included via ADD_SUBDIRECTORY)</td>
  </tr>
  <tr>
    <th>CMAKE_CURRENT_BINARY_DIR</th>
    <td>The build (sub)directory corresponding to CMAKE_CURRENT_SOURCE_DIR</td>
  </tr>
  <tr> <td colspan="2" class="separator"></td> </tr>
  <tr>
    <th>RUNTIME_OUTPUT_DIRECTORY</th>
    <td></td>
  </tr>
  <tr>
    <th>LIBRARY_OUTPUT_DIRECTORY</th>
    <td></td>
  </tr>
  <tr>
    <th>ARCHIVE_OUTPUT_DIRECTORY</th>
    <td></td>
  </tr>
  <tr> <td colspan="2" class="separator"></td> </tr>
  <tr>
    <th colspan="2">
      Target properties<br>
      <tt>get_property(_outvar TARGET &lt;target&gt; PROPERTY &lt;property&gt;)</tt>
    </th>
  </tr>
  <tr>
    <th>property</th>
    <th>&nbsp;</th>
  </tr>
  <tr>
    <th>TYPE</th>
    <td>
      <code>get_property(_outvar TARGET &lt;some-target&gt; PROPERTY &lt;TYPE&gt;</code><br>
      <tt>_outvar</tt> would then be one of:
      <tt>STATIC_LIBRARY, MODULE_LIBRARY, SHARED_LIBRARY, EXECUTABLE</tt>
    </td>
  </tr>
  <!--
  <tr>
    <th></th>
    <td></td>
  </tr>
  -->
</tbody>
  </table>


### Dump all CMake variables

```cmake
function(dude_please_dump_cmake_variables)

  get_cmake_property(_variableNames VARIABLES)

  foreach (_variableName ${_variableNames})
    message(STATUS "${_variableName} = ${${_variableName}}")
  endforeach()

endfunction()
```

Found that at [SO](http://stackoverflow.com/questions/9298278/cmake-print-out-all-accessible-variables-in-a-script).

#### Dump all variables that match a regular expression

```cmake
function(dump_cmake_variables_matching regex)
  get_cmake_property(varnames VARIABLES)
  foreach (vname ${varnames})
    if ("${vname}" MATCHES "${regex}")
      message(STATUS "${vname} = ${${vname}}")
    endif()
  endforeach()
endfunction()
```

##### Example: finding out variables of interest produced by `pkg_check_modules()`.

```cmake
## CMakeLists.txt

# Postgresql libpqxx
find_package( PkgConfig )
pkg_check_modules( LIBPQXX libpqxx )
if (NOT LIBPQXX_FOUND)
  message( WARNING "Dude: Couldn't Postgresql libpqxx." )
endif()

dump_cmake_variables_matching(".*PQXX.*")
```

Would dump those variables that were set by `pkg_check_modules(...)`, of which :

```cmake
LIBPQXX_CFLAGS =
LIBPQXX_FOUND = 1
LIBPQXX_INCLUDE_DIRS =
LIBPQXX_LDFLAGS = -lpqxx
LIBPQXX_LIBDIR = /usr/lib
LIBPQXX_LIBRARIES = pqxx
```

And configuring a CMake target :

```cmake
if (NOT LIBPQXX_FOUND)
  message( WARNING "Dude: Couldn't find Postgresql libpqxx." )
else()
  message( STATUS "Dude: Found Postgresql libpqxx." )
  target_compile_options(     pimplio PUBLIC ${LIBPQXX_CFLAGS} )
  target_include_directories( pimplio PUBLIC ${LIBPQXX_INCLUDE_DIRS} )
  target_link_libraries(      pimplio        ${LIBPQXX_LIBRARIES} )
endif()
```

## Environment variables

Test if environment variable is defined :

```cmake
if (NOT DEFINED ENV{BOOSTROOT})
  message(STATUS "FYI: Environment variable $BOOTROOT is set to: ${ENV{BOOSTROOT}}")
endif()
```

### add\_executable() & add\_library()

```cmake
add_executable( brand-new-program
  src/main.cpp
  src/impl1.cpp
  src/impl2.cpp
  )
```

### Include current directory

```cmake
set(CMAKE_INCLUDE_CURRENT_DIR ON)
```

#### Changing the output directory to for ex. 'bin/' or 'lib/'

##### Per specific target :

```cmake
set_property(TARGET brand-new-program
  PROPERTY RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin )

set_property(TARGET util
  PROPERTY LIBRARY_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib )

set(TARGET static-stuff
  PROPERTY ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib )
```

##### Globally (for the whole CMake project)

[via](http://bytefreaks.net/programming-2/manually-set-the-cmake-output-folder)

```cmake
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin)
```

##### Searching for targets in a subdirectory, changing output directory :

__A/__ Fetch the __list of sub-directories__.  The property `SUBDIRECTORIES` does just that, though it returns absolute path names (which is ok and concise (unless for printing out on screen)) :

```cmake
get_property(subdirs_list
  DIRECTORY "${CMAKE_CURRENT_SOURCE_DIR}"
  PROPERTY SUBDIRECTORIES)
```

__A”/__Another possiblity is to manual list all files :

```cmake
file(GLOB files_list
  RELATIVE ${CMAKE_CURRENT_SOURCE_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}/*)
```

__A””/__ And filter out anything that is not a directory :

```cmake
set(subdirs_list "")

foreach(subdir ${files_list})
  if (IS_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/${subdir})
    list(APPEND subdirs_list "${subdir}")
  endif()
endforeach()
```

__B/__ Finally we can obtain targets by means of the `BUILDSYSTEM_TARGETS` property :


```cmake
foreach(subdir ${subdirs_list})
  get_property(_targets_list
    DIRECTORY "${subdir}"
    PROPERTY BUILDSYSTEM_TARGETS )

  if (NOT _targets_list)
    message(STATUS "FYI: ${subdir}: Provided no build target.")

  else()
```

__B”/__ And loop over these, doing our cryptic stuff : Finding out the `TYPE` property for each target, and setting the appropriate value for prop. `xxx_OUTPUT_DIRECTORY` :

```cmake
    foreach(_target ${_targets_list})
      get_property(_target_type
        TARGET ${_target}
        PROPERTY TYPE)

      if (${_target_type} MATCHES "_LIBRARY$")
        set_property(
          TARGET ${_target}
          PROPERTY LIBRARY_OUTPUT_DIRECTORY
            ${CMAKE_BINARY_DIR}/lib )

      elseif (${_target_type} STREQUAL "EXECUTABLE")
        set_property(
          TARGET ${_target}
          PROPERTY RUNTIME_OUTPUT_DIRECTORY
            ${CMAKE_BINARY_DIR}/bin )

      else()
        message(STATUS "FYI: ${subdir}: Target '${_target}' is-a '${_target_type}'.")
      endif()
    endforeach()

  endif()
endforeach()
```

### Test if a file or directory exists

```cmake
if (NOT EXISTS ...)
endif()
```

```cmake
if (NOT IS_DIRECTORY ${SOCI_SOURCES_DIR}/src)
  message( WARNING "Dude: Could not find SOCI sources at ${SOCI_SOURCES_DIR}" )
endif()
```

### Find out the name of the "current directory"

TODO: some blah blah

```cmake
get_filename_component(ThisModuleName "${CMAKE_CURRENT_SOURCE_DIR}" NAME)
set(ThisModuleName "dude-${ThisModuleName}")
message(STATUS "Dude! I'm ${ThisModuleName}.")
```

```cmake
add_library(${ThisModuleName} library-entry-point.cpp other-stuff.cpp )
target_compile_options( ${ThisModuleName} PUBLIC ${LIBxxxx_CFLAGS} )
target_include_directories( ${ThisModuleName} PUBLIC ${LIBxxxx_INCLUDE_DIRS} )
```

```cmake
add_executable(${ThisModuleName} main.cpp)
```

### How to add a new sub-directory ?

__In short,__ just create the subdirectory and :

* drop an empty `CMakeLists.txt` file in it ;
* __either__ append `add_subdirectory( src/new-subdir )` to the top-level project's `CMakeLists.txt` ;
* __or__, if applicable, edit the "_parent"_ `CMakeLists.txt` file one directory up, for ex. `src/`.
* __Done:__ that's all it takes.
* Then you may proceed with the usual directives :
- __add_executable()__
- etc...

Copy-paste "one-liner" (remove the trailing blankees `_`):

```bash
$ cat <<EOF  >>  CMakeLists.txt__

add_executable( exe-bin-name
  exe-bin-main.cpp
  exe-bin-util.cpp
  )

target_compile_options( exe-bin-name
  PUBLIC
    -std=c++1z
    -Wall -pedantic
    -Wmissing-field-initializers -Wnon-virtual-dtor
    -Wdelete-non-virtual-dtor    -Woverloaded-virtual
  )

target_include_directories( exe-bin-name
  PUBLIC
    ${XCB_INCLUDE_DIRS}
    ${XCB_AUX_INCLUDE_DIRS}
    ${CAIRO_XCB_INCLUDE_DIRS}
  )

target_link_libraries( exe-bin-name
    ${Boost_LIBRARIES}
    ${XCB_LIBRARIES}
    ${XCB_AUX_LIBRARIES}
    ${CAIRO_XCB_LIBRARIES}
  )

option( DUDE_TURN_ME_ON
  "Dude, please turn this option ON for Knuth's sake (!)" OFF )

if (DUDE_TURN_ME_ON)
  message( WARNING "Dude: You should not have enabled this feature, really." )
  add_executable( magic-bin magic-bin.cpp )
endif()

if (NOT DUDE_TURN_ME_ON)
  add_subdirectory( tests )
fi

# EOF message.
message(STATUS "Dude: This is a communication from ${CMAKE_CURRENT_LIST_FILE}")
EOF
```

### List sub-directories (loop over these and do stuff)

__Or:__ how to list "direct" sub-directories for automated addition to CMake project ([via](https://stackoverflow.com/a/7788165/643087)).

Basically we have `file(GLOB files_list ... RELATIVE ... )` which output value (here `files_list`) we'll filter out files from dirs
([file()](https://cmake.org/cmake/help/latest/command/file.html)])

```cmake
file(GLOB files_list
  RELATIVE ${CMAKE_CURRENT_SOURCE_DIR}
           ${CMAKE_CURRENT_SOURCE_DIR}/*)
```

We may want to collect sub-directory names in a variable `subdirs_list` :

```cmake
set(subdirs_list "")
```

Loop over the `files_list`, acting upon directories which contain a `CMakeLists.txt` file :

```cmake
foreach(subdir ${files_list})
  if(EXISTS ${CMAKE_CURRENT_SOURCE_DIR}/${subdir}/CMakeLists.txt)

    add_subdirectory( ${subdir} )
    message(STATUS "Dude: Added experiment ${subdir}")
    list(APPEND subdirs_list "${subdir}")

  elseif(IS_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/${subdir})

    message(STATUS "Dude: Skipping ${subdir}/ as it contains no CMakeLists.txt (FYI)")

  endif()
endforeach()
```

The `subdirs_list` was indeed set :

```cmake
message(STATUS "Dude: Experiments list: ${subdirs_list}")
```

Futile :

```cmake
unset(files_list)
unset(subdirs_list)
```


#### Adding a CMake-based third-party code as a sub-project

Say we'll manage a few of our project dependencies by means of Git submodule, for ex. :

```bash
git submodule add git@... vendor/sub-project
```

Then if the sub-project is "properly" configured, a simple
`add_subdirectory( vendor/sub-project )` may be enough.

__But (!)__ ... stating dependencies ? (TODO)


### find\_library() & find\_file() : Manual library search & setup

* `CMAKE_INCLUDE_PATH` is searched by [find\_file()](https://cmake.org/cmake/help/latest/command/find_file.html) ;
* `CMAKE_LIBRARY_PATH` is searched by [find\_library()](https://cmake.org/cmake/help/latest/command/find_library.html) ;

Given we install stuff into a `./local/` sub-directory (for ex. `./configure --prefix=$PWD/local` or `cmake -DCMAKE_INSTALL_PREFIX=$PWD/local`),
then we end up with an [FHS](https://en.wikipedia.org/wiki/Filesystem_Hierarchy_Standard)-like sub-directory structure where we are interested
in the `local/include/` headers directory and `local/lib/` for static/shared
libraries (and maybe also `lib32/` and/or `lib64/`).

```cmake
# CMakeLists.txt
if ( IS_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}/local/include )
  message( STATUS "Dude: Found ./local/ FHS-like sub-directory." )

  list(APPEND CMAKE_INCLUDE_PATH "${CMAKE_SOURCE_DIR}/local/include")
  list(APPEND CMAKE_LIBRARY_PATH "${CMAKE_SOURCE_DIR}/local/lib")

  message( STATUS "Updated CMAKE_INCLUDE_PATH = ${CMAKE_INCLUDE_PATH}" )
  message( STATUS "Updated CMAKE_LIBRARY_PATH = ${CMAKE_LIBRARY_PATH}" )
endif()
```

Then we may resort to `find_library()` and `find_file()`
for configuring our stuff :

```cmake
find_library( LibSociCore  libsoci_core.so
  HINTS
    ${CMAKE_CURRENT_SOURCE_DIR}/build/some/other/lib/dir
    ${CMAKE_SOURCE_DIR}/build/lib/wtf
  )
```

And for header files :

```cmake
find_file( LibSociMainHeaderFile soci/soci.h )

# ^ this will effectivelly find a `soci.h` file
#   into `./local/include/soci/soci.h`
```

* The optional __HINTS__ (or __PATHS__) may be followed by a list of additional
directories to search into.
* __Note that__ neither __find_library()__ nor __find_file()__ will traverse the search paths recursively (!) &ndash; which is often quite annoying, since the alternative is typically to resort to __[file()](https://cmake.org/cmake/help/latest/command/file.html)'s__ `GLOB` or `GLOB_RECURSE`.

Checking that something was found is also a bit awkward to say the least : if the file you're looking for could not be found, then __CMake__ will have the return variable (ex. `LibSociCore` have its value set to `"LibSociCore-NOTFOUND"` -\_-

Hence we end up doing tests like `if (LibSociCore MATCHES "-NOTFOUND$")` :

```cmake
if ( LibSociCore MATCHES "-NOTFOUND$" )
  message( ERROR "Dude: Could not find a couple of SOCI's shared libraries.")
  message( ERROR "      ` we won't build some stuff.")
else()
  message( STATUS "Dude: Found `libsoci_core.so` at ${LibSociCore}")
  message( STATUS "Dude: Found `libsoci_mysql.so` at ${LibSociMysql}")

  # Remove the file base-name from the path
  get_filename_component( LibSociInclude "${LibSociMainHeaderFile}" DIRECTORY)

  include_directories(
      ${LibSociInclude}
      ${LibSociInclude}/mysql
    )

  add_executable( some-example-exe some-example-exe.cpp)
endif()
```

### "Importing" / linking against objects (\*.o)

It is possible to define targets that produce `.o` object files, and later on reference these from other targets definitions (e.g. libraries, executables) with `$<TARGET_OBJECTS:foo>`.

(__That is to say:__ your intent is to <u>not</u> have _“yet another small library target just for sharing binary code”_ ; and at the same time be able to speed up build-time by avoiding that several targets depend on the same __.cpp__ implementations).

The following would have `src/util/foo.cpp` be built twice, one for target __buzz__, the other for target __booz__ :

```cmake
add_library(buzz src/buzz_main.cpp src/util/foo.cpp)
add_library(booz src/booz.cpp      src/util/foo.cpp)
```

One possible solution (_out of (?)_) is then to have __foo.cpp__ produce only one __foo.o__ object file during the whole build in this way &ndash; __(1) first__ define the __OBJECT__ targets :

```cmake
add_library( foo  OBJECT src/util/foo.cpp )
add_library( bar  OBJECT src/other/bar.cpp )
```

__(2) Then__ one may have for ex. a library target named __buzz__ (that would produce a `libbuzz.so` shared library) import (depend up on) objects of our target `foo`, by resorting to the obscur idiom `$<TARGET_OBJECTS:foo>` :

```cmake
add_library( buzz
  SHARED
    src/buzz_main.cpp
    src/buzz_extra.cpp
    src/buzz_util.cpp

    $<TARGET_OBJECTS:foo>

    ${LIBXXX_LIBRARIES}
    ...
    ...
  )
```

Likewise for an executable `another-exec` which would depend upon objects from target `bar` :

```cmake
add_executable( another-exec
    buzz
    $<TARGET_OBJECTS:bar>
  )
```

### set\_source\_files\_properties()

Or, for ex. "how does one have specific compiler flags for a given source code file ?"

Say we want a given source `yeah.cpp` that require some bleeding edge feature from C++1z :

```cmake
set_source_files_properties( yeah.cpp
  PROPERTIES
    COMPILE_FLAGS -std=c++1z -Wall -Weverything -pedantic )
```

### EXCLUDE\_FROM\_ALL property :

This :

```cmake
set(EXCLUDE_FROM_ALL ON)
```

Or :

```cmake
set_property(DIRECTORY PROPERTY EXCLUDE_FROM_ALL ON)
```

Will have targets defined within the current subtree be excluded from the
list of targets to built per default.

### Set (append) custom target link flags (TODO)

This won't work (?) :

    set_property(TARGET myDLL APPEND PROPERTY LINK_FLAGS /INCLUDE:\_InitLibrary)

Trying this way :

    get_property(link_flags TARGET myDLL PROPERTY LINK_FLAGS)
    set(link_flags "${link_flags} /INCLUDE:\_InitLibrary")
    set_target_properties(myDLL PROPERTIES LINK_FLAGS ${link_flags})

There's an `APPEND_STRING` possibility for `set_property()` ??

[via](http://cmake.3232098.n2.nabble.com/Appending-to-the-LINK-FLAGS-target-property-td6520731.html)

### Setting the default build type (CMAKE\_BUILD\_TYPE)

The build type/variant is generally left for the user to specify, typically through the command line argument `-DCMAKE_BUILD_TYPE=Release`.  It is possible though to have a different default for this.

Via [cmake.org/Wiki/CMake\_FAQ](https://cmake.org/Wiki/CMake_FAQ#How_can_I_change_the_default_build_mode_and_see_it_reflected_in_the_GUI.3F) :

```cmake
IF(NOT CMAKE_BUILD_TYPE)
  #SET(CMAKE_BUILD_TYPE RelWithDebInfo CACHE STRING
  SET(CMAKE_BUILD_TYPE MinSizeDeb CACHE STRING
    "Choose the type of build, options are: None Debug Release RelWithDebInfo MinSizeRel MinSizeDeb."
      FORCE)
ENDIF(NOT CMAKE_BUILD_TYPE)
```


### CCACHE

Found this solution :
[crascit.com/2016/04/09/using-ccache-with-cmake/](https://crascit.com/2016/04/09/using-ccache-with-cmake/)

```cmake
find_program(CCACHE_PROGRAM ccache)

if (true)
  if (CCACHE_PROGRAM)
    set_property( GLOBAL PROPERTY RULE_LAUNCH_COMPILE "${CCACHE_PROGRAM}" )
    message(STATUS "Dude! Found ccache, so we're setting CMake global property RULE_LAUNCH_COMPILE to '${CCACHE_PROGRAM}'.")
  endif()
endif()
```

### Custom targets for running arbitrary commands

Here's an example using [__add_custom_target()__](https://cmake.org/cmake/help/latest/command/add_custom_target.html) for listing executable files in the build directory :

```cmake
add_custom_target( list-binaries
  ALL
  DEPENDS our-final-program
    COMMAND find ${CMAKE_CURRENT_BINARY_DIR} -type d -name CMakeFiles -prune -o -type f -perm -a+x -ls
  WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
  COMMENT "-- List os executable files (+x) that were built : --"
  VERBATIM USES_TERMINAL
  )
```

The trick here is to have it __depend__ upon some final target (here `our-final-program`) that would end up built last.  Option __ALL__ would have that command be run at each build (else one would have to explicitely ask for it, for ex. `cmake --build build/ --target list-binaries`).

Note that it is possible to have several commands be run for a custom target :

```cmake
COMMAND echo -e "\\e[32;1;7m  Here's your static stuff :  \\e[0m"
COMMAND find ${CMAKE_BINARY_DIR} -type f -name '*.a' -ls
COMMAND echo -e "\\e[2m  PS: have a break away from computer at times!  \\e[0m"
COMMAND echo -e "\\e[32;1;7m  DONE  \\e[0m"
```

* [Colors and formatting (ANSI/VT100 Control sequences) @ misc.flogisoft.com](https://misc.flogisoft.com/bash/tip_colors_and_formatting)

#### Custom targets for having two pre-configured Debug / Release builds

One example usage I've seen for [__add_custom_target()__](https://cmake.org/cmake/help/latest/command/add_custom_target.html) is having two custom `debug` and `release` targets that are “maintainer pre-configured” with some intent of having common defaults by turning some options ON/OFF.

This may be useful at least in situations were one has several of those options that configure fragments of code, and you can't always easily tell which combinations of options would, in the best case, break the build ; or worse: not break the build and introduce some pesky bugs &ndash; hence this idea of having pre-configured, say “maintainer likes” builds, for chasing bugs.

```cmake
# Run with: `cmake --build build --target debug`
add_custom_target( debug
  COMMAND ${CMAKE_COMMAND} -DCMAKE_BUILD_TYPE=Debug ${CMAKE_SOURCE_DIR}
  COMMAND ${CMAKE_COMMAND} --build ${CMAKE_BINARY_DIR} --target all
  COMMENT "Run CMake for a Debug build."
  )
```

```cmake
# Run with: `cmake --build build --target release`
add_custom_target(release
  COMMAND ${CMAKE_COMMAND} -DCMAKE_BUILD_TYPE=Release ${CMAKE_SOURCE_DIR}
  COMMAND ${CMAKE_COMMAND} --build ${CMAKE_BINARY_DIR} --target all
  COMMENT "Run CMake for a Release build."
  )
```

## CMake functions

* [CMake - functions](https://cmake.org/cmake/help/latest/command/function.html)
* [SO: Usage of CMakeParseArguments standard module](https://stackoverflow.com/a/23336987/643087)

## Pointers

* [CMAKE\_INCLUDE\_PATH](https://cmake.org/cmake/help/latest/variable/CMAKE_INCLUDE_PATH.html)
* [CMAKE\_LIBRARY\_PATH](https://cmake.org/cmake/help/latest/variable/CMAKE_LIBRARY_PATH.html#variable:CMAKE_LIBRARY_PATH)
* [if, elseif, else, endif](https://cmake.org/cmake/help/latest/command/if.html)
* [CMake regular expressions (`string()`)](https://cmake.org/cmake/help/latest/command/string.html#regex-specification)
* [get\_filename\_component()](https://cmake.org/cmake/help/latest/command/get_filename_component.html)
* [CMAKE\_BUILD\_TYPE](https://cmake.org/cmake/help/latest/variable/CMAKE_BUILD_TYPE.html)
* [Imported targets](https://cmake.org/cmake/help/latest/manual/cmake-buildsystem.7.html#imported-targets)
* [add\_library](https://cmake.org/cmake/help/latest/command/add_library.html)

* <https://github.com/cginternals/cmake-init>
  – “Template for reliable, cross-platform C++ project setup using cmake.”
* [LLVM's CMake Primer (*must read*)](https://llvm.org/docs/CMakePrimer.html)
* <http://foonathan.net/blog/2016/03/03/cmake-install.html>
* <https://cmake.org/Wiki/CMake/Tutorials/How_to_create_a_ProjectConfig.cmake_file>
* <https://cmake.org/Wiki/CMake:How_To_Write_Platform_Checks>
* <https://cmake.org/Wiki/CMake/Language_Syntax>
* <http://www.kaizou.org/2014/11/typical-cmake-project/>
* <https://www.vtk.org/Wiki/CMake/Examples>
* [2016: CMake Dependencies Done Right](http://floooh.github.io/2016/01/12/cmake-dependency-juggling.html)

### Q&A

* [SO: Difference between IMPORTED and INTERFACE libraries](https://stackoverflow.com/a/36649194/643087)

__EOF__
