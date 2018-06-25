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
#include "hnd_tracking_data.hpp"
#include <queue>
#include "hnd_scaling_coordinate_mapper.hpp"
#include "hnd_morphology.hpp"
#include <cmath>
#include "hnd_segmentation.hpp"
#include "hnd_constants.hpp"
#include <Shiny.h>
#include <astra_core/plugins/PluginLogging.hpp>
#include <limits>

#define MAX_DEPTH 10000

namespace astra { namespace hand { namespace segmentation {

    struct point_ttl
    {
        int x;
        int y;
        float ttl;

        point_ttl(int x, int y, float ttl) :
            x(x),
            y(y),
            ttl(ttl)
        { }
    };

    static void enqueue_neighbors(BitmapMask& matVisited,
                                  std::queue<point_ttl>& pointQueue,
                                  const point_ttl& pt)
    {
        PROFILE_FUNC();
        const int& x = pt.x;
        const int& y = pt.y;
        const int width = matVisited.width();
        const int height = matVisited.height();

        if (x < 1 || x > width - 2 ||
            y < 1 || y > height - 2)
        {
            return;
        }

        const float& ttlRef = pt.ttl;

        auto& rightVisited = matVisited.at(x+1, y);
        if (0 == rightVisited)
        {
            rightVisited = 1;
            pointQueue.push(point_ttl(x+1, y, ttlRef));
        }

        auto& leftVisited = matVisited.at(x-1, y);
        if (0 == leftVisited)
        {
            leftVisited = 1;
            pointQueue.push(point_ttl(x-1, y, ttlRef));
        }

        auto& downVisited = matVisited.at(x, y+1);
        if (0 == downVisited)
        {
            downVisited = 1;
            pointQueue.push(point_ttl(x, y+1, ttlRef));
        }

        auto& upVisited = matVisited.at(x, y-1);
        if (0 == upVisited)
        {
            upVisited = 1;
            pointQueue.push(point_ttl(x, y-1, ttlRef));
        }
    }

    static Point2i find_nearest_in_range_pixel(tracking_data& data,
                                               BitmapMask& matVisited)
    {
        PROFILE_FUNC();
        assert(matVisited.size() == data.matrices.depth.size());
        const float referenceAreaSqrt = data.referenceAreaSqrt;
        if (referenceAreaSqrt == 0)
        {
            return INVALID_POINT;
        }

        const float minDepth = data.referenceWorldPosition.z - data.settings.segmentationBandwidthDepthNear;
        const float maxDepth = data.referenceWorldPosition.z + data.settings.segmentationBandwidthDepthFar;
        const float maxSegmentationDist = data.settings.maxSegmentationDist;
        BitmapF& depthMatrix = data.matrices.depth;
        BitmapMask& searchedMatrix = data.matrices.foregroundSearched;

        std::queue<point_ttl> pointQueue;

        pointQueue.push(point_ttl(data.seedPosition.x,
                                  data.seedPosition.y,
                                  maxSegmentationDist));

        matVisited.at(data.seedPosition) = 1;

        while (!pointQueue.empty())
        {
            point_ttl pt = pointQueue.front();
            pointQueue.pop();
            const int& x = pt.x;
            const int& y = pt.y;
            float& ttlRef = pt.ttl;

            if (ttlRef <= 0)
            {
                continue;
            }

            searchedMatrix.at(x, y) =
                pixel_type::searched_from_out_of_range;

            float depth = depthMatrix.at(x, y);
            bool pointInRange = depth != 0 && depth > minDepth && depth < maxDepth;

            if (pointInRange)
            {
                return Point2i(x, y);
            }

            ttlRef -= referenceAreaSqrt;

            enqueue_neighbors(matVisited, pointQueue, pt);
        }

        return INVALID_POINT;
    }

