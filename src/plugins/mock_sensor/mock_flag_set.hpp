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
#ifndef MOCK_FLAG_SET_H
#define MOCK_FLAG_SET_H

#include <set>
#include <string>
#include <iosfwd>

namespace orbbec { namespace mocks {

    class flag_set
    {
    public:
        void add_flag(const std::string& flag);
        void remove_flag(const std::string& flag);
        void set_flag(const std::string& flag, bool isSet);

        bool has_flag(const std::string& flag);
        bool has_flags() const { return flags_.size() > 0; }

    private:
        std::set<std::string> flags_;

        friend std::ostream& operator<<(std::ostream& os, const flag_set& fs);
    };
}}

#endif /* MOCK_FLAG_SET_H */
