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
#ifndef ONI_INFRARED_STREAM_H
#define ONI_INFRARED_STREAM_H

#include "oni_devicestream.hpp"
#include <astra_core/capi/plugins/astra_plugin.h>
#include <astra/capi/astra_ctypes.h>
#include <astra/capi/streams/stream_types.h>

namespace orbbec { namespace ni {

    class infrared_stream : public devicestream<astra_imageframe_wrapper_t>
    {
    public:
        infrared_stream(astra::PluginServiceProxy& pluginService,
                        astra_streamset_t streamSet,
                        openni::Device& oniDevice,
                        stream_listener& listener);

        infrared_stream(const infrared_stream&) = delete;
        infrared_stream& operator=(const infrared_stream&) = delete;
    };
}}

#endif /* ONI_INFRARED_STREAM_H */
