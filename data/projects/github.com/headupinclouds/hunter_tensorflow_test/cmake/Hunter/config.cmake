hunter_config(nsync VERSION ${HUNTER_nsync_VERSION} CMAKE_ARGS NSYNC_LANGUAGE=c++11)
hunter_config(c-ares VERSION ${HUNTER_c-ares_VERSION} CMAKE_ARGS CARES_STATIC=ON CARES_SHARED=OFF)
hunter_config(Eigen VERSION v0.0.0-429aa5254200-p0)

set(tensorflow_CMAKE_ARGS
  tensorflow_ENABLE_GPU=${htt_ENABLE_GPU}
  tensorflow_ENABLE_GRPC_SUPPORT=${htt_ENABLE_GRPC_SUPPORT}
  tensorflow_BUILD_CC_TESTS=OFF
  tensorflow_BUILD_PYTHON_TESTS=OFF 
  tensorflow_BUILD_CC_EXAMPLE=OFF
)

# Forward these variables if defined
set(tensorflow_CMAKE_FWD
  CUDA_TOOLKIT_ROOT_DIR
  CUDA_TOOLKIT_TARGET_DIR
  tensorflow_PATH_STATIC_LIB
  tensorflow_CUDNN_INCLUDE
  tensorflow_PATH_CUDNN_STATIC_LIB
  )

foreach(key ${tensorflow_CMAKE_FWD})
  if(${key})
    list(APPEND tensorflow_CMAKE_ARGS "${key}=${${key}}")
  endif()
endforeach()

foreach(key_value ${tensorflow_CMAKE_ARGS})
  message("TF: ${key_value}")
endforeach()

if(htt_USE_TF_SUBMODULE)
  hunter_config(tensorflow GIT_SUBMODULE "src/3rdparty/tensorflow" CMAKE_ARGS ${tensorflow_CMAKE_ARGS})
else()
  hunter_config(tensorflow VERSION ${HUNTER_tensorflow_VERSION} CMAKE_ARGS ${tensorflow_CMAKE_ARGS})
endif()
