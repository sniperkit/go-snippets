set (_projname "astra")

set(INCLUDE_PATH "../../include/astra")

set(${_projname}_SOURCES
  ${INCLUDE_PATH}/Vector.hpp
  ${INCLUDE_PATH}/Vector2f.hpp
  ${INCLUDE_PATH}/Vector2i.hpp
  ${INCLUDE_PATH}/Vector3f.hpp
  ${INCLUDE_PATH}/Vector3i.hpp
  ${INCLUDE_PATH}/astra.hpp
  ${INCLUDE_PATH}/capi/astra.h
  ${INCLUDE_PATH}/capi/astra_ctypes.h
  ${INCLUDE_PATH}/capi/streams/color_capi.h
  ${INCLUDE_PATH}/capi/streams/color_types.h
  ${INCLUDE_PATH}/capi/streams/depth_capi.h
  ${INCLUDE_PATH}/capi/streams/depth_parameters.h
  ${INCLUDE_PATH}/capi/streams/depth_types.h
  ${INCLUDE_PATH}/capi/streams/hand_capi.h
  ${INCLUDE_PATH}/capi/streams/hand_parameters.h
  ${INCLUDE_PATH}/capi/streams/hand_types.h
  ${INCLUDE_PATH}/capi/streams/image_capi.h
  ${INCLUDE_PATH}/capi/streams/image_parameters.h
  ${INCLUDE_PATH}/capi/streams/image_types.h
  ${INCLUDE_PATH}/capi/streams/infrared_capi.h
  ${INCLUDE_PATH}/capi/streams/infrared_types.h
  ${INCLUDE_PATH}/capi/streams/point_capi.h
  ${INCLUDE_PATH}/capi/streams/point_types.h
  ${INCLUDE_PATH}/capi/streams/skeleton_capi.h
  ${INCLUDE_PATH}/capi/streams/skeleton_types.h
  ${INCLUDE_PATH}/capi/streams/stream_types.h
  ${INCLUDE_PATH}/streams/Color.hpp
  ${INCLUDE_PATH}/streams/DebugHand.hpp
  ${INCLUDE_PATH}/streams/Depth.hpp
  ${INCLUDE_PATH}/streams/Hand.hpp
  ${INCLUDE_PATH}/streams/Image.hpp
  ${INCLUDE_PATH}/streams/Infrared.hpp
  ${INCLUDE_PATH}/streams/Point.hpp
  ${INCLUDE_PATH}/streams/Skeleton.hpp
  astra_generic_stream_api.hpp
  astra_depth.cpp
  astra_color.cpp
  astra_infrared.cpp
  astra_image.cpp
  astra_hand.cpp
  astra_point.cpp
  astra_skeleton.cpp
  astra.cpp
  )

add_library(${_projname} SHARED ${${_projname}_SOURCES})

add_definitions(-DASTRA_BUILD_EX)

target_link_libraries(${_projname} astra_core_api Shiny)

include_directories(${_projname} ${SHINY_INCLUDE_DIR})

install_lib(${_projname})