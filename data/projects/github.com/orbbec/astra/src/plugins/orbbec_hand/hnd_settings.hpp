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
#ifndef HND_SETTINGS_H
#define HND_SETTINGS_H

#include <string>

namespace astra { namespace hand {

    struct depth_utility_settings
    {
        float depthSmoothingFactor{ 0.05f };
        float velocityThresholdFactor{ 0.005f };
        float maxDepthJumpPercent{ 0.1f };
        int erodeSize{ 1 };
        float depthAdjustmentFactor{ 2.5f };
        float minDepth { 500.0f };
        float maxDepth { 4000.0f };
    };

    struct trajectory_analyzer_settings
    {
        float maxSteadyDelta{ 5.0f };
        int minSteadyFrames{ 15 };
        int minHeadingDist{ 125 };
        float deltaHeadingFactor{ 0.5f };
        float minHeadingDiffForInflection{ 135.0f };
        float maxHeadingDiffForContinuation{ 45.0f };
        int minWaveInflectionsForGesture{ 2 };
        int maxFramesBetweenInflections{ 90 };
    };

    struct area_test_settings
    {
        float areaBandwidth{ 150.0f }; //mm
        float areaBandwidthDepth{ 100.0f }; //mm
        float minArea{ 3000.0f }; //mm^2
        float maxArea{ 35000.0f }; //mm^2
    };

    struct natural_edge_test_settings
    {
        float naturalEdgeBandwidth{ 150.0f }; //mm
        float minPercentNaturalEdges { 0.8f }; //mm^2
    };

    struct circumference_test_settings
    {
        float foregroundRadius1{ 100.0f };
        float foregroundRadius2{ 150.0f };
        float foregroundRadiusMinPercent1{ 0.08f };
        float foregroundRadiusMinPercent2{ 0.04f };
        float foregroundRadiusMaxPercent1{ 0.35f };
        float foregroundRadiusMaxPercent2{ 0.15f };
    };

    struct segmentation_settings
    {
        area_test_settings areaTestSettings;
        circumference_test_settings circumferenceTestSettings;
        natural_edge_test_settings naturalEdgeTestSettings;

        float segmentationBandwidthDepthNear{ 500.0f };
        float segmentationBandwidthDepthFar{ 100.0f };
        float maxSegmentationDist{ 250.0f }; //mm
        float heightScoreFactor{ 0.5f };
        float depthScoreFactor{ 1000.0f };
        float targetEdgeDistance{ 40.0f }; //mm
        float edgeDistanceScoreFactor{ 100.0f };
        float pointInertiaFactor{ 40.0f };
        float pointInertiaRadius{ 40.0f }; //mm
        float maxDepthToDownscaleTestPass { 1300.0f }; //mm
    };

    struct point_processor_settings
    {
        trajectory_analyzer_settings trajectoryAnalyzerSettings;
        segmentation_settings segmentationSettings;

        float maxMatchDistLostActive{ 500.0f }; //mm
        float maxMatchDistDefault{ 500.0f }; //mm
        float steadyDeadBandRadius{ 75.0f }; //mm
        int maxInactiveFramesForCandidatePoints{ 60 };
        int maxInactiveFramesForLostPoints{ 15 };
        int maxInactiveFramesForActivePoints{ 480 };
        float pointSmoothingFactor{ 0.5f };
        float pointDeadBandSmoothingFactor{ 0.05f };
        float pointSmoothingDeadZone{ 50.0f }; //mm
        int maxFailedTestsInProbation{ 5 };
        int probationFrameCount{ 30 };
        int maxFailedTestsInProbationActivePoints{ 3 };
        float secondChanceMinDistance{ 100.0f };
        float mergePointDistance { 100.0f }; //mm
        int maxhandpointUpdatesPerFrame { 10 };
    };

    struct hand_settings
    {
        int processingSizeWidth{ 160 };
        int processingSizeHeight{ 120 };

        depth_utility_settings depthUtilitySettings;
        point_processor_settings pointProcessorSettings;
    };

    hand_settings parse_settings(std::string path);
}}

#endif // HND_SETTINGS_H
