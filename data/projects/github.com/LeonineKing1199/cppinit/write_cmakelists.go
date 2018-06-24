package main

import "os"

func writeCMakeLists(proj string, file *os.File) (err error) {
	_, err = file.WriteString(`cmake_minimum_required(VERSION 3.11)

project(` + proj + `)

find_path(CATCH_INCLUDE_DIR "catch.hpp")

add_library(Catch INTERFACE)
target_include_directories(Catch INTERFACE ${CATCH_INCLUDE_DIR})

find_package(
  Boost 1.67
  REQUIRED
)

add_library(
  ` + proj + `_core

  INTERFACE
)

if (MSVC)
  # Win10
  target_compile_definitions(` + proj + `_core INTERFACE _WIN32_WINNT=0x0A00)
endif()

target_compile_definitions(
  ` + proj + `_core

  INTERFACE
  BOOST_COROUTINES_NO_DEPRECATION_WARNING=1
  BOOST_CONFIG_SUPPRESS_OUTDATED_MESSAGE=1
)

target_compile_features(` + proj + `_core INTERFACE cxx_std_14)

target_include_directories(
  ` + proj + `_core

  INTERFACE
  ${CMAKE_CURRENT_SOURCE_DIR}/include
)

target_link_libraries(
  ` + proj + `_core

  INTERFACE
  Boost::boost
)

# keep this around, just in case
if (MSVC)
  target_link_libraries(` + proj + `_core INTERFACE Boost::disable_autolinking)
endif()

add_library(test_utils INTERFACE)
target_include_directories(test_utils INTERFACE ${CMAKE_CURRENT_SOURCE_DIR}/test/include)

add_executable(
  ` + proj + `_tests

  ${CMAKE_CURRENT_SOURCE_DIR}/test/main.cpp
)

target_link_libraries(
  ` + proj + `_tests

  PRIVATE
  ` + proj + `_core
  test_utils
  Catch
)

set(CMAKE_MODULE_PATH ${CMAKE_MODULE_PATH} "${CMAKE_SOURCE_DIR}/cmake.modules/")
enable_testing()
include(ParseAndAddCatchTests)
ParseAndAddCatchTests(` + proj + `_tests)
`)

	return
}
