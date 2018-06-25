set (_projname "mock_sensor")

set(${_projname}_HEADERS
  mock_base_stream.hpp
  mock_color_generator.hpp
  mock_color_sensor.hpp
  mock_depth_generator.hpp
  mock_depth_sensor.hpp
  mock_device.hpp
  mock_device_color_stream.hpp
  mock_device_depth_stream.hpp
  mock_device_ir_stream.hpp
  mock_device_stream.hpp
  mock_device_streamset.hpp
  mock_flag_set.hpp
  mock_frame_generator.hpp
  mock_infrared_generator.hpp
  mock_sensor_plugin.hpp
  mock_stream.hpp
  )

set(${_projname}_SOURCES
  mock_base_stream.cpp
  mock_color_generator.cpp
  mock_color_sensor.cpp
  mock_depth_generator.cpp
  mock_depth_sensor.cpp
  mock_device.cpp
  mock_device_color_stream.cpp
  mock_device_depth_stream.cpp
  mock_device_ir_stream.cpp
  mock_device_streamset.cpp
  mock_flag_set.cpp
  mock_frame_generator.cpp
  mock_infrared_generator.cpp
  mock_sensor_plugin.cpp
  )

find_package(SFML COMPONENTS graphics window system)

add_definitions(-DASTRA_BUILD)

include_directories(${_projname} ${SFML_INCLUDE_DIR} ../device_sdk/include/)

add_library(${_projname} SHARED ${${_projname}_SOURCES} ${${_projname}_HEADERS})

set_target_properties(${_projname} PROPERTIES FOLDER "plugins")

target_link_libraries(${_projname} astra_core_api astra-device-sdk ${SFML_LIBRARIES})