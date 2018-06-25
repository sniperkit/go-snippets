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
#ifndef HND_DEPTH_UTILITY_H
#define HND_DEPTH_UTILITY_H

#include <astra/astra.hpp>
#include "hnd_settings.hpp"
#include "hnd_bitmap.hpp"
#include <cstdint>

#ifndef MIN
#define MIN(a,b) (((a)<(b))?(a):(b))
#endif

#ifndef MAX
#define MAX(a,b) (((a)>(b))?(a):(b))
#endif

namespace astra { namespace hand {

    class depth_utility
    {
    public:
        depth_utility(float width, float height, depth_utility_settings& settings);
        virtual ~depth_utility();

        void depth_to_velocity_signal(const DepthFrame& depthFrame,
                                      BitmapF& matDepth,
                                      BitmapF& matDepthFullSize,
                                      BitmapMask& matVelocitySignal);
        void reset();

        const BitmapF& matDepthVel() const { return matDepthVel_; }
        const BitmapF& matDepthAvg() const { return matDepthAvg_; }
        const BitmapF& matDepthVelErode() const { return matDepthVelErode_; }
        const BitmapF& matDepthFilled() const { return matDepthFilled_; }

    private:
        enum class fill_mask_type : std::uint8_t
        {
            normal = 0,
            filled = 1
        };

        static void depthframe_to_matrix(const DepthFrame& depthFrameSrc,
                                         const int width,
                                         const int height,
                                         BitmapF& matTarget);

        static void accumulate_averages(const BitmapF& input,
                                        const BitmapF& prevInput,
                                        const BitmapMask& mask,
                                        const float alpha,
                                        const float jumpThreshold,
                                        BitmapF& accumulated);

        static void fill_zero_values(BitmapF& matDepth,
                                     BitmapF& matDepthFilled,
                                     BitmapMask& matDepthFilledMask,
                                     BitmapF& matDepthPrevious);

        static void filter_zero_values_and_jumps(BitmapF& depthCurrent,
                                                 BitmapF& depthPrev,
                                                 BitmapF& depthAvg,
                                                 BitmapMask& matDepthFilledMask,
                                                 const float maxDepthJumpPercent);
        void threshold_velocity_signal(BitmapF& matVelocityFiltered,
                                       BitmapMask& matVelocitySignal,
                                       const float velocityThresholdFactor);

        void adjust_velocities_for_depth(BitmapF& matDepth,
                                         BitmapF& matVelocityFiltered);

        int depth_to_chunk_index(float depth);

        void analyze_velocities(BitmapF& matDepth,
                                BitmapF& matVelocityFiltered);

        const float processingWidth_;
        const float processingHeight_;

        BitmapMask rectElement_;
        BitmapMask rectElement2_;
        BitmapF matDepthOriginal_;
        BitmapF matDepthPrevious_;
        BitmapF matDepthFilled_;
        BitmapMask matDepthFilledMask_;
        BitmapF matDepthAvg_;
        BitmapF matDepthVel_;
        BitmapF matDepthVelErode_;

        float depthSmoothingFactor_;
        float velocityThresholdFactor_;
        float maxDepthJumpPercent_;
        int erodeSize_;
        float depthAdjustmentFactor_;
        float minDepth_;
        float maxDepth_;

        static const int NUM_DEPTH_VEL_CHUNKS = 8;
        const float MIN_CHUNK_DEPTH = 0.0f;
        const float MAX_CHUNK_DEPTH = 8000.0f;
        float velErodeFactor_{ 0.98f };
        float maxVel_[NUM_DEPTH_VEL_CHUNKS];
        float depthCount_[NUM_DEPTH_VEL_CHUNKS];
    };
}}

#endif // HND_DEPTH_UTILITY_H
