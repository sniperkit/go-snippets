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
#ifndef ASTRA_PARAMETER_BIN_H
#define ASTRA_PARAMETER_BIN_H

#include <memory>
#include <astra_core/capi/astra_types.h>
#include <cstring>

namespace astra {

    class parameter_bin
    {
    public:
        parameter_bin(size_t byteSize)
            : byteLength_(byteSize)
        {
            //TODO pooling
            data_ = DataPtr(new uint8_t[byteSize]);
	    std::memset(data_.get(), 0, byteSize);
        }

        size_t byteLength() { return byteLength_; }
        void* data() { return data_.get(); }

        astra_parameter_bin_t get_handle() { return reinterpret_cast<astra_parameter_bin_t>(this); }
        static parameter_bin* get_ptr(astra_parameter_bin_t bin)
        {
            return reinterpret_cast<parameter_bin*>(bin);
        }

    private:
        using DataPtr = std::unique_ptr<uint8_t[]>;

        DataPtr data_;
        size_t byteLength_;
    };
}

#endif /* ASTRA_PARAMETER_BIN_H */
