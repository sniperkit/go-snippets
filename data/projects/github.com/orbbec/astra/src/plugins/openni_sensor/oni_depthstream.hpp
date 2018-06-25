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
#ifndef ONI_DEPTHSTREAM_H
#define ONI_DEPTHSTREAM_H

#include "oni_devicestream.hpp"
#include <astra_core/capi/plugins/astra_plugin.h>
#include <astra/capi/streams/stream_types.h>
#include <astra/capi/astra_ctypes.h>
#include <cmath>
#include <Shiny.h>

namespace orbbec { namespace ni {

    class depthstream : public devicestream<astra_imageframe_wrapper_t>
    {
    public:
        depthstream(astra::PluginServiceProxy& pluginService,
                    astra_streamset_t streamSet,
                    openni::Device& oniDevice,
                    stream_listener& listener);

        depthstream(const depthstream&) = delete;
        depthstream& operator=(const depthstream&) = delete;

    private:
        conversion_cache_t conversionCache_;

        void refresh_conversion_cache(float horizontalFov,
                                      float verticalFov,
                                      int resolutionX,
                                      int resolutionY);

        virtual astra_status_t on_open() override;

        virtual void on_get_parameter(astra_streamconnection_t connection,
                                      astra_parameter_id id,
                                      astra_parameter_bin_t& parameterBin) override;

        void on_set_parameter(astra_streamconnection_t connection,
                              astra_parameter_id id,
                              size_t inByteLength,
                              astra_parameter_data_t inData) override;

        virtual void on_connection_removed(astra_bin_t bin,
                                           astra_streamconnection_t connection) override;
    };
}}

#endif /* ONIDEPTHSTREAM_H */