    static float segment_foreground_and_get_average_depth(tracking_data& data)
    {
        PROFILE_FUNC();
        const float& maxSegmentationDist = data.settings.maxSegmentationDist;
        const segmentation_velocity_policy& velocitySignalPolicy = data.velocityPolicy;
        const float seedDepth = data.matrices.depth.at(data.seedPosition);
        const float referenceAreaSqrt = data.referenceAreaSqrt;
        BitmapF& depthMatrix = data.matrices.depth;
        BitmapMask& velocitySignalMatrix = data.matrices.velocitySignal;
        BitmapMask& segmentationMatrix = data.matrices.layerSegmentation;
        BitmapMask& searchedMatrix = data.matrices.foregroundSearched;

        std::queue<point_ttl> pointQueue;

        double totalDepth = 0;
        int depthCount = 0;

        const float bandwidthDepth = data.settings.segmentationBandwidthDepthNear;

        //does the seed point start in range?
        //If not, it will search outward until it finds in range pixels
        const float minDepth = data.referenceWorldPosition.z - bandwidthDepth;
        const float maxDepth = data.referenceWorldPosition.z + data.settings.segmentationBandwidthDepthFar;

        BitmapMask matVisited(depthMatrix.size());
        matVisited.fill(0);

        Point2i seedPosition = data.seedPosition;

        bool seedInRange = seedDepth != 0 && seedDepth > minDepth && seedDepth < maxDepth;
        if (!seedInRange)
        {
            seedPosition = find_nearest_in_range_pixel(data, matVisited);
            if (seedPosition == INVALID_POINT)
            {
                //No in range pixels found, no foreground to set
                return 0.0f;
            }
        }

        pointQueue.push(point_ttl(seedPosition.x,
                                  seedPosition.y,
                                  maxSegmentationDist));

        matVisited.at(seedPosition) = 1;

        while (!pointQueue.empty())
        {
            point_ttl pt = pointQueue.front();
            pointQueue.pop();
            const int& x = pt.x;
            const int& y = pt.y;
            float& ttlRef = pt.ttl;

            if (velocitySignalPolicy == VELOCITY_POLICY_RESET_TTL &&
                velocitySignalMatrix.at(x, y) == pixel_type::foreground)
            {
                ttlRef = maxSegmentationDist;
            }

            float depth = depthMatrix.at(x, y);
            bool pointOutOfRange = depth == 0 ||
                depth < minDepth ||
                        depth > maxDepth;

            if (ttlRef <= 0)
            {
                segmentationMatrix.at(x, y) = pixel_type::foreground_out_of_range_edge;
                continue;
            }
            else if (pointOutOfRange)
            {
                segmentationMatrix.at(x, y) = pixel_type::foreground_natural_edge;
                continue;
            }

            totalDepth += depth;
            ++depthCount;

            searchedMatrix.at(x, y) = pixel_type::searched;
            segmentationMatrix.at(x, y) = pixel_type::foreground;

            ttlRef -= referenceAreaSqrt;

            enqueue_neighbors(matVisited, pointQueue, pt);
        }

        if (depthCount > 0)
        {
            float averageDepth = static_cast<float>(totalDepth / depthCount);
            return averageDepth;
        }
        else
        {
            return 0.0f;
        }
    }

    void calculate_layer_score(tracking_data& data, const float layerAverageDepth)
    {
        PROFILE_FUNC();
        BitmapF& edgeDistanceMatrix = data.matrices.layerEdgeDistance;
        const float depthFactor = data.settings.depthScoreFactor;
        const float heightFactor = data.settings.heightScoreFactor;
        const float edgeDistanceFactor = data.settings.edgeDistanceScoreFactor;
        const float targetEdgeDist = data.settings.targetEdgeDistance;
        BitmapF& layerScoreMatrix = data.matrices.layerScore;
        const float pointInertiaFactor = data.settings.pointInertiaFactor;
        const float pointInertiaRadius = data.settings.pointInertiaRadius;
        const astra::Vector3f* worldPoints = data.matrices.worldPoints;

        layerScoreMatrix.recreate(data.matrices.depth.size());
        layerScoreMatrix.fill(0.f);

        scaling_coordinate_mapper mapper = get_scaling_mapper(data.matrices);

        auto seedWorldPosition = astra::Vector3f(data.referenceWorldPosition.x,
                                                 data.referenceWorldPosition.y,
                                                 data.referenceWorldPosition.z);

        int width = data.matrices.depth.width();
        int height = data.matrices.depth.height();

        int edgeRadius = mapper.scale() * width / 32;
        int minX = edgeRadius - 1;
        int maxX = width - edgeRadius;
        int minY = edgeRadius - 1;
        int maxY = height - edgeRadius;

        for (int y = 0; y < height; y++)
        {
            float* edgeDistanceRow = edgeDistanceMatrix.data(y);
            float* layerScoreRow = layerScoreMatrix.data(y);

            for (int x = 0; x < width; ++x,
                     ++worldPoints,
                     ++edgeDistanceRow,
                     ++layerScoreRow)
            {
                astra::Vector3f worldPosition = *worldPoints;
                if (worldPosition.z != 0 && x > minX && x < maxX && y > minY && y < maxY)
                {
                    //start with arbitrary large value to prevent scores from going negative
                    //(which is ok algorithmically, but debug visuals don't like it)
                    float score = 10000;

                    score += worldPosition.y * heightFactor;

                    float depthDiff = 1.0 - worldPosition.z / layerAverageDepth;
                    score += depthDiff * depthFactor;

                    if (pointInertiaRadius > 0)
                    {
                        auto vector = worldPosition - seedWorldPosition;
                        float length = vector.length();
                        float distFromSeedNorm = std::max(0.0f, std::min(1.0f,
                                                                         length / pointInertiaRadius));
                        score += (1.0f - distFromSeedNorm) * pointInertiaFactor;
                    }

                    float edgeDistance = *edgeDistanceRow;
                    float edgeScore = edgeDistanceFactor * std::min(edgeDistance, targetEdgeDist) / targetEdgeDist;

                    score += edgeScore;

                    *layerScoreRow = score;
                }
                else
                {
                    *layerScoreRow = 0;
                }
            }
        }
    }


