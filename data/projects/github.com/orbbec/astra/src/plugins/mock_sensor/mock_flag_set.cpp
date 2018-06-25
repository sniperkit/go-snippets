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
#include "mock_flag_set.hpp"
#include <iostream>

namespace orbbec { namespace mocks {

    void flag_set::add_flag(const std::string& flag)
    {
        if (flag.length() > 0)
            flags_.insert(flag);
    }

    void flag_set::remove_flag(const std::string& flag)
    {
        if (flag.length() > 0)
            flags_.erase(flag);
    }

    bool flag_set::has_flag(const std::string& flag)
    {
        return flags_.find(flag) != flags_.end();
    }

    void flag_set::set_flag(const std::string& flag, bool isSet)
    {
        if (isSet)
        {
            add_flag(flag);
        }
        else
        {
            remove_flag(flag);
        }
    }

    std::ostream& operator<<(std::ostream& os, const flag_set& fs)
    {
        for(auto& flag : fs.flags_)
        {
            os << flag << " ";
        }

        os.seekp(-1, std::ios_base::end);

        return os;
    }
}}
