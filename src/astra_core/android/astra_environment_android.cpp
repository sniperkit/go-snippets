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
#include "../astra_environment.hpp"

namespace astra { namespace environment {

    std::string application_name()
    {
        FILE *file;

        const size_t NAME_LENGTH = 1024;
        char appName[NAME_LENGTH];

        file = fopen("/proc/self/cmdline", "r");
        fread(appName, 1, NAME_LENGTH, file);
        fclose(file);

        return std::string(appName);
    }

    std::string application_path()
    {
        const size_t PATH_LENGTH = 1024;
        char path[PATH_LENGTH];

        std::string appName = application_name();
        snprintf(path, PATH_LENGTH, "/data/data/%s/files/", appName.c_str());

        return std::string(path);
    }

    std::string libpath()
    {
        const size_t PATH_LENGTH = 1024;
        char path[PATH_LENGTH];

        std::string appName = application_name();
        snprintf(path, PATH_LENGTH, "/data/data/%s/lib/", appName.c_str());

        return std::string(path);
    }

}}
