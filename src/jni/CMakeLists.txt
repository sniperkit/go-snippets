set(C_SOURCES
  astra_wrapper.h
  astra_wrapper.cpp
  )

add_library(astra_jni SHARED ${C_SOURCES})
target_link_libraries(astra_jni ${ASTRA_LIBRARIES})

install(TARGETS astra_jni
        RUNTIME DESTINATION bin/
        LIBRARY DESTINATION lib/
        ARCHIVE DESTINATION lib/
        )
