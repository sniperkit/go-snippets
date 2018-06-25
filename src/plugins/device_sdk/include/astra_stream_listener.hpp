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
#ifndef ASTRA_STREAM_LISTENER_H
#define ASTRA_STREAM_LISTENER_H

#include "astra_sensor_stream.hpp"
#include <memory>

namespace astra { namespace devices {

    class stream_listener
    {
    public:
        virtual ~stream_listener() {}

        virtual void state_changed(sensor_stream::shared_ptr stream) {}
        virtual void property_changed(sensor_stream::shared_ptr stream, sensor_stream::property_id id) {}
        virtual void streaming_changed(sensor_stream::shared_ptr stream) {}
        virtual void new_frame_available(sensor_stream::shared_ptr stream) {}
    };
}}

#endif /* ASTRA_STREAM_LISTENER_H */
