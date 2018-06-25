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
#ifndef MOCK_DEVICE_COLOR_STREAM_H
#define MOCK_DEVICE_COLOR_STREAM_H

#include <memory>
#include <chrono>

#include "mock_base_stream.hpp"
#include "mock_color_generator.hpp"

using astra::devices::device_status;
using astra::devices::device_status_value;

namespace orbbec { namespace mocks {

    class mock_color_stream : public orbbec::mocks::base_stream
    {
    public:
        mock_color_stream();
        virtual ~mock_color_stream();

        virtual device_status on_initialize() override;
        virtual device_status on_read_into(void* dest, std::size_t size, std::int32_t timeout) override;

    private:
        std::unique_ptr<color_generator> generator_;
    };
}}

#endif /* MOCK_DEVICE_COLOR_STREAM_H */
