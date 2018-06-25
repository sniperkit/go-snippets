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
#ifndef ORBBEC_SKELETON_PLUGIN_HPP
#define ORBBEC_SKELETON_PLUGIN_HPP

#include <astra_core/plugins/Plugin.hpp>
#include <astra/astra.hpp>
#include "orbbec_skeleton_tracker.hpp"
#include <memory>
#include <vector>

namespace orbbec { namespace skeleton {

    class skeleton_plugin : public astra::plugins::plugin_base
    {
    public:
        static const size_t MAX_SKELETONS = 5;

        skeleton_plugin(astra::PluginServiceProxy* pluginProxy)
            : plugin_base(pluginProxy, "orbbec_skeleton")
        {
            register_for_stream_events();
        }

        virtual ~skeleton_plugin()
        {
            LOG_INFO("orbbec.skeleton.skeleton_plugin", "poof destroyed");
            unregister_for_stream_events();
        }

    private:
        virtual void on_stream_added(astra_streamset_t setHandle,
                                     astra_stream_t streamHandle,
                                     astra_stream_desc_t desc) override;

        virtual void on_stream_removed(astra_streamset_t setHandle,
                                       astra_stream_t streamHandle,
                                       astra_stream_desc_t desc) override;

        using skeleton_trackerPtr = std::unique_ptr<skeleton_tracker>;
        using skeleton_trackerList = std::vector<skeleton_trackerPtr>;

        skeleton_trackerList skeletonTrackers_;
    };
}}


#endif /* ORBBEC_SKELETON_PLUGIN_HPP */