    bool test_point_in_range(tracking_matrices& matrices,
                             const Point2i& targetPoint,
                             test_behavior outputLog)
    {
        PROFILE_FUNC();
        if (targetPoint == segmentation::INVALID_POINT ||
            targetPoint.x < 0 || targetPoint.x >= matrices.depth.width() ||
            targetPoint.y < 0 || targetPoint.y >= matrices.depth.height())
        {
            if (outputLog == TEST_BEHAVIOR_LOG)
            {
                LOG_INFO("point_processor", "test_point_in_range failed: position: (%d, %d)",
                         targetPoint.x,
                         targetPoint.y);
            }
            return false;
        }

        if (outputLog == TEST_BEHAVIOR_LOG)
        {
            LOG_INFO("point_processor", "test_point_in_range success: position: (%d, %d)",
                     targetPoint.x,
                     targetPoint.y);
        }

        return true;
    }

    float get_point_area(tracking_matrices& matrices,
                         area_test_settings& settings,
                         const Point2i& point)
    {
        PROFILE_FUNC();
        auto scalingMapper = get_scaling_mapper(matrices);

        float area = count_neighborhood_area(matrices.layerSegmentation,
                                             matrices.depth,
                                             matrices.area,
                                             point,
                                             settings.areaBandwidth,
                                             settings.areaBandwidthDepth,
                                             scalingMapper);

        return area;
    }


    float get_point_area_integral(tracking_matrices& matrices,
                                  BitmapF& integralArea,
                                  area_test_settings& settings,
                                  const Point2i& point)
    {
        //PROFILE_FUNC();
        auto scalingMapper = get_scaling_mapper(matrices);

        float area = count_neighborhood_area_integral(matrices.depth,
                                                      integralArea,
                                                      point,
                                                      settings.areaBandwidth,
                                                      scalingMapper);

        return area;
    }

    bool test_point_area_core(float area,
                              area_test_settings& settings,
                              test_phase phase,
                              test_behavior outputLog)
    {
        //PROFILE_FUNC();
        float minArea = settings.minArea;
        const float maxArea = settings.maxArea;
        if (phase == TEST_PHASE_UPDATE)
        {
            //minimum of 0 during update phase
            minArea = 0;
        }

        bool validPointArea = area > minArea && area < maxArea;

        if (outputLog == TEST_BEHAVIOR_LOG)
        {
            if (validPointArea)
            {
                LOG_INFO("Segmentation", "test_point_area passed: area %f within [%f, %f]",
                         area,
                         minArea,
                         maxArea);
            }
            else
            {
                LOG_INFO("Segmentation", "test_point_area failed: area %f not within [%f, %f]",
                         area,
                         minArea,
                         maxArea);
            }
        }

        return validPointArea;
    }

    bool test_point_area(tracking_matrices& matrices,
                         area_test_settings& settings,
                         const Point2i& targetPoint,
                         test_phase phase,
                         test_behavior outputLog)
    {
        PROFILE_FUNC();
        float area = get_point_area(matrices, settings, targetPoint);

        return test_point_area_core(area, settings, phase, outputLog);
    }


    bool test_point_area_integral(tracking_matrices& matrices,
                                  BitmapF& integralArea,
                                  area_test_settings& settings,
                                  const Point2i& targetPoint,
                                  test_phase phase,
                                  test_behavior outputLog)
    {
        PROFILE_FUNC();
        float area = get_point_area_integral(matrices, integralArea, settings, targetPoint);

        return test_point_area_core(area, settings, phase, outputLog);
    }

    float get_percent_natural_edges(BitmapF& matDepth,
                                    BitmapMask& matSegmentation,
                                    const Point2i& center,
                                    const float bandwidth,
                                    const scaling_coordinate_mapper& mapper)
    {
        PROFILE_FUNC();
        int width = matDepth.width();
        int height = matDepth.height();
        if (center.x < 0 || center.y < 0 ||
            center.x >= width || center.y >= height)
        {
            return 0;
        }

        float startingDepth = matDepth.at(center);

        Point2i  topLeft = mapper.offset_pixel_location_by_mm(center, -bandwidth, bandwidth, startingDepth);

        int offsetX = center.x - topLeft.x;
        int offsetY = center.y - topLeft.y;
        Point2i bottomRight(center.x + offsetX, center.y + offsetY);

        int32_t x0 = MAX(0, topLeft.x);
        int32_t y0 = MAX(0, topLeft.y);
        int32_t x1 = MIN(width - 1, bottomRight.x);
        int32_t y1 = MIN(height - 1, bottomRight.y);

        int naturalEdgeCount = 0;
        int totalEdgeCount = 0;

        for (int y = y0; y <= y1; y++)
        {
            auto* segmentationRow = matSegmentation.data(y);

            segmentationRow += x0;
            for (int x = x0; x <= x1; ++x, ++segmentationRow)
            {
                pixel_type segmentation = static_cast<pixel_type>(*segmentationRow);
                if (segmentation == pixel_type::foreground_natural_edge)
                {
                    ++naturalEdgeCount;
                    ++totalEdgeCount;
                }
                else if (segmentation == pixel_type::foreground_out_of_range_edge)
                {
                    ++totalEdgeCount;
                }
            }
        }

        float percentNaturalEdges = 0;
        if (totalEdgeCount > 0)
        {
            percentNaturalEdges = naturalEdgeCount / static_cast<float>(totalEdgeCount);
        }

        return percentNaturalEdges;
    }

