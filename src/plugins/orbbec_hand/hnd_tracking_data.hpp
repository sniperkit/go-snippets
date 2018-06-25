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
#ifndef HND_TRACKING_DATA_H
#define HND_TRACKING_DATA_H

#include "hnd_bitmap.hpp"
#include "hnd_scaling_coordinate_mapper.hpp"
#include "hnd_settings.hpp"
#include <cstdint>

namespace astra { namespace hand {

    enum pixel_type
    {
        background = 0,
        foreground = 1,
        searched = 2,
        searched_from_out_of_range = 3,
        foreground_natural_edge = 4,
        foreground_out_of_range_edge = 5
    };

    enum class tracked_point_type
    {
        candidate_point,
        active_point
    };

    enum class tracking_status
    {
        not_tracking,
        tracking,
        lost,
        dead
    };

    enum test_behavior
    {
        TEST_BEHAVIOR_NONE = 0,
        TEST_BEHAVIOR_LOG = 1
    };

    enum test_phase
    {
        TEST_PHASE_CREATE = 0,
        TEST_PHASE_UPDATE = 1
    };

    inline std::string tracking_status_to_string(tracking_status status)
    {
        switch (status)
        {
        case tracking_status::not_tracking:
            return "Not Tracking";
        case tracking_status::tracking:
            return "Tracking";
        case tracking_status::lost:
            return "Lost";
        case tracking_status::dead:
            return "Dead";
        default:
            return "Unknown";
        }
    }

    enum segmentation_velocity_policy
    {
        VELOCITY_POLICY_IGNORE = 0,
        VELOCITY_POLICY_RESET_TTL = 1
    };

    struct tracking_matrices
    {
        BitmapF& depthFullSize;
        BitmapF& depth;
        BitmapF& area;
        BitmapF& areaSqrt;
        BitmapMask& velocitySignal;
        BitmapMask& layerSegmentation;
        BitmapF& layerScore;
        BitmapF& layerEdgeDistance;
        BitmapF& layerIntegralArea;
        BitmapMask& layerTestPassMap;
        BitmapMask& foregroundSearched;
        BitmapMask& debugSegmentation;
        BitmapF& debugScore;
        BitmapF& debugScoreValue;
        BitmapMask& debugTestPassMap;
        bool enableTestPassMap;
        const astra::Vector3f* fullSizeWorldPoints;
        astra::Vector3f* worldPoints;
        bool debugLayersEnabled;
        int layerCount;
        const astra::CoordinateMapper& fullSizeMapper;
        const conversion_cache_t depthToWorldData;
        std::vector<astra::Vector2i> layerCirclePoints;

        tracking_matrices(BitmapF& depthFullSize,
                          BitmapF& depth,
                          BitmapF& area,
                          BitmapF& areaSqrt,
                          BitmapMask& velocitySignal,
                          BitmapMask& foregroundSearched,
                          BitmapMask& layerSegmentation,
                          BitmapF& layerScore,
                          BitmapF& layerEdgeDistance,
                          BitmapF& layerIntegralArea,
                          BitmapMask& layerTestPassMap,
                          BitmapMask& debugSegmentation,
                          BitmapF& debugScore,
                          BitmapF& debugScoreValue,
                          BitmapMask& debugTestPassMap,
                          bool enableTestPassMap,
                          const astra::Vector3f* fullSizeWorldPoints,
                          astra::Vector3f* worldPoints,
                          bool debugLayersEnabled,
                          const astra::CoordinateMapper& fullSizeMapper,
                          const conversion_cache_t depthToWorldData)
        :
            depthFullSize(depthFullSize),
            depth(depth),
            area(area),
            areaSqrt(areaSqrt),
            velocitySignal(velocitySignal),
            layerSegmentation(layerSegmentation),
            layerScore(layerScore),
            layerEdgeDistance(layerEdgeDistance),
            layerIntegralArea(layerIntegralArea),
            layerTestPassMap(layerTestPassMap),
            foregroundSearched(foregroundSearched),
            debugSegmentation(debugSegmentation),
            debugScore(debugScore),
            debugScoreValue(debugScoreValue),
            debugTestPassMap(debugTestPassMap),
            enableTestPassMap(enableTestPassMap),
            fullSizeWorldPoints(fullSizeWorldPoints),
            worldPoints(worldPoints),
            debugLayersEnabled(debugLayersEnabled),
            layerCount(0),
            fullSizeMapper(fullSizeMapper),
            depthToWorldData(depthToWorldData)
        { }
    };

    inline float get_resize_factor(tracking_matrices& matrices)
    {
        float resizeFactor = matrices.depthFullSize.width() / static_cast<float>(matrices.depth.width());

        return resizeFactor;
    }

    inline scaling_coordinate_mapper get_scaling_mapper(tracking_matrices& matrices)
    {
        const float resizeFactor = get_resize_factor(matrices);

        return scaling_coordinate_mapper(matrices.depthToWorldData, resizeFactor);
    }

    struct tracking_data
    {
        tracking_matrices& matrices;
        const Point2i& seedPosition;
        const Vector3f referenceWorldPosition;
        const float referenceAreaSqrt;
        const segmentation_velocity_policy velocityPolicy;
        const segmentation_settings settings;
        const test_phase phase;

        tracking_data(tracking_matrices& matrices,
                      const Point2i& seedPosition,
                      const Vector3f referenceWorldPosition,
                      const float referenceAreaSqrt,
                      const segmentation_velocity_policy velocityPolicy,
                      const segmentation_settings settings,
                      const test_phase phase)
            : matrices(matrices),
              seedPosition(seedPosition),
              referenceWorldPosition(referenceWorldPosition),
              referenceAreaSqrt(referenceAreaSqrt),
              velocityPolicy(velocityPolicy),
              settings(settings),
              phase(phase)
        {}
    };
}}

#endif // HND_TRACKING_DATA_H
