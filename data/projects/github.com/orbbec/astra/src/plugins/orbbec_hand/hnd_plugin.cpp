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
#include "hnd_settings.hpp"
#include "hnd_plugin.hpp"
#include "hnd_hand_tracker.hpp"
#include <Shiny.h>

EXPORT_PLUGIN(astra::hand::plugin);

namespace astra { namespace hand {

    const char HANDPLUGIN_CONFIG_FILE[] = "plugins/orbbec_hand.toml";

    plugin::plugin(PluginServiceProxy* pluginProxy)
        : plugin_base(pluginProxy, "orbbec_hand")
    {
        settings_ = parse_settings(HANDPLUGIN_CONFIG_FILE);
    }

    plugin::~plugin()
    {
#ifndef __ANDROID__
        PROFILE_UPDATE();
        PROFILE_OUTPUT("profile_orbbec_hand.txt");
#endif

        pluginService().unregister_stream_registered_callback(streamAddedCallbackId_);
        pluginService().unregister_stream_unregistering_callback(streamRemovingCallbackId_);
    }

    void plugin::on_initialize()
    {
        pluginService().register_stream_registered_callback(&plugin::stream_registered_handler_thunk,
                                                            this,
                                                            &streamAddedCallbackId_);

        pluginService().register_stream_unregistering_callback(&plugin::stream_unregistering_handler_thunk,
                                                               this,
                                                               &streamRemovingCallbackId_);
    }

    void plugin::stream_registered_handler_thunk(void* clientTag,
                                                 astra_streamset_t setHandle,
                                                 astra_stream_t streamHandle,
                                                 astra_stream_desc_t desc)
    {
        plugin* plugin = static_cast<hand::plugin*>(clientTag);
        plugin->stream_registered_handler(setHandle, streamHandle, desc);
    }

    void plugin::stream_unregistering_handler_thunk(void* clientTag,
                                                    astra_streamset_t setHandle,
                                                    astra_stream_t streamHandle,
                                                    astra_stream_desc_t desc)

    {
        plugin* plugin = static_cast<hand::plugin*>(clientTag);
        plugin->stream_unregistering_handler(setHandle, streamHandle, desc);
    }

    void plugin::stream_registered_handler(astra_streamset_t setHandle,
                                           astra_stream_t streamHandle,
                                           astra_stream_desc_t streamDesc)
    {
        if (streamDesc.type == ASTRA_STREAM_DEPTH &&
            streamTrackerMap_.find(streamHandle) == streamTrackerMap_.end())
        {
            StreamDescription depthDescription = streamDesc;

            hand_tracker* tracker = new hand_tracker(pluginService(),
                                                     setHandle,
                                                     depthDescription,
                                                     settings_);

            streamTrackerMap_[streamHandle] = tracker;
        }
    }

    void plugin::stream_unregistering_handler(astra_streamset_t setHandle,
                                              astra_stream_t streamHandle,
                                              astra_stream_desc_t desc)
    {
        auto it = streamTrackerMap_.find(streamHandle);
        if (it != streamTrackerMap_.end())
        {
            hand_tracker* tracker = it->second;

            delete tracker;
            streamTrackerMap_.erase(it);
        }
    }
}}
