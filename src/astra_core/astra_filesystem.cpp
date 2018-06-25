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
#include "astra_filesystem.hpp"
#include "astra_string.hpp"

#include <string>

namespace astra { namespace filesystem {

    bool ends_with_path_separator(const std::string& path)
    {
        return path.length() > 0 && path[path.length() - 1] == filesystem::path_separator;
    }

    std::string trim_path_separator(const std::string& path)
    {
        return trim_end(path, filesystem::path_separator);
    }

    std::string append_path_separator(const std::string& path)
    {
        return trim_path_separator(path) + filesystem::path_separator;
    }

    std::string combine_paths(const std::string& left, const std::string& right)
    {
        return trim_path_separator(left) + filesystem::path_separator + right;
    }
}}
