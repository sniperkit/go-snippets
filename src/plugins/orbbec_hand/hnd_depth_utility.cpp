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
#include "hnd_depth_utility.hpp"
#include <astra/astra.hpp>
#include "hnd_tracking_data.hpp"
#include "hnd_morphology.hpp"
#include <cmath>
#include <Shiny.h>

namespace astra { namespace hand {

    depth_utility::depth_utility(float width, float height, depth_utility_settings& settings) :
        processingWidth_(width),
        processingHeight_(height),
        depthSmoothingFactor_(settings.depthSmoothingFactor),
        velocityThresholdFactor_(settings.velocityThresholdFactor),
        maxDepthJumpPercent_(settings.maxDepthJumpPercent),
        erodeSize_(settings.erodeSize),
        depthAdjustmentFactor_(settings.depthAdjustmentFactor),
        minDepth_(settings.minDepth),
        maxDepth_(settings.maxDepth)
    {
        PROFILE_FUNC();
        rectElement_ = get_structuring_element(MorphShape::Rect,
                                               Size2i(erodeSize_ * 2 + 1, erodeSize_ * 2 + 1));

        reset();
    }

    depth_utility::~depth_utility()
    {
        PROFILE_FUNC();
    }

    void depth_utility::reset()
    {
        PROFILE_FUNC();

        Size2i size(processingWidth_, processingHeight_);

        matDepthFilled_.recreate(size);
        matDepthFilled_.fill(0.f);

        matDepthFilledMask_.recreate(size);
        matDepthFilledMask_.fill(0);

        matDepthPrevious_.recreate(size);
        matDepthPrevious_.fill(0.f);

        matDepthAvg_.recreate(size);
        matDepthAvg_.fill(0.f);

        matDepthVel_.recreate(size);
        matDepthVel_.fill(0.f);

        matDepthVelErode_.recreate(size);
        matDepthVelErode_.fill(0.f);

        for (int i = 0; i < NUM_DEPTH_VEL_CHUNKS; i++)
        {
            maxVel_[i] = 0;
        }
    }

    void depth_utility::accumulate_averages(const BitmapF& input,
                                            const BitmapF& prevInput,
                                            const BitmapMask& mask,
                                            const float alpha,
                                            const float jumpThreshold,
                                            BitmapF& accumulated)
    {
        assert(input.width() == accumulated.width());
        assert(input.height() == accumulated.height());

        const float oneMinusAlpha = 1.0f - alpha;

        const unsigned width = input.width();
        const unsigned height = input.height();

        const float* inputData = input.data();
        const auto* prevInputData = prevInput.data();
        const MaskType* maskData = mask.data();
        float* accumulatedData = accumulated.data();

        for(unsigned y = 0; y < height; ++y)
        {
            for (unsigned x = 0; x < width; ++x, ++inputData, ++prevInputData, ++maskData, ++accumulatedData)
            {
                //update running average
                *accumulatedData = oneMinusAlpha * (*accumulatedData) +
                                   alpha * (*inputData);

                const fill_mask_type mergeMask = fill_mask_type(*maskData);
                const bool updatedValue = mergeMask == fill_mask_type::normal;
                const bool indeterminateValues = *inputData == 0.0f || *prevInputData == 0.0f;

                if (!updatedValue || indeterminateValues)
                {
                    *accumulatedData = *inputData;
                }
                else
                {
                    const float currentValue = *inputData;
                    const float previousValue = *prevInputData;
                    const float percentChange = (currentValue - previousValue) / previousValue;
                    const bool isIncreasing = percentChange > 0.0f;

                    if (isIncreasing && std::fabs(percentChange) > jumpThreshold)
                    {
                        *accumulatedData = *inputData;
                    }
                }
            }
        }
    }

