set(XS_SRC
  xs_plugin.cpp
  xs_point_processor.cpp
  )

set(XS_INCLUDE
  xs_plugin.hpp
  xs_point_processor.hpp
  xs_pointstream.hpp
  )

include_directories(orbbec_xs ${SHINY_INCLUDE_DIR})
add_library(orbbec_xs SHARED ${XS_SRC} ${XS_INCLUDE})
target_link_libraries(orbbec_xs astra_core_api astra Shiny)

set_target_properties(orbbec_xs PROPERTIES FOLDER "plugins")
install_lib(orbbec_xs "Plugins")
