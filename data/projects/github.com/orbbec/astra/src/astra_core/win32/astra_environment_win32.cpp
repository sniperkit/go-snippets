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
#include "../astra_filesystem.hpp"

#include <windows.h>
#include <string>

namespace astra { namespace environment {

    namespace detail {

        std::string module_path_for_proc_address(void* procAddress)
        {
            HMODULE hModule;
            BOOL rc = ::GetModuleHandleEx(GET_MODULE_HANDLE_EX_FLAG_FROM_ADDRESS,
                                          (LPCSTR)procAddress,
                                          &hModule);

            if (!rc)
            {
                return std::string();
            }

            std::string path;
            const size_t PATH_LENGTH = 1024;
            path.resize(PATH_LENGTH);

            DWORD len = ::GetModuleFileName(hModule, &path[0], PATH_LENGTH);

            if (len == 0)
            {
                return std::string();
            }

            path.resize(len);

            return path;
        }
    }

    std::string lib_path()
    {
        return filesystem::directory_name(detail::module_path_for_proc_address(reinterpret_cast<void*>(&lib_path)));
    }
}}
