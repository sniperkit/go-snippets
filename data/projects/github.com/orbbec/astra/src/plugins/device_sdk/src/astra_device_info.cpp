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
#include "astra_device_info.hpp"
#include <cassert>

namespace astra { namespace devices {

    device_info::device_info(std::string uri, std::uint32_t vendorId, std::uint32_t productId)
        : uri_(uri), vendorId_(vendorId), productId_(productId)
    {
        assert(uri_.length() > 0);
    }

    const std::string& device_info::uri() const
    {
        return uri_;
    }

    const std::uint32_t& device_info::vendor_id() const
    {
        return vendorId_;
    }

    const std::uint32_t& device_info::product_id() const
    {
        return productId_;
    }

    bool operator==(const device_info& lhs, const device_info& rhs)
    {
        return lhs.uri() ==  rhs.uri();
    }

    bool operator!=(const device_info& lhs, const device_info& rhs)
    {
        return !(lhs == rhs);
    }
}}
