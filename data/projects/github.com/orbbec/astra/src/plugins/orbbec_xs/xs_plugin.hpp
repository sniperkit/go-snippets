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
#ifndef XS_PLUGIN_H
#define XS_PLUGIN_H

#include <astra_core/plugins/Plugin.hpp>
#include <astra/astra.hpp>
#include "xs_point_processor.hpp"
#include <memory>
#include <map>

namespace astra { namespace xs {

    class plugin : public astra::plugins::plugin_base
    {
    public:
        plugin(PluginServiceProxy* pluginProxy)
            : plugin_base(pluginProxy, "orbbec_xs")
        {
            LOG_INFO("astra.xs.plugin", "Initializing xs plugin");
            register_for_stream_events();

        }

        virtual ~plugin()
        {
            unregister_for_stream_events();
            LOG_INFO("astra.xs.plugin", "Terminated xs plugin");
        }

    private:
        virtual void on_stream_added(astra_streamset_t setHandle,
                                     astra_stream_t streamHandle,
                                     astra_stream_desc_t streamDesc) override;

        virtual void on_stream_removed(astra_streamset_t setHandle,
                                       astra_stream_t streamHandle,
                                       astra_stream_desc_t streamDesc) override;

        using point_processor_ptr = std::unique_ptr<point_processor>;
        using point_processor_map = std::map<astra_stream_t,
                                              point_processor_ptr>;

        point_processor_map pointProcessorMap_;
    };
}}

#endif /* XS_PLUGIN_H */