    void depth_utility::depth_to_velocity_signal(const DepthFrame& depthFrame,
                                                 BitmapF& matDepth,
                                                 BitmapF& matDepthFullSize,
                                                 BitmapMask& matVelocitySignal)
    {
        PROFILE_FUNC();
        const int width = depthFrame.width();
        const int height = depthFrame.height();

        depthframe_to_matrix(depthFrame, width, height, matDepthFullSize);

        if (width == processingWidth_ && height == processingHeight_)
        {
            //target size is original size, just use the same data
            matDepth = matDepthFullSize;
        }
        else
        {
            matDepth.recreate(processingWidth_, processingHeight_);
            //convert to the target processing size with nearest neighbor
            resize(matDepthFullSize, matDepth);
        }

        matVelocitySignal.recreate(matDepth.size());
        matVelocitySignal.fill(0);

        //fill 0 depth pixels with the value from the previous frame
        fill_zero_values(matDepth, matDepthFilled_, matDepthFilledMask_, matDepthPrevious_);

        accumulate_averages(matDepthFilled_,
                            matDepthPrevious_,
                            matDepthFilledMask_,
                            depthSmoothingFactor_,
                            maxDepthJumpPercent_,
                            matDepthAvg_);

        //current minus average, scaled by average = velocity as a percent change

        auto* velocityData = matDepthVel_.data();
        const float* t0Data = matDepthAvg_.data();
        const float* t1Data = matDepthFilled_.data();

        // abstractly:
        //matDepthVel_ = (matDepthFilled_ - matDepthAvg_) / matDepthAvg_;
        for (unsigned y = 0; y < matDepthVel_.height(); ++y)
        {
            for(unsigned x = 0; x < matDepthVel_.width(); ++x, ++velocityData, ++t0Data, ++t1Data)
            {
                float t0Value = *t0Data + std::numeric_limits<float>::epsilon();
                float t1Value = *t1Data;
                // percentage of the previous velocities
                float percentVelocity = ((t1Value - t0Value) + std::numeric_limits<float>::epsilon()) / t0Value;

                *velocityData = percentVelocity;
            }
        }

        // scales velocity by depth, removed signed velocities
        adjust_velocities_for_depth(matDepth, matDepthVel_);

        erode(matDepthVel_, matDepthVelErode_, rectElement_);

        threshold_velocity_signal(matDepthVelErode_,
                                  matVelocitySignal,
                                  velocityThresholdFactor_);

        copy_to(matDepth, matDepthPrevious_);

        //analyze_velocities(matDepth, matDepthVelErode_);
    }

    void depth_utility::depthframe_to_matrix(const DepthFrame& depthFrameSrc,
                                             const int width,
                                             const int height,
                                             BitmapF& matTarget)
    {
        PROFILE_FUNC();
        //ensure initialized
        matTarget.recreate(width, height);

        const int16_t* depthData = depthFrameSrc.data();

        for (int y = 0; y < height; ++y)
        {
            float* row = matTarget.data(y);
            for (int x = 0; x < width; ++x)
            {
                float depth = static_cast<float>(*depthData);
                *row = depth;
                ++row;
                ++depthData;
            }
        }
    }

    void depth_utility::fill_zero_values(BitmapF& matDepth,
                                         BitmapF& matDepthFilled,
                                         BitmapMask& matDepthFilledMask,
                                         BitmapF& matDepthPrevious)
    {
        PROFILE_FUNC();
        const int width = matDepth.width();
        const int height = matDepth.height();

        for (int y = 0; y < height; ++y)
        {
            const auto* depthRow = matDepth.data(y);
            const auto* prevDepthRow = matDepthPrevious.data(y);
            auto* filledDepthRow = matDepthFilled.data(y);
            auto* filledDepthMaskRow = matDepthFilledMask.data(y);

            for (int x = 0; x < width; ++x)
            {
                auto depth = *depthRow;

                if (depth == 0)
                {
                    depth = *prevDepthRow;
                    *filledDepthMaskRow = static_cast<MaskType>(fill_mask_type::filled);
                }
                else
                {
                    *filledDepthMaskRow = static_cast<MaskType>(fill_mask_type::normal);
                }

                *filledDepthRow = depth;

                ++depthRow;
                ++prevDepthRow;
                ++filledDepthRow;
                ++filledDepthMaskRow;
            }
        }
    }

    void depth_utility::filter_zero_values_and_jumps(BitmapF& matDepth,
                                                     BitmapF& matDepthPrevious,
                                                     BitmapF& matDepthAvg,
                                                     BitmapMask& matDepthFilledMask,
                                                     const float maxDepthJumpPercent)
    {
        PROFILE_FUNC();
        const int width = matDepth.width();
        const int height = matDepth.width();

        for (int y = 0; y < height; ++y)
        {
            auto* depthRow = matDepth.data(y);
            auto* prevDepthRow = matDepthPrevious.data(y);
            auto* avgRow = matDepthAvg.data(y);
            auto* filledDepthMaskRow = matDepthFilledMask.data(y);

            for (int x = 0; x < width; ++x, ++depthRow, ++prevDepthRow, ++filledDepthMaskRow)
            {
                float depth = *depthRow;
                float previousDepth = *prevDepthRow;
                fill_mask_type fillType = static_cast<fill_mask_type>(*filledDepthMaskRow);

                //calculate percent change since last frame
                float deltaPercent = (depth - previousDepth) / previousDepth;
                float absDeltaPercent = std::fabs(deltaPercent);
                bool movingAway = deltaPercent > 0;

                //suppress signal if either current or previous pixel are invalid
                bool isZeroDepth = (0 == depth || 0 == previousDepth);

                //suppress signal when a pixel was artificially filled
                bool isFilled = (fillType == fill_mask_type::filled);

                //suppress signal when a pixel jumps a long distance from near to far
                bool isJumpingAway = (absDeltaPercent > maxDepthJumpPercent && movingAway);

                if (isZeroDepth || isFilled || isJumpingAway)
                {
                    //set the average to the current depth, and set velocity to zero
                    //this suppresses the velocity signal for edge jumping artifacts
                    avgRow[x] = depth;
                }

                *prevDepthRow = depth;
            }
        }
    }

