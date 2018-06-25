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
#ifndef HND_HAND_TRACKER_H
#define HND_HAND_TRACKER_H

#include <astra/astra.hpp>
#include <astra/capi/streams/stream_types.h>
#include <astra_core/plugins/Plugin.hpp>

#include "hnd_depth_utility.hpp"
#include "hnd_tracked_point.hpp"
#include "hnd_point_processor.hpp"
#include "hnd_scaling_coordinate_mapper.hpp"
#include "hnd_handstream.hpp"
#include "hnd_debug_handstream.hpp"
#include "hnd_debug_visualizer.hpp"
#include "hnd_settings.hpp"
#include <memory>
#include "hnd_bitmap.hpp"

namespace astra { namespace hand {

    class hand_tracker : public FrameListener
    {
    public:
        hand_tracker(PluginServiceProxy& pluginService,
                     astra_streamset_t streamSet,
                     StreamDescription& depthDesc,
                     hand_settings& settings);

        virtual ~hand_tracker();
        virtual void on_frame_ready(StreamReader& reader, Frame& frame) override;
    private:
        void create_streams(PluginServiceProxy& pluginService, astra_streamset_t streamSet);
        void reset();
        void generate_hand_frame(astra_frame_index_t frameIndex);
        static void copy_position(Vector3f& source, astra_vector3f_t& target);
        static astra_handstatus_t convert_hand_status(tracking_status status, tracked_point_type type);
        static void reset_hand_point(astra_handpoint_t& point);

        void overlay_circle(_astra_imageframe& imageFrame);
        void update_debug_image_frame(_astra_imageframe& astraColorframe);
        void generate_hand_debug_image_frame(astra_frame_index_t frameIndex);
        void update_tracking(const DepthFrame& depthFrame, const PointFrame& pointFrame);
        void update_hand_frame(std::vector<tracked_point>& internaltracked_points, _astra_handframe& frame);

        void debug_probe_point(tracking_matrices& matrices);
        void debug_spawn_point(tracking_matrices& matrices);

        void track_points(BitmapF& matDepth,
                          BitmapF& matDepthFullSize,
                          BitmapMask& matForeground,
                          const astra::Vector3f* worldPoints);
        Point2i get_mouse_probe_position();
        Point2i get_spawn_position();

        //fields

        StreamSet streamset_;
        StreamReader reader_;
        DepthStream depthStream_;

        hand_settings& settings_;
        PluginServiceProxy& pluginService_;
        depth_utility depthUtility_;
        point_processor pointProcessor_;

        float processingSizeWidth_;
        float processingSizeHeight_;

        using ColorStreamPtr = std::unique_ptr<debug_handstream>;
        ColorStreamPtr debugimagestream_;

        using handstream_ptr = std::unique_ptr<handstream>;
        handstream_ptr handStream_;

        BitmapF matDepth_;
        BitmapF matDepthFullSize_;
        BitmapF matDepthWindow_;
        BitmapMask matVelocitySignal_;
        BitmapF matArea_;
        BitmapF matAreaSqrt_;
        BitmapF layerIntegralArea_;
        BitmapMask debugUpdateSegmentation_;
        BitmapMask debugCreateSegmentation_;
        BitmapMask debugRefineSegmentation_;
        BitmapMask updateForegroundSearched_;
        BitmapMask createForegroundSearched_;
        BitmapF debugUpdateScore_;
        BitmapF debugCreateScore_;
        BitmapF debugRefineScore_;
        BitmapF debugUpdateScoreValue_;
        BitmapF debugCreateScoreValue_;
        BitmapF debugRefineScoreValue_;
        BitmapMask debugCreateTestPassMap_;
        BitmapMask debugUpdateTestPassMap_;
        BitmapMask debugRefineTestPassMap_;

        BitmapMask layerSegmentation_;
        BitmapF layerScore_;
        BitmapF layerEdgeDistance_;
        BitmapMask layerTestPassMap_;

        BitmapMask refineForegroundSearched_;
        BitmapMask refineSegmentation_;
        BitmapF refineScore_;
        BitmapF refineEdgeDistance_;

        astra::Vector3f* worldPoints_{nullptr};
        int numWorldPoints_{0};

        debug_visualizer debugVisualizer_;
    };

}}

#endif // HND_HAND_TRACKER_H
