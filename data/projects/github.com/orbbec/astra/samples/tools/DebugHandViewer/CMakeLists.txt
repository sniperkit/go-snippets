set (_projname "DebugHandViewer")

set (${_projname}_SOURCES
  main.cpp
  )

find_package(SFML COMPONENTS graphics window system)

add_executable(${_projname} ${${_projname}_SOURCES})

set_target_properties(${_projname} PROPERTIES FOLDER "${SAMPLE_DIR_FOLDER}tools")

include_directories(${SFML_INCLUDE_DIR})

target_link_libraries(${_projname} ${ASTRA_LIBRARIES} ${SFML_LIBRARIES})

set_debug_working_dir(${_projname} "$(TargetDir)")

add_astra_sdk_dependency(${_projname})

install_sample(${_projname})
