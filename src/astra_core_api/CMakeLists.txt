set (_projname "astra_core_api")

set(${_projname}_SOURCES
  astra_core_api.h
  astra_core_api.cpp
  )

set(${_projname}_SOURCES_GEN
  astra_core_api.cpp.lpp
  )

add_definitions(-DASTRA_BUILD_API_PROXY)
add_definitions(-DASTRA_BUILD)

source_group(templates FILES ${${_projname}_SOURCES_GEN})

add_library(${_projname} SHARED ${${_projname}_SOURCES} ${${_projname}_SOURCES_GEN})

install_lib(${_projname})