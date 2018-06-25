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
#include <string>
#include <dlfcn.h>
#include <libgen.h>
#include "../astra_filesystem.hpp"

namespace astra { namespace environment {

    std::string module_path_for_proc_address(void* procAddress)
    {
        Dl_info info;
        if (!dladdr(procAddress, &info))
        {
            return std::string();
        }

        return std::string(info.dli_fname);
    }

    std::string lib_path()
    {
        return astra::filesystem::directory_name(module_path_for_proc_address(reinterpret_cast<void*>(&lib_path)));
    }
}}