    bool test_natural_edges(tracking_matrices& matrices,
                            natural_edge_test_settings& settings,
                            const Point2i& targetPoint,
                            test_phase phase,
                            test_behavior outputLog)
    {
        PROFILE_FUNC();

        auto scalingMapper = get_scaling_mapper(matrices);
        float percentNaturalEdges = get_percent_natural_edges(matrices.depth,
                                                              matrices.layerSegmentation,
                                                              targetPoint,
                                                              settings.naturalEdgeBandwidth,
                                                              scalingMapper);

        float minPercentNaturalEdges = settings.minPercentNaturalEdges;

        bool passed = percentNaturalEdges >= minPercentNaturalEdges;

        if (outputLog == TEST_BEHAVIOR_LOG)
        {
            if (passed)
            {
                LOG_INFO("Segmentation", "test_natural_edges passed: %f (minimum %f)",
                         percentNaturalEdges,
                         minPercentNaturalEdges);
            }
            else
            {
                LOG_INFO("Segmentation", "test_natural_edges failed: %f (minimum %f)",
                         percentNaturalEdges,
                         minPercentNaturalEdges);
            }
        }
        return passed;
    }

    bool test_foreground_radius_percentage(tracking_matrices& matrices,
                                           circumference_test_settings& settings,
                                           const Point2i& targetPoint,
                                           test_phase phase,
                                           test_behavior outputLog)
    {
        PROFILE_FUNC();
        auto scalingMapper = get_scaling_mapper(matrices);

        std::vector<astra::Vector2i>& points = matrices.layerCirclePoints;

        float percentForeground1 = get_max_sequential_circumference_percentage(matrices.depth,
                                                                               matrices.layerSegmentation,
                                                                               targetPoint,
                                                                               settings.foregroundRadius1,
                                                                               scalingMapper,
                                                                               points);

        float percentForeground2 = get_max_sequential_circumference_percentage(matrices.depth,
                                                                               matrices.layerSegmentation,
                                                                               targetPoint,
                                                                               settings.foregroundRadius2,
                                                                               scalingMapper,
                                                                               points);

        float minPercent1 = settings.foregroundRadiusMinPercent1;
        float minPercent2 = settings.foregroundRadiusMinPercent2;
        const float maxPercent1 = settings.foregroundRadiusMaxPercent1;
        const float maxPercent2 = settings.foregroundRadiusMaxPercent2;

        if (phase == TEST_PHASE_UPDATE)
        {
            //no minimum during update phase
            minPercent1 = 0;
            minPercent2 = 0;
        }

        bool passTest1 = percentForeground1 >= minPercent1 &&
            percentForeground1 <= maxPercent1;

        bool passTest2 = percentForeground2 >= minPercent2 &&
            percentForeground2 <= maxPercent2;

        bool passed = passTest1 && passTest2;

        if (outputLog == TEST_BEHAVIOR_LOG)
        {
            if (passed)
            {
                LOG_INFO("Segmentation", "test_foreground_radius_percentage passed: perc1 %f [%f,%f] perc2 %f [%f,%f]",
                         percentForeground1,
                         minPercent1,
                         maxPercent1,
                         percentForeground2,
                         minPercent2,
                         maxPercent2);
            }
            else
            {
                LOG_INFO("Segmentation", "test_foreground_radius_percentage failed: perc1 %f [%f,%f] perc2 %f [%f,%f]",
                         percentForeground1,
                         minPercent1,
                         maxPercent1,
                         percentForeground2,
                         minPercent2,
                         maxPercent2);
            }
        }
        return passed;
    }

    BitmapF& calculate_integral_area(tracking_matrices& matrices)
    {
        PROFILE_FUNC();
        BitmapMask& segmentationMatrix = matrices.layerSegmentation;
        BitmapF& areaMatrix = matrices.area;
        BitmapF& integralAreaMatrix = matrices.layerIntegralArea;

        integralAreaMatrix.recreate(matrices.depth.size());
        integralAreaMatrix.fill(0.f);

        int width = matrices.depth.width();
        int height = matrices.depth.height();

        float* lastIntegralAreaRow = nullptr;
        for (int y = 0; y < height; y++)
        {
            auto* segmentationRow = segmentationMatrix.data(y);
            float* areaRow = areaMatrix.data(y);
            float* integralAreaRow = integralAreaMatrix.data(y);
            float* integralAreaRowStart = integralAreaRow;

            float leftArea = 0;
            float upLeftArea = 0;

            for (int x = 0; x < width; ++x,
                     ++areaRow,
                     ++integralAreaRow,
                     ++segmentationRow)
            {
                float upArea = 0;
                if (lastIntegralAreaRow != nullptr)
                {
                    upArea = *lastIntegralAreaRow;
                    ++lastIntegralAreaRow;
                }

                float currentArea = 0;

                pixel_type segmentation = static_cast<pixel_type>(*segmentationRow);
                if (segmentation == pixel_type::foreground)
                {
                    currentArea = *areaRow;
                }

                float area = currentArea + leftArea + upArea - upLeftArea;

                *integralAreaRow = area;

                leftArea = area;
                upLeftArea = upArea;
            }

            lastIntegralAreaRow = integralAreaRowStart;
        }

        return integralAreaMatrix;
    }