    void depth_utility::threshold_velocity_signal(BitmapF& matVelocityFiltered,
                                                  BitmapMask& matVelocitySignal,
                                                  const float velocityThresholdFactor)
    {
        PROFILE_FUNC();
        const int width = matVelocitySignal.width();
        const int height = matVelocitySignal.height();

        for (int y = 0; y < height; ++y)
        {
            auto* velFilteredRow = matVelocityFiltered.data(y);
            auto* velocitySignalRow = matVelocitySignal.data(y);

            for (int x = 0; x < width; ++x, ++velFilteredRow, ++velocitySignalRow)
            {
                //matVelocityFiltered is already abs(vel)
                const float velFiltered = *velFilteredRow;

                if (velFiltered > velocityThresholdFactor)
                {
                    *velocitySignalRow = pixel_type::foreground;
                }
                else
                {
                    *velocitySignalRow = pixel_type::background;
                }
            }
        }
    }

    void depth_utility::adjust_velocities_for_depth(BitmapF& matDepth, BitmapF& matVelocityFiltered)
    {
        PROFILE_FUNC();
        if (depthAdjustmentFactor_ == 0)
        {
            return;
        }

        const int width = matDepth.width();
        const int height = matDepth.height();

        for (int y = 0; y < height; ++y)
        {
            auto* depthRow = matDepth.data(y);
            auto* velFilteredRow = matVelocityFiltered.data(y);

            for (int x = 0; x < width; ++x, ++depthRow, ++velFilteredRow)
            {
                const float depth = *depthRow;
                if (depth != 0.0f)
                {
                    float& velFiltered = *velFilteredRow;

                    if (depth > minDepth_ && depth < maxDepth_)
                    {
                        float depthM = depth / 1000.0f;
                        // scale by depth, constant, remove signed velocities
                        velFiltered = std::fabs(velFiltered / (depthM * depthAdjustmentFactor_));
                    }
                    else
                    {
                        velFiltered = 0;
                    }
                }
            }
        }
    }

    int depth_utility::depth_to_chunk_index(float depth)
    {
        PROFILE_FUNC();
        if (depth == 0 || depth < MIN_CHUNK_DEPTH || depth > MAX_CHUNK_DEPTH)
        {
            return -1;
        }
        const float chunkRange = MAX_CHUNK_DEPTH - MIN_CHUNK_DEPTH;
        const float normDepth = (depth - MIN_CHUNK_DEPTH) / chunkRange;

        return std::min(NUM_DEPTH_VEL_CHUNKS - 1, static_cast<int>(NUM_DEPTH_VEL_CHUNKS * normDepth));
    }

    void depth_utility::analyze_velocities(BitmapF& matDepth, BitmapF& matVelocityFiltered)
    {
        PROFILE_FUNC();
        int width = matDepth.width();
        int height = matDepth.height();

        for (int i = 0; i < NUM_DEPTH_VEL_CHUNKS; i++)
        {
            maxVel_[i] *= velErodeFactor_;
            depthCount_[i] = 0;
        }

        for (int y = 0; y < height; ++y)
        {
            float* depthRow = matDepth.data(y);
            float* velFilteredRow = matVelocityFiltered.data(y);

            for (int x = 0; x < width; ++x, ++depthRow, ++velFilteredRow)
            {
                //matVelocityFiltered is already abs(vel)

                float depth = *depthRow;
                int chunkIndex = depth_to_chunk_index(depth);
                if (chunkIndex >= 0 && depth != 0)
                {
                    ++depthCount_[chunkIndex];

                    float& maxVel = maxVel_[chunkIndex];
                    float velFiltered = *velFilteredRow;

                    if (velFiltered > maxVel)
                    {
                        maxVel = velFiltered;
                    }
                }
            }
        }

        const float chunkRange = MAX_CHUNK_DEPTH - MIN_CHUNK_DEPTH;
        const float chunk_size = chunkRange / NUM_DEPTH_VEL_CHUNKS;

        for (int i = 0; i < NUM_DEPTH_VEL_CHUNKS; i++)
        {
            float maxVel = maxVel_[i];
            int count = depthCount_[i];
            float startDepth = (MIN_CHUNK_DEPTH + i * chunk_size) / 1000.0f;
            float endDepth = (MIN_CHUNK_DEPTH + (1 + i) * chunk_size) / 1000.0f;
            float ratio = 0;
            if (endDepth != 0.0f)
            {
                ratio = maxVel / endDepth;
            }

            printf("[%.1fm %d,%f,%f] ", startDepth, count, maxVel, ratio);
            if (i % 3 == 2)
            {
                printf("\n");
            }
        }
        printf("[%.1fm]\n\n", static_cast<int>(MAX_CHUNK_DEPTH)/1000.0f);
    }
}}
