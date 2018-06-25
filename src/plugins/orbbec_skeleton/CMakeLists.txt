include_directories(${ASTRA_INCLUDE_DIR})

set(ORBBEC_SKELETON_SRC
  orbbec_skeleton_plugin.hpp
  orbbec_skeleton_plugin.cpp
  orbbec_skeletonstream.hpp
  orbbec_skeleton_tracker.hpp
  orbbec_skeleton_tracker.cpp
  )

add_library(orbbec_skeleton SHARED ${ORBBEC_SKELETON_SRC})
target_link_libraries(orbbec_skeleton ${ASTRA_LIBRARIES})
set_target_properties(orbbec_skeleton PROPERTIES FOLDER "plugins")