    bool test_single_point(tracking_data& data, Point2i seedPosition)
    {
        auto matrices = data.matrices;

        auto areaTestSettings = data.settings.areaTestSettings;
        auto circumferenceTestSettings = data.settings.circumferenceTestSettings;
        auto naturalEdgeTestSettings = data.settings.naturalEdgeTestSettings;
        auto integralArea = matrices.layerIntegralArea;

        test_phase phase = data.phase;
        test_behavior outputTestLog = TEST_BEHAVIOR_NONE;

        bool validPointArea = test_point_area_integral(matrices,
                                                       integralArea,
                                                       areaTestSettings,
                                                       seedPosition,
                                                       phase,
                                                       outputTestLog);
        bool validRadiusTest = false;

        if (validPointArea)
        {
            validRadiusTest = test_foreground_radius_percentage(matrices,
                                                                circumferenceTestSettings,
                                                                seedPosition,
                                                                phase,
                                                                outputTestLog);
        }

        bool validNaturalEdgeTest = false;
        if (validRadiusTest)
        {
            validNaturalEdgeTest = test_natural_edges(matrices,
                                                      naturalEdgeTestSettings,
                                                      seedPosition,
                                                      phase,
                                                      outputTestLog);
        }

        bool passesAllTests = (validPointArea && validRadiusTest && validNaturalEdgeTest);

        return passesAllTests;
    }

    foreground_status create_test_pass_from_foreground(tracking_data& data)
    {
        PROFILE_FUNC();
        auto matrices = data.matrices;
        BitmapMask& segmentationMatrix = matrices.layerSegmentation;

        BitmapMask& testPassMatrix = matrices.layerTestPassMap;
        testPassMatrix.recreate(data.matrices.depth.size());
        testPassMatrix.fill(0);

        int width = matrices.depth.width();
        int height = matrices.depth.height();

        auto areaTestSettings = data.settings.areaTestSettings;
        auto circumferenceTestSettings = data.settings.circumferenceTestSettings;
        auto naturalEdgeTestSettings = data.settings.naturalEdgeTestSettings;

        auto integralArea = matrices.layerIntegralArea;

        test_phase phase = data.phase;
        test_behavior outputTestLog = TEST_BEHAVIOR_NONE;

        foreground_status status = foreground_status::empty;

        int xskip = 1;
        int yskip = 1;
        bool downscale = false;
        if (data.phase == TEST_PHASE_CREATE ||
            data.referenceWorldPosition.z < data.settings.maxDepthToDownscaleTestPass)
        {
            //during create cycle or if our point is close enough,
            //downscale test pass map by 2 to save computation
            downscale = true;
            xskip = 2;
            yskip = 2;
        }
        const int maxY = height - yskip;

        for (int y = 0; y <= maxY; y += yskip)
        {
            auto* segmentationRow = segmentationMatrix.data(y);
            auto* testPassRow = testPassMatrix.data(y);
            auto* testPassRowNext = testPassMatrix.data(y+1);

            for (int x = 0; x < width; x += xskip,
                     segmentationRow += xskip,
                     testPassRow += xskip,
                     testPassRowNext += xskip)
            {
                pixel_type segmentation = static_cast<pixel_type>(*segmentationRow);
                if (segmentation != pixel_type::foreground)
                {
                    continue;
                }

                Point2i seedPosition(x, y);
                bool validPointArea = test_point_area_integral(matrices,
                                                               integralArea,
                                                               areaTestSettings,
                                                               seedPosition,
                                                               phase,
                                                               outputTestLog);
                bool validRadiusTest = false;

                if (validPointArea)
                {
                    validRadiusTest = test_foreground_radius_percentage(matrices,
                                                                        circumferenceTestSettings,
                                                                        seedPosition,
                                                                        phase,
                                                                        outputTestLog);
                }

                bool validNaturalEdgeTest = false;
                if (validRadiusTest)
                {
                    validNaturalEdgeTest = test_natural_edges(matrices,
                                                              naturalEdgeTestSettings,
                                                              seedPosition,
                                                              phase,
                                                              outputTestLog);
                }

                if (validPointArea && validRadiusTest && validNaturalEdgeTest)
                {
                    status = foreground_status::has_points;
                    *testPassRow = pixel_type::foreground;
                    if (downscale)
                    {
                        *(testPassRow+1) = pixel_type::foreground;
                        *testPassRowNext = pixel_type::foreground;
                        *(testPassRowNext+1) = pixel_type::foreground;
                    }
                }
            }
        }

        return status;
    }

