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
#ifndef HND_PLUGIN_H
#define HND_PLUGIN_H

#include <astra_core/plugins/Plugin.hpp>
#include <astra/astra.hpp>
#include <map>
#include "hnd_hand_tracker.hpp"
#include "hnd_settings.hpp"

namespace astra { namespace hand {

    class plugin : public plugins::plugin_base
    {
    public:
        plugin(PluginServiceProxy* pluginProxy);
        virtual ~plugin();

        virtual void temp_update() override { }

    protected:
        virtual void on_initialize() override;

    private:
        static void stream_registered_handler_thunk(void* clientTag,
                                                    astra_streamset_t setHandle,
                                                    astra_stream_t streamHandle,
                                                    astra_stream_desc_t desc);

        static void stream_unregistering_handler_thunk(void* clientTag,
                                                       astra_streamset_t setHandle,
                                                       astra_stream_t streamHandle,
                                                       astra_stream_desc_t desc);

        void stream_registered_handler(astra_streamset_t setHandle,
                                       astra_stream_t streamHandle,
                                       astra_stream_desc_t desc);
        void stream_unregistering_handler(astra_streamset_t setHandle,
                                          astra_stream_t streamHandle,
                                          astra_stream_desc_t desc);


        astra_callback_id_t streamAddedCallbackId_{0};
        astra_callback_id_t streamRemovingCallbackId_{0};

        using hand_tracker_map = std::map<astra_stream_t,
                                          hand_tracker*>;

        hand_tracker_map streamTrackerMap_;

        hand_settings settings_;
    };
}}

#endif /* HND_PLUGIN_H */
