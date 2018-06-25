set (_projname "ClockUtil")

set (${_projname}_SOURCES
  ../../../include/common/clock/Pulser.h
  ../../../include/common/clock/Stdafx.h
  ../../../include/common/clock/Stopwatch.h
  Pulser.cpp
  Stopwatch.cpp
)

add_library(${_projname} ${${_projname}_SOURCES})

set_target_properties(${_projname} PROPERTIES FOLDER "${COMMON_DIR_FOLDER}clock")

set_debug_working_dir(${_projname} "$(TargetDir)")