    Point2i track_point_from_seed(tracking_data& data)
    {
        PROFILE_FUNC();

        Size2i size = data.matrices.depth.size();
        data.matrices.layerSegmentation.recreate(size);
        data.matrices.layerSegmentation.fill(0);

        data.matrices.layerEdgeDistance.recreate(size);
        data.matrices.layerEdgeDistance.fill(0.f);

        data.matrices.layerScore.recreate(size);
        data.matrices.layerScore.fill(0.f);

        const bool debugLayersEnabled = data.matrices.debugLayersEnabled;

        const float layerAverageDepth = segment_foreground_and_get_average_depth(data);

        if (layerAverageDepth == 0.0f || all_zero(data.matrices.layerSegmentation))
        {
            return INVALID_POINT;
        }

        BitmapF& matScore = data.matrices.layerScore;

        calculate_edge_distance(data.matrices.layerSegmentation,
                                data.matrices.areaSqrt,
                                data.matrices.layerEdgeDistance,
                                data.settings.targetEdgeDistance);

        calculate_integral_area(data.matrices);

        calculate_layer_score(data, layerAverageDepth);

        auto minMaxLoc = find_min_max_loc(matScore, data.matrices.layerSegmentation);

        bool foundPoint = minMaxLoc.maxLoc.x != -1 && minMaxLoc.maxLoc.y != -1;

        if (foundPoint)
        {
            bool passesTests = test_single_point(data, minMaxLoc.maxLoc);

            if (!passesTests)
            {
                //our initial point failed the tests, so now we will test all the points
                //and find the max score with the testPassMap

                foreground_status status = create_test_pass_from_foreground(data);
                if (status == foreground_status::empty)
                {
                    foundPoint = false;
                }
                else
                {
                    minMaxLoc = find_min_max_loc(matScore, data.matrices.layerTestPassMap);
                }
            }
        }

        if (debugLayersEnabled)
        {
            ++data.matrices.layerCount;

            BitmapMask layerCountMat(size);
            layerCountMat.fill(data.matrices.layerCount);

            bitwise_or(layerCountMat,
                       data.matrices.debugSegmentation,
                       data.matrices.debugSegmentation,
                       data.matrices.layerSegmentation);

            bitwise_or(layerCountMat,
                       data.matrices.debugTestPassMap,
                       data.matrices.debugTestPassMap,
                       data.matrices.layerTestPassMap);

            BitmapMask scoreMask;
            in_range(matScore, scoreMask, 1, std::numeric_limits<int>::max());
            copy_to(matScore, data.matrices.debugScoreValue, scoreMask);

            range_normalize(matScore, data.matrices.debugScore, 0, 1, scoreMask);
        }

        if (!foundPoint)
        {
            return INVALID_POINT;
        }

        return minMaxLoc.maxLoc;
    }

    bool find_next_velocity_seed_pixel(BitmapMask& velocitySignalMatrix,
                                       BitmapMask& searchedMatrix,
                                       Point2i& foregroundPosition,
                                       Point2i& nextSearchStart)
    {
        PROFILE_FUNC();
        assert(velocitySignalMatrix.size() == searchedMatrix.size());

        int width = velocitySignalMatrix.width();
        int height = velocitySignalMatrix.height();

        int startX = MAX(0, MIN(width - 1, nextSearchStart.x));
        const int startY = MAX(0, MIN(height - 1, nextSearchStart.y));

        for (int y = startY; y < height; y++)
        {
            for (int x = startX; x < width; x++)
            {
                pixel_type velocitySignal = static_cast<pixel_type>(velocitySignalMatrix.at(x, y));
                pixel_type searched = static_cast<pixel_type>(searchedMatrix.at(x, y));
                if (velocitySignal == pixel_type::foreground && searched != pixel_type::searched)
                {
                    foregroundPosition.x = x;
                    foregroundPosition.y = y;

                    nextSearchStart.x = x + 1;
                    if (nextSearchStart.x < width)
                    {
                        nextSearchStart.y = y;
                    }
                    else
                    {
                        nextSearchStart.x = 0;
                        nextSearchStart.y = y + 1;
                        if (nextSearchStart.y >= height)
                        {
                            return false;
                        }
                    }
                    return true;
                }
            }
            startX = 0;
        }

        foregroundPosition = segmentation::INVALID_POINT;
        nextSearchStart.x = width;
        nextSearchStart.y = height;
        return false;
    }

    void calculate_edge_distance(BitmapMask& segmentationMatrix,
                                 BitmapF& areaSqrtMatrix,
                                 BitmapF& edgeDistanceMatrix,
                                 const float maxEdgeDistance)
    {
        PROFILE_FUNC();
        BitmapF eroded;
        BitmapMask crossElement = get_structuring_element(MorphShape::Cross, Size2i(3, 3));

        edgeDistanceMatrix.recreate(segmentationMatrix.size());
        edgeDistanceMatrix.fill(0.f);

        BitmapF ones(segmentationMatrix.size());
        ones.fill(1.f);

        convert_to(segmentationMatrix, eroded);

        //close small holes
        int dilateCount = 1;
        for (int i = 0; i < dilateCount; i++)
        {
            dilate(eroded, eroded, crossElement);
        }

        int nonZeroCount = 0;
        const int imageLength = eroded.width() * eroded.height();
        int iterations = 0;
        const int maxIterations = segmentationMatrix.width() / 2;
        bool done;

        do
        {
            PROFILE_BEGIN(edge_dist_loop);
            //erode makes the image smaller
            erode(eroded, eroded, crossElement);
            //accumulate the eroded image to the edgeDistance buffer
            scalar_add(areaSqrtMatrix, edgeDistanceMatrix, edgeDistanceMatrix, eroded);

            nonZeroCount = count_non_zero(eroded);
            done = (nonZeroCount == 0);

            auto minMax = find_min_max(edgeDistanceMatrix);

            if (minMax.max > maxEdgeDistance)
            {
                done = true;
            }

            PROFILE_END();
            //nonZeroCount < imageLength guards against segmentation with all 1's, which will never erode
        } while (!done && nonZeroCount < imageLength && ++iterations < maxIterations);
    }

