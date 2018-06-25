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
#include "mock_depth_sensor.hpp"
#include "mock_device_depth_stream.hpp"
#include "mock_device_ir_stream.hpp"

namespace orbbec { namespace mocks {

    mock_depth_sensor::mock_depth_sensor()
        : sensor(sensor_info(astra::devices::sensor_type::depth, "depth"))
    {
    }

    mock_depth_sensor::~mock_depth_sensor()
    {
    }

    device_status mock_depth_sensor::on_initialize()
    {
        auto depth = std::make_shared<mock_depth_stream>();
        add_stream(depth);

        auto ir = std::make_shared<mock_ir_stream>();
        add_stream(ir);

        return device_status_value::ok;
    }
}}
