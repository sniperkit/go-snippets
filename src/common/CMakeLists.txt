include(ConfigUserFile)

set(COMMON_DIR_FOLDER "common/")

if (NOT ASTRA_ANDROID)
    add_subdirectory(serialization)
endif()

add_subdirectory(clock)
