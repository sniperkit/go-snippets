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
#include "xs_plugin.hpp"
#include <astra_core/astra_core.hpp>
#include <Shiny.h>

EXPORT_PLUGIN(astra::xs::plugin);

namespace astra { namespace xs {

    void plugin::on_stream_added(astra_streamset_t setHandle,
                                 astra_stream_t streamHandle,
                                 astra_stream_desc_t streamDesc)
    {
        if (streamDesc.type == ASTRA_STREAM_DEPTH &&
            pointProcessorMap_.find(streamHandle) == pointProcessorMap_.end())
        {
            LOG_INFO("astra.xs.plugin", "creating point processor");

            StreamDescription depthDescription = streamDesc;

            auto ppPtr = astra::make_unique<point_processor>(pluginService(),
                                                           setHandle,
                                                           depthDescription);

            pointProcessorMap_[streamHandle] = std::move(ppPtr);
        }
    }

    void plugin::on_stream_removed(astra_streamset_t setHandle,
                                   astra_stream_t streamHandle,
                                   astra_stream_desc_t streamDesc)
    {
        auto it = pointProcessorMap_.find(streamHandle);
        if (it != pointProcessorMap_.end())
        {
            pointProcessorMap_.erase(it);
        }
    }

}}
