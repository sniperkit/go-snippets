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
#ifndef MOCKSENSORPLUGIN_H
#define MOCKSENSORPLUGIN_H

#include <astra_core/astra_core.hpp>
#include <astra_core/plugins/PluginBase.hpp>
#include <astra_core/plugins/PluginLogging.hpp>
#include "mock_device_stream.hpp"
#include "mock_device_streamset.hpp"
#include <memory>
#include <vector>
#include <astra/capi/streams/depth_types.h>
#include <astra/capi/streams/color_types.h>
#include <astra/capi/streams/stream_types.h>

namespace orbbec { namespace mocks {

    class mock_sensor_plugin : public astra::plugins::plugin_base
    {
    public:
        mock_sensor_plugin(astra::PluginServiceProxy* pluginService)
            : plugin_base(pluginService, "mock_sensor")
        {
            register_for_host_events();

            add_or_get_device("mock_sensor");
        }

        virtual ~mock_sensor_plugin();
        virtual void temp_update() override;

        mock_sensor_plugin(const mock_sensor_plugin&) = delete;
        mock_sensor_plugin& operator=(const mock_sensor_plugin&) = delete;

    private:
        virtual void on_host_event(astra_event_id id, const void* data, size_t dataSize) override;

        device_streamset* add_or_get_device(const char* testUri);
        device_streamset* find_device(const char* testUri);

        astra_status_t read_streams();

        using streamset_ptr = std::unique_ptr<device_streamset>;
        std::vector<streamset_ptr> streamsets_;
    };
}}


#endif /* MOCKSENSORPLUGIN_H */
