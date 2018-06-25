set(MODULE_NAME "astra-device-sdk")

set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)

set(${MODULE_NAME}_HEADERS
  include/astra_device.hpp
  include/astra_device_info.hpp
  include/astra_device_listener.hpp
  include/astra_device_state.hpp
  include/astra_device_status.hpp
  include/astra_sensor.hpp
  include/astra_sensor_info.hpp
  include/astra_sensor_listener.hpp
  include/astra_sensor_state.hpp
  include/astra_sensor_stream.hpp
  include/astra_stream_listener.hpp
  include/astra_stream_state.hpp
  )

set(${MODULE_NAME}_SOURCES
  src/astra_device.cpp
  src/astra_device_info.cpp
  src/astra_device_state.cpp
  src/astra_device_status.cpp
  src/astra_sensor.cpp
  src/astra_sensor_info.cpp
  src/astra_sensor_state.cpp
  src/astra_sensor_stream.cpp
  src/astra_stream_mode.cpp
  src/astra_stream_state.cpp
  )

include_directories(include/)

add_library(${MODULE_NAME} STATIC ${${MODULE_NAME}_SOURCES} ${${MODULE_NAME}_HEADERS})