    void get_circumference_points(BitmapF& matDepth,
                                  const Point2i& center,
                                  const float& radius,
                                  const scaling_coordinate_mapper& mapper,
                                  std::vector<astra::Vector2i>& points)
    {
        PROFILE_FUNC();

        int width = matDepth.width();
        int height = matDepth.height();
        if (center.x < 0 || center.x >= width ||
            center.y < 0 || center.y >= height ||
            radius < 1)
        {
            return;
        }

        float referenceDepth = matDepth.at(center);
        if (referenceDepth == 0)
        {
            return;
        }

        Point2i offsetRight = mapper.offset_pixel_location_by_mm(center, radius, 0, referenceDepth);

        //http://en.wikipedia.org/wiki/Midpoint_circle_algorithm
        int pixelRadius = offsetRight.x - center.x;
        int cx = center.x;
        int cy = center.y;

        std::vector<astra::Vector2i> offsets;
        //reserve a slight overestimation of number of points for 1/8 of circumference
        offsets.reserve(pixelRadius);

        {
            int dx = pixelRadius; //radius in pixels
            int dy = 0;
            int radiusError = 1 - dx;

            while (dx >= dy)
            {
                offsets.push_back(Vector2i(dx, dy));

                dy++;
                if (radiusError < 0)
                {
                    radiusError += 2 * dy + 1;
                }
                else
                {
                    dx--;
                    radiusError += 2 * (dy - dx) + 1;
                }
            }
        }

        //PROFILE_BEGIN(circ_checks);

        //clear & reuse capacity across calls
        points.clear();
        points.reserve(static_cast<int>(pixelRadius * 2.0f * PI_F));

        int length = offsets.size();

        //Order and the permutations of dx,dy are critical here
        //so the points list will contain the points in order

        for (int i = 1; i < length; ++i)
        {
            //dx, dy
            const astra::Vector2i delta = offsets[i];
            const int x = cx + delta.x;
            const int y = cy + delta.y;

            if (x >= 0 && x < width &&
                y >= 0 && y < height)
            {
                points.push_back(astra::Vector2i(x, y));
            }
        }

        //even quadrants are reversed order
        for (int i = length-1; i >= 0; --i)
        {
            //dy, dx
            const astra::Vector2i delta = offsets[i];

            const int dx = delta.x;
            const int dy = delta.y;
            if (dx != dy)
            {
                const int x = cx + dy;
                const int y = cy + dx;
                if (x >= 0 && x < width &&
                    y >= 0 && y < height)
                {
                    points.push_back(astra::Vector2i(x, y));
                }
            }
        }

        for (int i = 1; i < length; ++i)
        {
            //-dy, dx
            const astra::Vector2i delta = offsets[i];
            const int x = cx - delta.y;
            const int y = cy + delta.x;

            if (x >= 0 && x < width &&
                y >= 0 && y < height)
            {
                points.push_back(astra::Vector2i(x, y));
            }
        }

        for (int i = length-1; i >= 0; --i)
        {
            //-dx, dy
            const astra::Vector2i delta = offsets[i];

            const int dx = delta.x;
            const int dy = delta.y;
            if (dx != dy)
            {
                const int x = cx - dx;
                const int y = cy + dy;
                if (x >= 0 && x < width &&
                    y >= 0 && y < height)
                {
                    points.push_back(astra::Vector2i(x, y));
                }
            }
        }

        for (int i = 1; i < length; ++i)
        {
            //-dx, -dy
            const astra::Vector2i delta = offsets[i];
            const int x = cx - delta.x;
            const int y = cy - delta.y;

            if (x >= 0 && x < width &&
                y >= 0 && y < height)
            {
                points.push_back(astra::Vector2i(x, y));
            }
        }

        for (int i = length-1; i >= 0; --i)
        {
            //-dy, -dx
            const astra::Vector2i delta = offsets[i];

            const int dx = delta.x;
            const int dy = delta.y;
            if (dx != dy)
            {
                const int x = cx - dy;
                const int y = cy - dx;
                if (x >= 0 && x < width &&
                    y >= 0 && y < height)
                {
                    points.push_back(astra::Vector2i(x, y));
                }
            }
        }

        for (int i = 1; i < length; ++i)
        {
            //dy, -dx
            const astra::Vector2i delta = offsets[i];
            const int x = cx + delta.y;
            const int y = cy - delta.x;

            if (x >= 0 && x < width &&
                y >= 0 && y < height)
            {
                points.push_back(astra::Vector2i(x, y));
            }
        }

        for (int i = length-1; i >= 0; --i)
        {
            //dx, -dy
            const astra::Vector2i delta = offsets[i];

            const int dx = delta.x;
            const int dy = delta.y;
            if (dx != dy)
            {
                const int x = cx + dx;
                const int y = cy - dy;
                if (x >= 0 && x < width &&
                    y >= 0 && y < height)
                {
                    points.push_back(astra::Vector2i(x, y));
                }
            }
        }
        //PROFILE_END();
    }

