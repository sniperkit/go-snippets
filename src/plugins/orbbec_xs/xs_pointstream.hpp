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
#ifndef XS_POINTSTREAM_H
#define XS_POINTSTREAM_H

#include <astra_core/plugins/SingleBinStream.hpp>
#include <astra/capi/streams/point_types.h>
#include <astra/capi/astra_ctypes.h>
#include <astra/capi/streams/stream_types.h>
#include <Shiny.h>

namespace astra { namespace xs {

    class PointStream : public astra::plugins::single_bin_stream<astra_imageframe_wrapper_t>
    {
    public:
        PointStream(PluginServiceProxy& pluginService,
                    astra_streamset_t streamSet,
                    uint32_t width,
                    uint32_t height)
            : single_bin_stream(pluginService,
                                streamSet,
                                StreamDescription(ASTRA_STREAM_POINT,
                                                  DEFAULT_SUBTYPE),
                                width * height * sizeof(astra_vector3f_t))
        {}
    };
}}

#endif /* XS_POINTSTREAM_H */
