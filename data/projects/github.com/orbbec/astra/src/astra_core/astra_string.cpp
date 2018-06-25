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
#include "astra_string.hpp"

namespace astra {

    std::string trim_start(const std::string& str, char trimChar)
    {
        if (str.length() == 0)
            return str;

        auto it = str.cbegin();
        while(it++ != str.cbegin() && *it == trimChar) {}

        return std::string(it - 1, str.cend());
    }

    std::string trim_end(const std::string& str, char trimChar)
    {
        if (str.length() == 0)
            return str;

        auto it = str.cend();
        while(--it != str.cbegin() && *it == trimChar) {}

        return std::string(str.cbegin(), it + 1);
    }
}
