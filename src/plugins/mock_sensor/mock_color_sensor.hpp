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
#ifndef MOCK_COLOR_SENSOR_H
#define MOCK_COLOR_SENSOR_H

#include "astra_sensor.hpp"

using astra::devices::device_status;
using astra::devices::device_status_value;
using astra::devices::sensor_info;

namespace orbbec { namespace mocks {

    class mock_color_sensor : public astra::devices::sensor
    {
    public:
        mock_color_sensor();
        virtual ~mock_color_sensor();

        device_status on_initialize() override;
    };
}}

#endif /* MOCK_COLOR_SENSOR_H */
