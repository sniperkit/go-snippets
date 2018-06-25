cmake_minimum_required(VERSION 3.9)

### Hunter snapshot that will be used ###
include("cmake/HunterGate.cmake")
HunterGate(
    URL "https://github.com/ruslo/hunter/archive/v0.19.216.tar.gz"
    SHA1 "247d3ca881e9c83e6ac419d55cd08d1775ea077a"
    LOCAL # see cmake/Hunter/config.cmake for package customization
)
project(hunter-tensorflow-test)

option(htt_USE_TF_SUBMODULE "Use tensorflow as submodule" ON)

# option(tensorflow_PATH_STATIC_LIB
# option(tensorflow_CUDNN_INCLUDE
# option(tensorflow_PATH_CUDNN_STATIC_LIB
    
option(htt_ENABLE_GPU "Enable the GPU" ON)
option(htt_ENABLE_GRPC_SUPPORT "Enable GRPC" ON)

### Download dependencies ###

hunter_add_package(tensorflow)
find_package(tensorflow CONFIG REQUIRED) # tensorflow::tensorflow

### Targets ###
add_executable(simple simple.cpp)
target_link_libraries(simple PUBLIC tensorflow::tensorflow)

add_subdirectory(src)

return() # todo

### Testing ###
enable_testing()

hunter_add_package(GTest)
find_package(GTest CONFIG REQUIRED) # GTest::main

hunter_add_package(gauze)
find_package(gauze CONFIG REQUIRED)

add_executable(test test.cpp)
target_link_libraries(test PUBLIC tensorflow::tensorflow gauze::gauze GTest::main)

gauze_add_test(
  NAME SimpleTest 
  COMMAND test
  WORKING_DIRECTORY "${CMAKE_CURRENT_LIST_DIR}"
  )
