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
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>
#include <dlfcn.h>

#include "../astra_shared_library.hpp"
#include "../astra_logger.hpp"

namespace astra { namespace process {

    void load_library(const char* fileName, lib_handle& libHandle)
    {
        libHandle = ::dlopen(fileName, RTLD_NOW);

        if (libHandle == nullptr)
        {
            // error
            const char* err = ::dlerror();
            LOG_TRACE("shared_library_unix", "load_library failed: %s", err);
            return;
        }
    }

    void free_library(const lib_handle libHandle)
    {
        if (::dlclose(libHandle) != 0)
        {
            //error
            const char* err = ::dlerror();
            LOG_TRACE("shared_library_unix", "load_library failed: %s", err);
            return;
        }
    }

    void get_proc_address(const lib_handle libHandle, const char* procName, far_proc& procAddr)
    {
        // Get the requested procedure address from the shared library via the OS
        procAddr = (far_proc)::dlsym(libHandle, procName);
    }

}}
