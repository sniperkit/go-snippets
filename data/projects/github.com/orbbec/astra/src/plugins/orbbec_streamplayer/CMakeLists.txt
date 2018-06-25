set (_projname "orbbec_streamplayer")

set(${_projname}_SOURCES
  StreamPlayerPlugin.cpp
  StreamPlayerPlugin.h
  DepthStream.h
  PlaybackStreamSet.h
  PlaybackStream.h
)

add_library(${_projname} SHARED ${${_projname}_SOURCES})

set_target_properties(${_projname} PROPERTIES FOLDER "plugins")

target_link_libraries(${_projname} astra_core_api astra FrameSerialization ClockUtil)

include_directories(${_projname})