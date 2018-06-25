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
#include "mock_device.hpp"
#include "mock_color_sensor.hpp"
#include "mock_depth_sensor.hpp"
#include "astra_device_info.hpp"

namespace orbbec { namespace mocks {

    mock_device::mock_device()
        : astra::devices::device(astra::devices::device_info("device/mock_device", 0, 0))
    {
    }

    mock_device::~mock_device() = default;

    device_status mock_device::on_initialize()
    {
        auto color = std::make_shared<mock_color_sensor>();
        add_sensor(color);

        auto depth = std::make_shared<mock_depth_sensor>();
        add_sensor(depth);

        return device_status_value::ok;
    }

    device_status mock_device::on_connect()
    {
        return device_status_value::ok;
    }

    device_status mock_device::on_disconnect()
    {
        return device_status_value::ok;
    }
}}
