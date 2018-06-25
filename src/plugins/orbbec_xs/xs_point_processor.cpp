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
#include "xs_point_processor.hpp"
#include <Shiny.h>

namespace astra { namespace xs {

    point_processor::point_processor(PluginServiceProxy& pluginService,
                                     astra_streamset_t streamset,
                                     StreamDescription& depthDesc)
        : streamset_(plugins::get_uri_for_streamset(pluginService, streamset)),
          setHandle_(streamset),
          reader_(streamset_.create_reader()),
          depthStream_(reader_.stream<DepthStream>(depthDesc.subtype())),
          pluginService_(pluginService)
    {
        depthStream_.start();
        reader_.add_listener(*this);
    }

    point_processor::~point_processor() = default;

    void point_processor::on_frame_ready(StreamReader& reader, Frame& frame)
    {
        const DepthFrame depthFrame = frame.get<DepthFrame>();

        create_point_stream_if_necessary(depthFrame);

        if (pointStream_->has_connections())
        {
            LOG_TRACE("astra.xs.point_processor", "updating point frame");
            update_pointframe_from_depth(depthFrame);
        }
    }

    void point_processor::create_point_stream_if_necessary(const DepthFrame& depthFrame)
    {
        if (pointStream_) { return; }

        //TODO check for changes in depthFrame width and height and update bin size
        LOG_INFO("astra.xs.point_processor", "creating point stream");

        int width = depthFrame.width();
        int height = depthFrame.height();

        auto ps = plugins::make_stream<PointStream>(pluginService_, setHandle_, width, height);
        pointStream_ = std::unique_ptr<PointStream>(std::move(ps));

        LOG_INFO("astra.xs.point_processor", "created point stream");

        depthConversionCache_ = depthStream_.depth_to_world_data();
    }

    void point_processor::update_pointframe_from_depth(const DepthFrame& depthFrame)
    {
        //use same frameIndex as source depth frame
        astra_frame_index_t frameIndex = depthFrame.frame_index();

        astra_imageframe_wrapper_t* pointFrameWrapper = pointStream_->begin_write(frameIndex);

        if (pointFrameWrapper != nullptr)
        {
            pointFrameWrapper->frame.frame = nullptr;
            pointFrameWrapper->frame.data = &pointFrameWrapper->frame_data[0];

            astra_image_metadata_t metadata;

            metadata.width = depthFrame.width();
            metadata.height = depthFrame.height();
            metadata.pixelFormat = ASTRA_PIXEL_FORMAT_POINT;

            pointFrameWrapper->frame.metadata = metadata;

            Vector3f* p_points = reinterpret_cast<Vector3f*>(pointFrameWrapper->frame.data);
            calculate_point_frame(depthFrame, p_points);

            pointStream_->end_write();
        }
    }

    void point_processor::calculate_point_frame(const DepthFrame& depthFrame,
                                                Vector3f* p_points)
    {
        int width = depthFrame.width();
        int height = depthFrame.height();
        const int16_t* p_depth = depthFrame.data();

        const conversion_cache_t conversionData = depthConversionCache_;

        for (int y = 0; y < height; ++y)
        {
            for (int x = 0; x < width; ++x, ++p_points, ++p_depth)
            {
                uint16_t depth = *p_depth;
                Vector3f& point = *p_points;

                float normalizedX = static_cast<float>(x) / conversionData.resolutionX - .5f;
                float normalizedY = .5f - static_cast<float>(y) / conversionData.resolutionY;

                point.x = normalizedX * depth * conversionData.xzFactor;
                point.y = normalizedY * depth * conversionData.yzFactor;
                point.z = depth;
            }
        }
    }
}}
