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
#include "orbbec_skeleton_tracker.hpp"
#include <astra/capi/streams/skeleton_parameters.h>
#include <cstddef>

namespace orbbec { namespace skeleton {

    const std::size_t skeleton_tracker::MAX_SKELETONS = 6;

    void skeleton_tracker::on_set_parameter(astra::plugins::stream* stream,
                                            astra_streamconnection_t connection,
                                            astra_parameter_id id,
                                            size_t inByteLength,
                                            astra_parameter_data_t inData)
    {
        switch(id)
        {
        case ASTRA_PARAMETER_SKELETON_Z_MIN:
        {
            std::uint16_t* zMin = reinterpret_cast<std::uint16_t*>(inData);
            zMin_ = *zMin;
            break;
        }
        case ASTRA_PARAMETER_SKELETON_Z_MAX:
        {
            std::uint16_t* zMax = reinterpret_cast<std::uint16_t*>(inData);
            zMax_ = *zMax;
            break;
        }
        default:
            break;
        }
    }

    void skeleton_tracker::on_get_parameter(astra::plugins::stream* stream,
                                            astra_streamconnection_t connection,
                                            astra_parameter_id id,
                                            astra_parameter_bin_t& parameterBin)
    {
        switch(id)
        {
        case ASTRA_PARAMETER_SKELETON_Z_MIN:
        {
            std::size_t resultByteLength = sizeof(std::uint16_t);

            astra_parameter_data_t parameterData;
            astra_status_t rc = pluginService_.get_parameter_bin(resultByteLength,
                                                                 &parameterBin,
                                                                 &parameterData);
            if (rc == ASTRA_STATUS_SUCCESS)
            {
                std::uint16_t* zMin = reinterpret_cast<std::uint16_t*>(parameterData);
                *zMin = zMin_;
            }

            break;
        }
        case ASTRA_PARAMETER_SKELETON_Z_MAX:
        {
            std::size_t resultByteLength = sizeof(std::uint16_t);

            astra_parameter_data_t parameterData;
            astra_status_t rc = pluginService_.get_parameter_bin(resultByteLength,
                                                                 &parameterBin,
                                                                 &parameterData);
            if (rc == ASTRA_STATUS_SUCCESS)
            {
                std::uint16_t* zMax = reinterpret_cast<std::uint16_t*>(parameterData);
                *zMax = zMax_;
            }

            break;
        }
        default:
            break;
        }
    }

    void skeleton_tracker::on_frame_ready(astra::StreamReader& reader, astra::Frame& frame)
    {
        if (!skeletonStream_->has_connections())
            return; // don't waste cycles if no one is listening

        LOG_TRACE("orbbec.skeleton.skeleton_tracker", "generating skeleton frame");
        const astra::DepthFrame depthFrame = frame.get<astra::DepthFrame>();

        if (!depthFrame.is_valid())
            return;

        // do something cool
        astra_skeletonframe_wrapper_t* skeletonFrame = skeletonStream_->begin_write(depthFrame.frame_index());

        if (skeletonFrame != nullptr)
        {
            skeletonFrame->frame.skeletons = reinterpret_cast<astra_skeleton_t*>(&(skeletonFrame->frame_data));
            skeletonFrame->frame.skeletonCount = skeleton_tracker::MAX_SKELETONS;

            astra_skeleton_t& skeleton = skeletonFrame->frame.skeletons[0];
            skeleton.trackingId = 1;
            skeleton.status = ASTRA_SKELETON_STATUS_TRACKED;
            skeleton.jointCount = ASTRA_MAX_JOINTS;
            skeleton.joints[0].status = ASTRA_JOINT_STATUS_TRACKED;
            skeleton.joints[0].jointType = ASTRA_JOINT_TYPE_LEFT_HAND;
            skeleton.joints[0].position.x = 0;
            skeleton.joints[0].position.y = 0;
            skeleton.joints[0].position.z = 2000;

            for(unsigned i = 1; i < MAX_SKELETONS; i++)
            {
                astra_skeleton_t& skeleton = skeletonFrame->frame.skeletons[i];
                skeleton.trackingId = i + 1;
                skeleton.status = ASTRA_SKELETON_STATUS_NOT_TRACKED;
            }

            skeletonStream_->end_write();
        }
    }
}}
