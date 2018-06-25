#!/bin/bash

TOOLCHAIN=xcode-hid-sections
CONFIG=Release

tf_opts=(
    CUDA_TOOLKIT_ROOT_DIR=/Developer/NVIDIA/CUDA-9.0
    CUDA_TOOLKIT_TARGET_DIR=/Developer/NVIDIA/CUDA-9.0
    tensorflow_PATH_STATIC_LIB=/Developer/NVIDIA/CUDA-9.0/lib
    tensorflow_CUDNN_INCLUDE=$(dirname $(dirname ${CUDNN_PATH}))/include
    tensorflow_PATH_CUDNN_STATIC_LIB=$(dirname ${CUDNN_PATH})
    tensorflow_BUILD_CC_TESTS=OFF
    htt_ENABLE_GPU=ON    
    htt_ENABLE_GRPC_SUPPORT=ON
)

polly.py --toolchain ${TOOLCHAIN} --config ${CONFIG} --fwd ${tf_opts[@]} --verbose ${*} --jobs 8
