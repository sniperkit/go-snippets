// This file is part of the Orbbec Astra SDK [https://orbbec3d.com]
// Copyright (c) 2015 Orbbec 3D
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Be excellent to each other.
#include <jni.h>
#include <astra_core/astra_core.hpp>
#include <astra_core/capi/astra_host_events.h>
#include "astra_wrapper.h"

#ifdef __cplusplus
extern "C" {
#endif

    JNIEXPORT void JNICALL Java_com_orbbec_jni_Astra_initialize(JNIEnv* env, jclass cls)
    {
        astra_initialize();
    }

    JNIEXPORT void JNICALL Java_com_orbbec_jni_Astra_terminate(JNIEnv* env, jclass cls)
    {
        astra_terminate();
    }

    JNIEXPORT void JNICALL Java_com_orbbec_jni_Astra_notify_1resource_1available(JNIEnv* env,
                                                                                    jclass cls,
                                                                                    jstring uri)
    {
        const char* uriString = env->GetStringUTFChars(uri, NULL);
        astra_notify_resource_available(uriString);
        env->ReleaseStringUTFChars(uri, uriString);
    }

    JNIEXPORT void JNICALL Java_com_orbbec_jni_Astra_notify_1resource_1unavailable(JNIEnv* env,
                                                                                      jclass cls,
                                                                                      jstring uri)
    {
        const char* uriString = env->GetStringUTFChars(uri, NULL);
        astra_notify_resource_unavailable(uriString);
        env->ReleaseStringUTFChars(uri, uriString);
    }

#ifdef __cplusplus
}
#endif
