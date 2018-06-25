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
#include "orbbec_skeleton_plugin.hpp"
#include <astra/astra.hpp>

EXPORT_PLUGIN(orbbec::skeleton::skeleton_plugin);

namespace orbbec { namespace skeleton {

    void skeleton_plugin::on_stream_added(astra_streamset_t setHandle,
                                          astra_stream_t streamHandle,
                                          astra_stream_desc_t desc)
    {
        if (desc.type != ASTRA_STREAM_DEPTH)
            return; // if new stream is not depth, we don't care.

        LOG_DEBUG("orbbec.skeleton.skeleton_plugin", "creating skeleton tracker for %p", streamHandle);
        skeletonTrackers_.push_back(astra::make_unique<skeleton_tracker>(pluginService(),
                                                                       setHandle,
                                                                       streamHandle));
    }

    void skeleton_plugin::on_stream_removed(astra_streamset_t setHandle,
                                            astra_stream_t streamHandle,
                                            astra_stream_desc_t desc)
    {
        if (desc.type != ASTRA_STREAM_DEPTH)
            return;

        LOG_DEBUG("orbbec.skeleton.skeleton_plugin", "looking for skeleton tracker for %p", streamHandle);
        auto it = std::find_if(skeletonTrackers_.cbegin(),
                               skeletonTrackers_.cend(),
                               [&streamHandle] (const skeleton_trackerPtr& trackerPtr)
                               {
                                   return trackerPtr->sourceStream() == streamHandle;
                               });

        LOG_DEBUG("orbbec.skeleton.skeleton_plugin", "destroying skeleton tracker for %p", streamHandle);
        if (it != skeletonTrackers_.cend())
        {
            skeletonTrackers_.erase(it);
        }
    }
}}
