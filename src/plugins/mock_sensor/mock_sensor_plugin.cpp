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
#include <astra/capi/astra_ctypes.h>
#include <cstring>
#include <sstream>

#include "mock_sensor_plugin.hpp"
#include "mock_device_stream.hpp"

EXPORT_PLUGIN(orbbec::mocks::mock_sensor_plugin)

namespace orbbec { namespace mocks {

        mock_sensor_plugin::~mock_sensor_plugin()
        {
            LOG_INFO("orbbec.mocks.mock_sensor_plugin", "shutting down mock sensor plugin");
        }

        void mock_sensor_plugin::on_host_event(astra_event_id id, const void* data, size_t dataSize)
        {
#ifdef __ANDROID__
            switch (id)
            {
            case ASTRA_EVENT_RESOURCE_AVAILABLE:
                const char* resourceUri = static_cast<const char*>(data);

                LOG_INFO("orbbec.mocks.mock_sensor_plugin", "resource uri received: %s", resourceUri);

                unsigned int vid = 0;
                unsigned int pid = 0;
                unsigned int bus = 0;
                unsigned int address = 0;

                int scanned = sscanf(resourceUri, "usb/%u/%u/%u/%u", &vid, &pid, &bus, &address);
            }
#endif
        }

        device_streamset* mock_sensor_plugin::add_or_get_device(const char* mockUri)
        {
            device_streamset* device = find_device(mockUri);

            if (device)
                return device;

            std::stringstream sstream;
            sstream << "device/mock_sensor" << streamsets_.size();

            streamset_ptr streamSet = astra::make_unique<device_streamset>(sstream.str(),
                                                                         pluginService(),
                                                                         mockUri);
            streamSet->open();
            device = streamSet.get();
            streamsets_.push_back(std::move(streamSet));

            return device;
        }

        device_streamset* mock_sensor_plugin::find_device(const char* mockUri)
        {
            auto it = std::find_if(streamsets_.begin(), streamsets_.end(),
                                   [&mockUri](const streamset_ptr& setPtr) -> bool
                                   {
                                       return setPtr->get_uri() == mockUri;
                                   });

            return it != streamsets_.end() ? it->get() : nullptr;
        }

        void mock_sensor_plugin::temp_update()
        {
            read_streams();
        }

        astra_status_t mock_sensor_plugin::read_streams()
        {
            for(auto& set : streamsets_)
            {
                set->read();
            }

            return ASTRA_STATUS_SUCCESS;
        }
    }}