    float get_max_sequential_circumference_percentage(BitmapF& matDepth,
                                                      BitmapMask& matSegmentation,
                                                      const Point2i& center,
                                                      const float& radius,
                                                      const scaling_coordinate_mapper& mapper,
                                                      std::vector<astra::Vector2i>& points)
    {
        PROFILE_FUNC();
        int foregroundCount = 0;
        int maxCount = 0;
        int firstSegmentCount = 0;
        bool isFirstSegment = true;
        int totalCount = 0;
        bool firstIsForeground = false;
        bool lastIsForeground = false;

        get_circumference_points(matDepth, center, radius, mapper, points);

        for (auto p : points)
        {
            bool isForeground =
                matSegmentation.at(p.x, p.y) == pixel_type::foreground;
            if (isForeground)
            {
                ++foregroundCount;
                if (totalCount == 0)
                {
                    //started on a foreground segment, might need to add to the last segment
                    firstIsForeground = true;
                }
            }
            else
            {
                if (foregroundCount > maxCount)
                {
                    maxCount = foregroundCount;
                }
                if (isFirstSegment)
                {
                    firstSegmentCount = foregroundCount;
                    isFirstSegment = false;
                }
                foregroundCount = 0;
            }
            lastIsForeground = isForeground;
            ++totalCount;
        }

        if (lastIsForeground)
        {
            if (firstIsForeground)
            {
                //add the first and last segment
                foregroundCount += firstSegmentCount;
            }

            if (foregroundCount > maxCount)
            {
                maxCount = foregroundCount;
            }
        }

        float percentForeground = maxCount / static_cast<float>(totalCount);

        return percentForeground;
    }

    float count_neighborhood_area(BitmapMask& matSegmentation,
                                  BitmapF& matDepth,
                                  BitmapF& matArea,
                                  const Point2i& center,
                                  const float bandwidth,
                                  const float bandwidthDepth,
                                  const scaling_coordinate_mapper& mapper)
    {
        PROFILE_FUNC();
        if (center.x < 0 || center.y < 0 ||
            center.x >= matDepth.width() || center.y >= matDepth.height())
        {
            return 0;
        }

        float startingDepth = matDepth.at(center);

        Point2i topLeft = mapper.offset_pixel_location_by_mm(center, -bandwidth, bandwidth, startingDepth);

        int offsetX = center.x - topLeft.x;
        int offsetY = center.y - topLeft.y;
        Point2i bottomRight(center.x + offsetX, center.y + offsetY);

        int32_t x0 = MAX(0, topLeft.x);
        int32_t y0 = MAX(0, topLeft.y);
        int32_t x1 = MIN(matDepth.width() - 1, bottomRight.x);
        int32_t y1 = MIN(matDepth.height() - 1, bottomRight.y);

        float area = 0;

        for (int y = y0; y <= y1; y++)
        {
            auto* segmentationRow = matSegmentation.data(y);
            float* areaRow = matArea.data(y);

            segmentationRow += x0;
            areaRow += x0;
            for (int x = x0; x <= x1; ++x, ++areaRow, ++segmentationRow)
            {
                pixel_type segmentation = static_cast<pixel_type>(*segmentationRow);
                if (segmentation == pixel_type::foreground)
                {
                    area += *areaRow;
                }
            }
        }

        return area;
    }


    float count_neighborhood_area_integral(BitmapF& matDepth,
                                           BitmapF& matAreaIntegral,
                                           const Point2i& center,
                                           const float bandwidth,
                                           const scaling_coordinate_mapper& mapper)
    {
        //PROFILE_FUNC();
        int width = matDepth.width();
        int height = matDepth.height();
        if (center.x < 0 || center.y < 0 ||
            center.x >= width || center.y >= height)
        {
            return 0;
        }

        float startingDepth = matDepth.at(center);

        Point2i topLeft = mapper.offset_pixel_location_by_mm(center, -bandwidth, bandwidth, startingDepth);

        int offsetX = center.x - topLeft.x;
        int offsetY = center.y - topLeft.y;
        Point2i bottomRight(center.x + offsetX, center.y + offsetY);

        //subtract one from topLeft because formula below is exclusive on the lower bounds
        int32_t x0 = MAX(0, topLeft.x-1);
        int32_t y0 = MAX(0, topLeft.y-1);
        int32_t x1 = MIN(width - 1, bottomRight.x);
        int32_t y1 = MIN(height - 1, bottomRight.y);

        float area = 0;

        area += matAreaIntegral.at(x1, y1);
        area += matAreaIntegral.at(x0, y0);
        area -= matAreaIntegral.at(x1, y0);
        area -= matAreaIntegral.at(x0, y1);

        return area;
    }

}}}
