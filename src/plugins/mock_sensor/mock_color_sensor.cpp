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
#include "mock_color_sensor.hpp"
#include "mock_device_color_stream.hpp"

namespace orbbec { namespace mocks {

    mock_color_sensor::mock_color_sensor()
        : sensor(sensor_info(astra::devices::sensor_type::color, "color"))
    {
    }

    mock_color_sensor::~mock_color_sensor()
    {
    }

    device_status mock_color_sensor::on_initialize()
    {
        auto color = std::make_shared<mock_color_stream>();
        add_stream(color);

        return device_status_value::ok;
    }
}}
