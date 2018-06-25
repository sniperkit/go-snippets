set (_projname "DepthReaderEventCPP")

set (${_projname}_SOURCES
  main.cpp
  ${SAMPLE_COMMON_DIR}/key_handler.h
  )

add_executable(${_projname} ${${_projname}_SOURCES})

set_target_properties(${_projname} PROPERTIES FOLDER "${SAMPLE_DIR_FOLDER}cpp-api")

target_link_libraries(${_projname} ${ASTRA_LIBRARIES})

set_debug_working_dir(${_projname} "$(TargetDir)")

add_astra_sdk_dependency(${_projname})

install_sample(${_projname})
