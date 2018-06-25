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
#include "hnd_trajectory_analyzer.hpp"
#include "hnd_constants.hpp"
#include <astra_core/plugins/PluginLogging.hpp>

namespace astra { namespace hand {

    using namespace std;

    trajectory_analyzer::trajectory_analyzer(int trackingId, trajectory_analyzer_settings& settings) :
        trackingId_(trackingId),
        pointSteady_(false),
        numSteadyFrames_(0),
        accumulatedDeltaHeading_(),
        lastAccumulatedDeltaHeading_(),
        avgDeltaHeadingValid_(false),
        lastAvgDeltaHeadingValid_(false),
        isTrackingHeading_(false),
        headingTrackStart_(),
        isInflecting_(false),
        numWaveInflections_(0),
        isWaveGesture_(false),
        framesSinceInflection_(0),
        recentDeltaHeading_(),
        /* settings below */
        maxSteadyDelta_(settings.maxSteadyDelta),
        minSteadyFrames_(settings.minSteadyFrames),
        minHeadingDist_(settings.minHeadingDist),
        deltaHeadingFactor_(settings.deltaHeadingFactor),
        minHeadingDiffForInflection_(settings.minHeadingDiffForInflection),
        maxHeadingDiffForContinuation_(settings.maxHeadingDiffForContinuation),
        minWaveInflectionsForGesture_(settings.minWaveInflectionsForGesture),
        maxFramesBetweenInflections_(settings.maxFramesBetweenInflections)
    {
    }

    trajectory_analyzer::~trajectory_analyzer()
    {
    }

    void trajectory_analyzer::reset_steady()
    {
        pointSteady_ = false;
        numSteadyFrames_ = 0;
    }

    void trajectory_analyzer::reset_wave()
    {
        LOG_INFO("trajectory_analyzer", "Reset wave gesture for point #%d", trackingId_);

        isWaveGesture_ = false;
        isInflecting_ = false;
        lastAccumulatedDeltaHeading_ = Vector3f();
        lastAvgDeltaHeadingValid_ = false;
        avgDeltaHeadingValid_ = false;
        isTrackingHeading_ = false;
        numWaveInflections_ = 0;
        recentDeltaHeading_ = Vector3f();
    }

    void trajectory_analyzer::set_for_next_inflection()
    {
        lastAccumulatedDeltaHeading_ = accumulatedDeltaHeading_;
        isInflecting_ = false;
        avgDeltaHeadingValid_ = false;
        lastAvgDeltaHeadingValid_ = true;
        isTrackingHeading_ = false;
    }

    void trajectory_analyzer::update(tracked_point& point)
    {
        if (point.trackingId != trackingId_)
        {
            throw std::logic_error("trajectory_analyzer updated with wrong tracking id");
        }
        ++framesSinceInflection_;
        if (framesSinceInflection_ == maxFramesBetweenInflections_)
        {
            LOG_INFO("trajectory_analyzer", "Wave gesture timed out for point #%d", trackingId_);

            reset_wave();
        }

        Vector3f deltaPosition = point.fullSizeWorldDeltaPosition;

        float delta = static_cast<float>(deltaPosition.length());

        if (delta > maxSteadyDelta_)
        {
            Vector3f deltaPositionNullY = deltaPosition;
            deltaPositionNullY.y = 0;

            if (!isTrackingHeading_)
            {
                isTrackingHeading_ = true;
                headingTrackStart_ = point.fullSizeWorldPosition;
                accumulatedDeltaHeading_ = deltaPositionNullY;
            }
            else
            {
                accumulatedDeltaHeading_ += deltaPositionNullY;

                recentDeltaHeading_ = deltaPositionNullY * (1 - deltaHeadingFactor_) + deltaPositionNullY * deltaHeadingFactor_;

                avgDeltaHeadingValid_ = is_valid_heading_dist(point.fullSizeWorldPosition);

                float headingDist = static_cast<float>((point.fullSizeWorldPosition - headingTrackStart_).length());

                LOG_TRACE("trajectory_analyzer", "#%d dist %f v1: %d v2: %d", trackingId_, headingDist, avgDeltaHeadingValid_, lastAvgDeltaHeadingValid_);

                if (avgDeltaHeadingValid_ && lastAvgDeltaHeadingValid_)
                {
                    float degreeDifference = get_degree_difference(accumulatedDeltaHeading_, lastAccumulatedDeltaHeading_);

                    if (degreeDifference > minHeadingDiffForInflection_)
                    {
                        if (!isInflecting_)
                        {
                            isInflecting_ = true;
                            framesSinceInflection_ = 0;
                            ++numWaveInflections_;
                            if (!isWaveGesture_)
                            {
                                LOG_INFO("trajectory_analyzer", "Wave count %d for point #%d", numWaveInflections_, trackingId_);
                                if (numWaveInflections_ == minWaveInflectionsForGesture_)
                                {
                                    LOG_INFO("trajectory_analyzer", "Wave gesture detected for point #%d", trackingId_);
                                    isWaveGesture_ = true;
                                }
                            }
                        }
                    }
                    else if (degreeDifference > maxHeadingDiffForContinuation_)
                    {
                        reset_wave();
                    }
                }
                if (avgDeltaHeadingValid_)
                {
                    float recentDegreeDifference = get_degree_difference(accumulatedDeltaHeading_, recentDeltaHeading_);

                    if (recentDegreeDifference > minHeadingDiffForInflection_)
                    {
                        set_for_next_inflection();
                    }
                }
            }
        }
        else if (avgDeltaHeadingValid_)
        {
            set_for_next_inflection();
        }

        if (delta < maxSteadyDelta_)
        {
            ++numSteadyFrames_;
            if (!pointSteady_ && numSteadyFrames_ > minSteadyFrames_)
            {
                pointSteady_ = true;

                LOG_INFO("trajectory_analyzer", "Steady gesture detected for point #%d", trackingId_);

                if (isWaveGesture_)
                {
                    reset_wave();
                }
            }
        }
        else
        {
            numSteadyFrames_ = 0;
            pointSteady_ = false;
        }
    }

    float trajectory_analyzer::get_delta_angle(float x, float y)
    {
        float radians = std::atan2(y, x);
        float degrees = radians * RAD_TO_DEG;
        return degrees;
    }

    bool trajectory_analyzer::is_valid_heading_dist(const Vector3f& currentWorldPosition)
    {
        float headingDist = static_cast<float>((currentWorldPosition - headingTrackStart_).length());
        bool validDist = headingDist > minHeadingDist_;
        return validDist;
    }

    float trajectory_analyzer::get_degree_difference(Vector3f& v1, Vector3f& v2)
    {
        float len1 = static_cast<float>(v1.length());
        float len2 = static_cast<float>(v2.length());

        if (len1 < EPSILON || len2 < EPSILON)
        {
            return 0.0f;
        }

        float invLen1 = 1.0 / len1;
        float invLen2 = 1.0 / len2;

        auto norm1 = v1 * invLen1;
        auto norm2 = v2 * invLen2;

        float angleBetween = std::acos(norm1.dot(norm2));
        float degreeBetween = angleBetween * RAD_TO_DEG;
        return degreeBetween;
    }
}}
