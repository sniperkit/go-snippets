set(JNI_LIB_ROOT "${CMAKE_CURRENT_BINARY_DIR}/lib/")
set(JNI_LIB_DESTINATION "${JNI_LIB_ROOT}/armeabi-v7a/")

make_directory(${JNI_LIB_DESTINATION})

add_custom_target(AstraAndroid ALL
  COMMAND ${CMAKE_COMMAND}
  -D"INPATH=${CMAKE_SOURCE_DIR}/vendor/openni2/android-armeabi-v7a"
  -D"OUTPATH=${JNI_LIB_DESTINATION}"
  -D"PATTERN=*.so"
  -P "CopyFiles.cmake"
  COMMAND ${CMAKE_COMMAND}
  -D"INPATH=${CMAKE_BINARY_DIR}/lib/"
  -D"OUTPATH=${JNI_LIB_DESTINATION}"
  -D"PATTERN=*.so"
  -P "CopyFiles.cmake"
  COMMAND ${CMAKE_COMMAND}
  -D"INPATH=${CMAKE_BINARY_DIR}/bin/Plugins/"
  -D"OUTPATH=${JNI_LIB_DESTINATION}"
  -D"PATTERN=*.so"
  -P "CopyFiles.cmake"
  COMMAND "${CMAKE_CURRENT_SOURCE_DIR}/gradlew"
  "build"
  "-q"
  "-PoutputBinaryDir=${CMAKE_CURRENT_BINARY_DIR}/build"
  "-PandroidNativeApiLevel=${ANDROID_NATIVE_API_LEVEL}"
  "-PnativeLibPath=${JNI_LIB_ROOT}"
  WORKING_DIRECTORY "${CMAKE_CURRENT_SOURCE_DIR}")

add_dependencies(AstraAndroid astra_jni ${SDK_DEPENDENT_TARGET})

install(FILES ${CMAKE_CURRENT_BINARY_DIR}/build/outputs/aar/astra-android-release.aar DESTINATION aar/)
