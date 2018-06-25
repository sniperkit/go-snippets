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
#ifndef XS_POINT_PROCESSOR_H
#define XS_POINT_PROCESSOR_H

#include <astra_core/plugins/Plugin.hpp>
#include <astra/astra.hpp>
#include "xs_pointstream.hpp"

namespace astra { namespace xs {

    class point_processor : public FrameListener
    {
    public:
        point_processor(PluginServiceProxy& pluginService,
                        astra_streamset_t streamset,
                        StreamDescription& depthDesc);
        virtual ~point_processor();
        virtual void on_frame_ready(StreamReader& reader, Frame& frame) override;

    private:
        void create_point_stream_if_necessary(const DepthFrame& depthFrame);

        void update_pointframe_from_depth(const DepthFrame& depthFrame);
        void calculate_point_frame(const DepthFrame& depthFrame,
                                   Vector3f* p_points);

        StreamSet streamset_;
        astra_streamset_t setHandle_;
        StreamReader reader_;
        DepthStream depthStream_;
        PluginServiceProxy& pluginService_;

        using PointStreamPtr = std::unique_ptr<PointStream>;
        PointStreamPtr pointStream_;

        conversion_cache_t depthConversionCache_;
    };
}}

#endif // XS_POINT_PROCESSOR_H
