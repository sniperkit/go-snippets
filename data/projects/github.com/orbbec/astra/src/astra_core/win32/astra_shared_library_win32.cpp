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
#define WIN32_LEAN_AND_MEAN

#include <Windows.h>

#include "../astra_shared_library.hpp"

namespace astra { namespace process {

    void load_library(const char* fileName, lib_handle& libHandle)
    {
        char strFileName[ASTRA_FILE_MAX_PATH];
        ::GetFullPathName(fileName, ASTRA_FILE_MAX_PATH, strFileName, nullptr);
        libHandle = LoadLibraryEx(strFileName, nullptr, LOAD_WITH_ALTERED_SEARCH_PATH);

        if (libHandle == nullptr)
        {
            // error
            //use GetLastError()
            return;
        }
    }

    void free_library(const lib_handle libHandle)
    {
        if (::FreeLibrary((HMODULE)libHandle) != 0)
        {
            //error
            //use GetLastError()
            return;
        }
    }

    void get_proc_address(const lib_handle libHandle, const char* procName, far_proc& procAddr)
    {
        // Get the requested procedure address from the shared library via the OS
        procAddr = (far_proc)::GetProcAddress((HMODULE)libHandle, procName);
        if (procAddr == nullptr)
        {
            //error
            //use GetLastError()
        }
    }

}}