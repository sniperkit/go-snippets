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
#ifndef HND_TRAJECTORY_ANALYZER_H
#define HND_TRAJECTORY_ANALYZER_H

#include "hnd_tracked_point.hpp"
#include "hnd_settings.hpp"
#include <astra_core/plugins/PluginLogging.hpp>

namespace astra { namespace hand {

    class trajectory_analyzer
    {
    public:
        trajectory_analyzer(int trackingId, trajectory_analyzer_settings& settings);
        ~trajectory_analyzer();

        void update(tracked_point& point);

        void reset_wave();
        void set_for_next_inflection();
        void reset_steady();

        bool pointSteady() const { return pointSteady_; }
        int trackingId() const { return trackingId_; }
        bool is_wave_gesture() const { return isWaveGesture_; }

    private:
        float get_delta_angle(float x, float y);
        bool is_valid_heading_dist(const Vector3f& currentWorldPosition);
        float get_degree_difference(Vector3f& v1, Vector3f& v2);

        int trackingId_;

        bool pointSteady_;
        int numSteadyFrames_;
        Vector3f accumulatedDeltaHeading_;
        Vector3f lastAccumulatedDeltaHeading_;
        bool avgDeltaHeadingValid_;
        bool lastAvgDeltaHeadingValid_;
        bool isTrackingHeading_;
        Vector3f headingTrackStart_;
        bool isInflecting_;
        int numWaveInflections_;
        bool isWaveGesture_;
        int framesSinceInflection_;
        Vector3f recentDeltaHeading_;

        float maxSteadyDelta_;
        int minSteadyFrames_;
        int minHeadingDist_;
        float deltaHeadingFactor_;
        float minHeadingDiffForInflection_;
        float maxHeadingDiffForContinuation_;
        int minWaveInflectionsForGesture_;
        int maxFramesBetweenInflections_;
    };
}}

#endif // HND_TRAJECTORY_ANALYZER_H
