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

#include <string>

namespace astra { namespace filesystem {

    std::string directory_name(const std::string& filePath)
    {
        const size_t PATH_LENGTH = 1024;
        char fullPath[PATH_LENGTH];
        char* fileNameBegin;

        DWORD result = ::GetFullPathName(filePath.c_str(), PATH_LENGTH, fullPath, &fileNameBegin);

        if (!result)
        {
            return std::string();
        }

        fileNameBegin[0] = '\0';
        std::string directoryName(trim_path_separator(fullPath));

        return directoryName;
    }
}}
