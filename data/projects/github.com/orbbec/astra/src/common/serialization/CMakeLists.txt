set (_projname "FrameSerialization")

set (${_projname}_SOURCES
  ../../../include/common/serialization/FrameInputStream.h
  ../../../include/common/serialization/FrameOutputStream.h
  ../../../include/common/serialization/FrameStreamReader.h
  ../../../include/common/serialization/FrameStreamWriter.h
  FrameStreamReader.cpp
  pb_util.cpp
  pb_util.h
  ProtoFrameInputStream.h
  ProtoFrameInputStream.cpp
  ProtoFrameOutputStream.h
  ProtoFrameOutputStream.cpp
  ../../../include/common/serialization/StreamFileModels.h
  ../../../include/common/serialization/FrameStreamWriter.h
  FrameStreamWriter.cpp
  )

find_package(Protobuf REQUIRED)

include_directories(${PROTOBUF_INCLUDE_DIRS})
include_directories(${CMAKE_CURRENT_BINARY_DIR})

set(PROTOBUF_DEFINITIONS
  proto/StreamHeader.proto
  proto/FrameDescription.proto
  proto/Frame.proto
  )

protobuf_generate_cpp(PROTO_SRCS PROTO_HDRS ${PROTOBUF_DEFINITIONS})
SET_SOURCE_FILES_PROPERTIES(${PROTO_SRC} ${PROTO_INCL} PROPERTIES GENERATED TRUE)

add_custom_target(autogen_pb_frame_serialization DEPENDS ${PROTO_SRCS} ${PROTO_HDRS})

set_target_properties(autogen_pb_frame_serialization PROPERTIES FOLDER "${COMMON_DIR_FOLDER}serialization")

set(${_projname}_SOURCES
  ${${_projname}_SOURCES}
  ${PROTO_SRCS}
  ${PROTO_HDRS}
  )

if (MSVC)
  add_library(${_projname} STATIC ${${_projname}_SOURCES})
else()
  add_library(${_projname} SHARED ${${_projname}_SOURCES})
endif()

set_target_properties(${_projname} PROPERTIES FOLDER "${COMMON_DIR_FOLDER}serialization")
target_link_libraries(${_projname} ${PROTOBUF_LIBRARIES} ClockUtil)

add_dependencies(${_projname} autogen_pb_frame_serialization)
