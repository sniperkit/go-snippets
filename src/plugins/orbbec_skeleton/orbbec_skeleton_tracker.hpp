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
#ifndef ORBBEC_SKELETON_TRACKER_HPP
#define ORBBEC_SKELETON_TRACKER_HPP

#include <astra_core/plugins/Plugin.hpp>
#include <astra/astra.hpp>
#include <astra/capi/astra_ctypes.h>
#include <astra/capi/streams/skeleton_types.h>
#include "orbbec_skeletonstream.hpp"
#include <cstdint>

namespace orbbec { namespace skeleton {

    class skeleton_tracker : public astra::FrameListener,
                             public astra::plugins::stream_event_handler
    {
    public:
        static const size_t MAX_SKELETONS;

        skeleton_tracker(astra::PluginServiceProxy& pluginService,
                         astra_streamset_t streamSet,
                         astra_stream_t sourceStream)
            : sourceStreamHandle_(sourceStream),
              sensor_(astra::plugins::get_uri_for_streamset(pluginService, streamSet)),
              reader_(sensor_.create_reader()),
              pluginService_(pluginService)
        {
            depthStream_ = reader_.stream<astra::DepthStream>();
            depthStream_.start();

            reader_.add_listener(*this);

            LOG_DEBUG("orbbec.skeleton.skeleton_tracker", "creating skeleton stream for %p", sourceStreamHandle_);
            auto s = astra::plugins::make_stream<skeletonstream>(pluginService_,
                                                                 streamSet,
                                                                 skeleton_tracker::MAX_SKELETONS);
            s->set_handler(this);
            skeletonStream_ = std::unique_ptr<skeletonstream>(std::move(s));
        }

        virtual ~skeleton_tracker() override
        {
            skeletonStream_->set_handler(nullptr);
            reader_.remove_listener(*this);

            LOG_DEBUG("orbbec.skeleton.skeleton_tracker", "destroying skeleton tracker for %p", sourceStreamHandle_);
        }

        astra_stream_t sourceStream() { return sourceStreamHandle_; }

        virtual void on_frame_ready(astra::StreamReader& reader, astra::Frame& frame) override;

        virtual void on_set_parameter(astra::plugins::stream* stream,
                                      astra_streamconnection_t connection,
                                      astra_parameter_id id,
                                      size_t inByteLength,
                                      astra_parameter_data_t inData) override;

        virtual void on_get_parameter(astra::plugins::stream* stream,
                                      astra_streamconnection_t connection,
                                      astra_parameter_id id,
                                      astra_parameter_bin_t& parameterBin) override;

    private:
        astra_stream_t sourceStreamHandle_;
        astra::DepthStream depthStream_{nullptr};
        astra::StreamSet sensor_;
        astra::StreamReader reader_;
        astra::PluginServiceProxy& pluginService_;

        std::uint16_t zMin_{0};
        std::uint16_t zMax_{65535};

        using skeletonstream_ptr = std::unique_ptr<skeletonstream>;
        skeletonstream_ptr skeletonStream_;
    };
}}

#endif /* ORBBEC_SKELETON_TRACKER_HPP */
