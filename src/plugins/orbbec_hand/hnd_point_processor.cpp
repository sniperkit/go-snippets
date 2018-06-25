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
#include "hnd_tracked_point.hpp"
#include "hnd_point_processor.hpp"
#include "hnd_segmentation.hpp"
#include <Shiny.h>

namespace astra { namespace hand {

    point_processor::point_processor(point_processor_settings& settings) :
        settings_(settings)
    {
        PROFILE_FUNC();
    }

    point_processor::~point_processor()
    {
        PROFILE_FUNC();
    }

    void point_processor::calculate_area(tracking_matrices& matrices, scaling_coordinate_mapper mapper)
    {
        PROFILE_FUNC();

        const astra::Vector3f* fullSizeWorldPoints = matrices.fullSizeWorldPoints;
        astra::Vector3f* worldPoints = matrices.worldPoints;
        auto depthToWorldData = matrices.depthToWorldData;
        BitmapF& areaMatrix = matrices.area;
        BitmapF& areaSqrtMatrix = matrices.areaSqrt;

        Size2i depthSize = matrices.depth.size();

        areaMatrix.recreate(depthSize);
        areaMatrix.fill(0.f);

        areaSqrtMatrix.recreate(depthSize);
        areaSqrtMatrix.fill(0.f);

        int fullSizeWidth = matrices.depthFullSize.width();
        int width = depthSize.width();
        int height = depthSize.height();

        float offsetX = mapper.offsetX();
        float offsetY = mapper.offsetY();
        float scale = mapper.scale();
        int intScale = static_cast<int>(scale);

        for (int y = 0; y < height; ++y)
        {
            float* areaRow = areaMatrix.data(y);
            float* areaSqrtRow = areaSqrtMatrix.data(y);

            for (int x = 0; x < width; ++x, ++worldPoints, ++areaRow, ++areaSqrtRow)
            {
                int fullSizeIndex = (x + y * fullSizeWidth) * intScale;
                const astra::Vector3f& p = fullSizeWorldPoints[fullSizeIndex];
                *worldPoints = p;
                const float depth = p.z;

                if (depth != 0)
                {
                    float depthX = (x + 1 + offsetX) * scale;
                    float depthY = (y + 1 + offsetY) * scale;
                    float normalizedX = depthX / depthToWorldData.resolutionX - .5f;
                    float normalizedY = .5f - depthY / depthToWorldData.resolutionY;

                    float wx = normalizedX * depth * depthToWorldData.xzFactor;
                    float wy = normalizedY * depth * depthToWorldData.yzFactor;

                    float deltaX = wx - p.x;
                    float deltaY = wy - p.y;

                    float area = fabs(deltaX * deltaY);

                    *areaRow = area;
                    *areaSqrtRow = sqrt(area);
                }
                else
                {
                    *areaRow = 0;
                    *areaSqrtRow = 0;
                }
            }
        }
    }

    void point_processor::initialize_common_calculations(tracking_matrices& matrices)
    {
        PROFILE_FUNC();
        auto scalingMapper = get_scaling_mapper(matrices);

        calculate_area(matrices, scalingMapper);
    }

    void point_processor::update_tracked_points(tracking_matrices& matrices)
    {
        PROFILE_FUNC();
        auto scalingMapper = get_scaling_mapper(matrices);

        //give priority updates to active points
        for (auto iter = trackedPoints_.begin(); iter != trackedPoints_.end(); ++iter)
        {
            tracked_point& trackedPoint = *iter;
            if (trackedPoint.pointType == tracked_point_type::active_point)
            {
                update_tracked_point(matrices, scalingMapper, trackedPoint);
            }
        }

        int numUpdatedPoints = 0;
        for (auto iter = trackedPoints_.begin(); iter != trackedPoints_.end(); ++iter)
        {
            tracked_point& trackedPoint = *iter;
            if (trackedPoint.pointType != tracked_point_type::active_point)
            {
                update_tracked_point(matrices, scalingMapper, trackedPoint);
            }
            ++numUpdatedPoints;
            if (numUpdatedPoints > settings_.maxhandpointUpdatesPerFrame)
            {
                break;
            }
        }
    }

    void point_processor::update_tracked_point(tracking_matrices& matrices,
                                               scaling_coordinate_mapper& scalingMapper,
                                               tracked_point& trackedPoint)
    {
        PROFILE_FUNC();
        const float width = matrices.depth.width();
        const float height = matrices.depth.height();

        ++trackedPoint.inactiveFrameCount;

        tracking_data updatetracking_data(matrices,
                                          trackedPoint.position,
                                          trackedPoint.worldPosition,
                                          trackedPoint.referenceAreaSqrt,
                                          VELOCITY_POLICY_IGNORE,
                                          settings_.segmentationSettings,
                                          TEST_PHASE_UPDATE);

        Point2i newTargetPoint = segmentation::track_point_from_seed(updatetracking_data);

        validate_and_update_tracked_point(matrices, scalingMapper, trackedPoint, newTargetPoint);

        //lost a tracked point, try to guess the position using previous position delta for second chance to recover

        auto xyDelta = trackedPoint.worldDeltaPosition;
        xyDelta.z = 0;
        double xyDeltaNorm = xyDelta.length();
        if (trackedPoint.trackingStatus != tracking_status::tracking &&
            newTargetPoint == segmentation::INVALID_POINT &&
            xyDeltaNorm > settings_.secondChanceMinDistance)
        {
            auto movementDirection = xyDelta * (1.0f / xyDeltaNorm);
            float maxSegmentationDist = settings_.segmentationSettings.maxSegmentationDist;
            auto estimatedWorldPosition = trackedPoint.worldPosition + movementDirection * maxSegmentationDist;

            Vector3f estimatedPosition = scalingMapper.convert_world_to_depth(estimatedWorldPosition);

            Point2i seedPosition;
            seedPosition.x = MAX(0, MIN(width - 1, static_cast<int>(estimatedPosition.x)));
            seedPosition.y = MAX(0, MIN(height - 1, static_cast<int>(estimatedPosition.y)));

            tracking_data recovertracking_data(matrices,
                                               seedPosition,
                                               estimatedWorldPosition,
                                               trackedPoint.referenceAreaSqrt,
                                               VELOCITY_POLICY_IGNORE,
                                               settings_.segmentationSettings,
                                               TEST_PHASE_UPDATE);

            newTargetPoint = segmentation::track_point_from_seed(recovertracking_data);

            //test for invalid point here so we don't increment failed test counts
            //for second chance recovery
            if (newTargetPoint != segmentation::INVALID_POINT)
            {
                validate_and_update_tracked_point(matrices, scalingMapper, trackedPoint, newTargetPoint);
            }

            if (trackedPoint.trackingStatus == tracking_status::tracking)
            {
                LOG_TRACE("point_processor", "update_tracked_point 2nd chance recovered #%d",
                          trackedPoint.trackingId);
            }
        }
    }

    void point_processor::reset()
    {
        PROFILE_FUNC();
        trackedPoints_.clear();
        nextTrackingId_ = 0;
    }

    void point_processor::update_full_resolution_points(tracking_matrices& matrices)
    {
        PROFILE_FUNC();
        for (auto iter = trackedPoints_.begin(); iter != trackedPoints_.end(); ++iter)
        {
            tracked_point& trackedPoint = *iter;

            float resizeFactor = get_resize_factor(matrices);

            //add 0.5 to center on the middle of the pixel
            trackedPoint.fullSizePosition.x = (trackedPoint.position.x + 0.5) * resizeFactor;
            trackedPoint.fullSizePosition.y = (trackedPoint.position.y + 0.5) * resizeFactor;

            bool resizeNeeded = matrices.depthFullSize.width() != matrices.depth.width();

            bool processRefinedPosition = false;

            if (resizeNeeded &&
                trackedPoint.trackingStatus == tracking_status::tracking &&
                trackedPoint.pointType == tracked_point_type::active_point)
            {
                Vector3f refinedWorldPosition = trackedPoint.worldPosition;
                if (processRefinedPosition)
                {
                    refinedWorldPosition = get_refined_high_res_position(matrices, trackedPoint);
                }

                Vector3f smoothedWorldPosition = smooth_world_positions(trackedPoint.fullSizeWorldPosition, refinedWorldPosition);

                update_tracked_point_from_world_position(trackedPoint,
                                                         smoothedWorldPosition,
                                                         resizeFactor,
                                                         matrices.depthToWorldData);
            }
            else
            {
                trackedPoint.fullSizeWorldPosition = trackedPoint.worldPosition;
                trackedPoint.fullSizeWorldDeltaPosition = trackedPoint.worldDeltaPosition;
            }
        }
    }

    void point_processor::update_trajectories()
    {
        PROFILE_FUNC();
        for (auto iter = trackedPoints_.begin(); iter != trackedPoints_.end(); ++iter)
        {
            //TODO take this and make it a method on tracked_point
            tracked_point& trackedPoint = *iter;
            int trackingId = trackedPoint.trackingId;
            auto it = trajectories_.find(trackingId);

            if (it == trajectories_.end())
            {
                trajectory_analyzer analyzer(trackingId, settings_.trajectoryAnalyzerSettings);
                analyzer.update(trackedPoint);
                trajectories_.insert(std::make_pair(trackingId, analyzer));
            }
            else
            {
                trajectory_analyzer& analyzer = it->second;
                analyzer.update(trackedPoint);

                if (analyzer.is_wave_gesture())
                {
                    end_probation(trackedPoint);
                    trackedPoint.pointType = tracked_point_type::active_point;
                }
            }
        }
    }

    Vector3f point_processor::get_refined_high_res_position(tracking_matrices& matrices,
                                                            const tracked_point& trackedPoint)
    {
        PROFILE_FUNC();
        assert(trackedPoint.pointType == tracked_point_type::active_point);

        if (trackedPoint.worldPosition.z == 0)
        {
            return trackedPoint.worldPosition;
        }

        int fullWidth = matrices.depthFullSize.width();
        int fullHeight = matrices.depthFullSize.height();
        int processingWidth = matrices.depth.width();
        int processingHeight = matrices.depth.height();

        int fullSizeX = trackedPoint.fullSizePosition.x;
        int fullSizeY = trackedPoint.fullSizePosition.y;
        int windowLeft = MAX(0, MIN(fullWidth - processingWidth, fullSizeX - processingWidth / 2));
        int windowTop = MAX(0, MIN(fullHeight - processingHeight, fullSizeY - processingHeight / 2));

        Point2i roiPosition(fullSizeX - windowLeft, fullSizeY - windowTop);

        float referenceAreaSqrt = matrices.areaSqrt.at(roiPosition);
        if (referenceAreaSqrt == 0)
        {
            return trackedPoint.worldPosition;
        }

        //Create a window into the full size data
        auto roi = matrices.depthFullSize.submap(windowLeft, windowTop, processingWidth, processingHeight);
        //copyTo for now so .at works with local coords in functions that use it
        copy_to(roi, matrices.depth);

        //initialize_common_calculations(matrices);
        scaling_coordinate_mapper roiMapper(matrices.depthToWorldData, 1.0, windowLeft, windowTop);

        calculate_area(matrices, roiMapper);

        tracking_data refinementtracking_data(matrices,
                                              roiPosition,
                                              trackedPoint.worldPosition,
                                              referenceAreaSqrt,
                                              VELOCITY_POLICY_IGNORE,
                                              settings_.segmentationSettings,
                                              TEST_PHASE_UPDATE);

        Point2i targetPoint = segmentation::track_point_from_seed(refinementtracking_data);

        if (targetPoint == segmentation::INVALID_POINT)
        {
            return trackedPoint.worldPosition;
        }

        int refinedFullSizeX = targetPoint.x + windowLeft;
        int refinedFullSizeY = targetPoint.y + windowTop;

        float refinedDepth = matrices.depthFullSize.at(refinedFullSizeY, refinedFullSizeX);

        if (refinedDepth == 0)
        {
            refinedDepth = trackedPoint.worldPosition.z;
        }

        Vector3f refinedWorldPosition = cv_convert_depth_to_world(matrices.depthToWorldData,
                                                                  refinedFullSizeX,
                                                                  refinedFullSizeY,
                                                                  refinedDepth);

        return refinedWorldPosition;
    }

    Vector3f point_processor::smooth_world_positions(const Vector3f& oldWorldPosition,
                                                     const Vector3f& newWorldPosition)
    {
        PROFILE_FUNC();
        float smoothingFactor = settings_.pointSmoothingFactor;

        float delta = (newWorldPosition - oldWorldPosition).length();
        if (delta < settings_.pointSmoothingDeadZone)
        {
            float factorRamp = delta / settings_.pointSmoothingDeadZone;
            smoothingFactor = settings_.pointSmoothingFactor * factorRamp +
                settings_.pointDeadBandSmoothingFactor * (1 - factorRamp);
        }

        return oldWorldPosition * (1 - smoothingFactor) + newWorldPosition * smoothingFactor;
    }

    void point_processor::update_tracked_point_from_world_position(tracked_point& trackedPoint,
                                                                   const Vector3f& newWorldPosition,
                                                                   const float resizeFactor,
                                                                   const conversion_cache_t& depthToWorldData)
    {
        PROFILE_FUNC();
        Vector3f fullSizeDepthPosition = cv_convert_world_to_depth(depthToWorldData, newWorldPosition);

        trackedPoint.fullSizePosition = Point2i(fullSizeDepthPosition.x, fullSizeDepthPosition.y);

        Vector3f deltaPosition = newWorldPosition - trackedPoint.fullSizeWorldPosition;
        trackedPoint.fullSizeWorldPosition = newWorldPosition;
        trackedPoint.fullSizeWorldDeltaPosition = deltaPosition;
    }

    void point_processor::start_probation(tracked_point& trackedPoint)
    {
        PROFILE_FUNC();
        if (!trackedPoint.isInProbation)
        {
            LOG_TRACE("point_processor", "started probation for: %d", trackedPoint.trackingId);
            trackedPoint.isInProbation = true;
            trackedPoint.probationFrameCount = 0;
            trackedPoint.failedTestCount = 0;
        }
    }

    void point_processor::end_probation(tracked_point& trackedPoint)
    {
        PROFILE_FUNC();
        trackedPoint.isInProbation = false;
        trackedPoint.failedTestCount = 0;
    }

    void point_processor::update_tracked_point_data(tracking_matrices& matrices, scaling_coordinate_mapper& scalingMapper, tracked_point& trackedPoint, const Point2i& newTargetPoint)
    {
        PROFILE_FUNC();
        float depth = matrices.depth.at(newTargetPoint);
        Vector3f worldPosition = scalingMapper.convert_depth_to_world(newTargetPoint.x, newTargetPoint.y, depth);

        Vector3f deltaPosition = worldPosition - trackedPoint.worldPosition;
        trackedPoint.worldPosition = worldPosition;
        trackedPoint.worldDeltaPosition = deltaPosition;

        trackedPoint.position = newTargetPoint;
        trackedPoint.referenceAreaSqrt = matrices.areaSqrt.at(trackedPoint.position);

        auto steadyDist = (worldPosition - trackedPoint.steadyWorldPosition).length();

        if (steadyDist > settings_.steadyDeadBandRadius)
        {
            trackedPoint.steadyWorldPosition = worldPosition;
            trackedPoint.inactiveFrameCount = 0;
        }
    }

    void point_processor::validate_and_update_tracked_point(tracking_matrices& matrices,
                                                            scaling_coordinate_mapper& scalingMapper,
                                                            tracked_point& trackedPoint,
                                                            const Point2i& newTargetPoint)
    {
        PROFILE_FUNC();
        if(trackedPoint.trackingStatus == tracking_status::dead)
        {
            return;
        }

        auto oldStatus = trackedPoint.trackingStatus;
        const test_behavior outputTestLog = TEST_BEHAVIOR_NONE;

        bool validPointInRange = segmentation::test_point_in_range(matrices,
                                                                   newTargetPoint,
                                                                   outputTestLog);

        if (validPointInRange)
        {
            trackedPoint.trackingStatus = tracking_status::tracking;
            if (trackedPoint.pointType == tracked_point_type::active_point)
            {
                trackedPoint.failedTestCount = 0;
            }
        }
        else
        {
            start_probation(trackedPoint);
            ++trackedPoint.failedTestCount;
        }

        if (trackedPoint.isInProbation)
        {
            bool probationFailed = false;
            if (trackedPoint.pointType == tracked_point_type::active_point &&
                trackedPoint.failedTestCount >= settings_.maxFailedTestsInProbationActivePoints)
            {
                //had valid in range points but must have failed the real tests

                //failed N consecutive tests within the probation period
                //gave the active point a few extra frames to recover
                trackedPoint.trackingStatus = tracking_status::lost;
                probationFailed = true;
                LOG_TRACE("point_processor", "lost an active point: %d", trackedPoint.trackingId);
            }
            else if (trackedPoint.pointType == tracked_point_type::candidate_point &&
                     trackedPoint.failedTestCount >= settings_.maxFailedTestsInProbation)
            {
                //failed N tests total (non-consecutive) within the probation period
                //too many failed tests, so long...
                trackedPoint.trackingStatus = tracking_status::lost;
                probationFailed = true;
            }

            ++trackedPoint.probationFrameCount;
            if (trackedPoint.probationFrameCount > settings_.probationFrameCount || probationFailed)
            {
                //you're out of probation, but we're keeping an eye on you...
                end_probation(trackedPoint);
                LOG_TRACE("point_processor", "ended probation: %d count: %d/%d probationFailed: %d",
                          trackedPoint.trackingId,
                          trackedPoint.probationFrameCount,
                          settings_.probationFrameCount,
                          probationFailed);
            }
        }

        if (validPointInRange)
        {
            update_tracked_point_data(matrices, scalingMapper, trackedPoint, newTargetPoint);
        }

        if (trackedPoint.trackingStatus != oldStatus)
        {
            LOG_TRACE("point_processor", "validate_and_update_tracked_point: #%d status %s --> status %s",
                      trackedPoint.trackingId,
                      tracking_status_to_string(oldStatus).c_str(),
                      tracking_status_to_string(trackedPoint.trackingStatus).c_str());
        }
    }

    void point_processor::remove_duplicate_points()
    {
        PROFILE_FUNC();

        for (auto iter = trackedPoints_.begin(); iter != trackedPoints_.end(); ++iter)
        {
            tracked_point& tracked = *iter;
            for (auto otherIter = trackedPoints_.begin(); otherIter != trackedPoints_.end(); ++otherIter)
            {
                tracked_point& otherTracked = *otherIter;
                bool bothNotDead = tracked.trackingStatus != tracking_status::dead && otherTracked.trackingStatus != tracking_status::dead;
                float pointDist = (tracked.worldPosition - otherTracked.worldPosition).length();
                if (tracked.trackingId != otherTracked.trackingId &&
                    bothNotDead &&
                    pointDist < settings_.mergePointDistance)
                {
                    tracked.inactiveFrameCount = MIN(tracked.inactiveFrameCount, otherTracked.inactiveFrameCount);
                    if (otherTracked.pointType == tracked_point_type::active_point && tracked.pointType != tracked_point_type::active_point)
                    {
                        tracked.trackingId = otherTracked.trackingId;
                        tracked.pointType = tracked_point_type::active_point;
                    }
                    otherTracked.trackingStatus = tracking_status::dead;
                }
            }
        }
    }

    void point_processor::remove_stale_or_dead_points()
    {
        PROFILE_FUNC();
        for (auto iter = trackedPoints_.begin(); iter != trackedPoints_.end();)
        {
            tracked_point& tracked = *iter;

            int max = settings_.maxInactiveFramesForCandidatePoints;
            if (tracked.pointType == tracked_point_type::active_point)
            {
                if (tracked.trackingStatus == tracking_status::lost)
                {
                    max = settings_.maxInactiveFramesForLostPoints;
                }
                else
                {
                    max = settings_.maxInactiveFramesForActivePoints;
                }
            }
            //if inactive for more than a certain number of frames, or dead, remove point
            if (tracked.inactiveFrameCount > max || tracked.trackingStatus == tracking_status::dead)
            {
                trajectories_.erase(tracked.trackingId);

                iter = trackedPoints_.erase(iter);
            }
            else
            {
                ++iter;
            }
        }
    }

    void point_processor::update_tracked_or_create_new_point_from_seed(tracking_matrices& matrices,
                                                                       const Point2i& seedPosition)
    {
        PROFILE_FUNC();
        float referenceDepth = matrices.depth.at(seedPosition);
        float referenceAreaSqrt = matrices.areaSqrt.at(seedPosition);
        if (referenceDepth == 0 || referenceAreaSqrt == 0 || seedPosition == segmentation::INVALID_POINT)
        {
            //Cannot expect to properly segment when the seedPosition has zero depth
            return;
        }

        auto scalingMapper = get_scaling_mapper(matrices);
        Vector3f referenceWorldPosition =
            scalingMapper.convert_depth_to_world(seedPosition.x,
                                                 seedPosition.y,
                                                 referenceDepth);

        tracking_data createtracking_data(matrices,
                                          seedPosition,
                                          referenceWorldPosition,
                                          referenceAreaSqrt,
                                          VELOCITY_POLICY_RESET_TTL,
                                          settings_.segmentationSettings,
                                          TEST_PHASE_CREATE);

        Point2i targetPoint = segmentation::track_point_from_seed(createtracking_data);

        const test_behavior outputTestLog = TEST_BEHAVIOR_NONE;

        bool validPointInRange = segmentation::test_point_in_range(matrices,
                                                                   targetPoint,
                                                                   outputTestLog);

        if (!validPointInRange)
        {
            return;
        }

        bool existingPoint = false;

        float depth = matrices.depth.at(targetPoint);

        Vector3f worldPosition = scalingMapper.convert_depth_to_world(targetPoint.x,
                                                                      targetPoint.y,
                                                                      depth);

        for (auto iter = trackedPoints_.begin(); iter != trackedPoints_.end(); ++iter)
        {
            tracked_point& trackedPoint = *iter;
            if (trackedPoint.trackingStatus != tracking_status::dead)
            {
                float dist = (trackedPoint.worldPosition - worldPosition).length();
                float maxDist = settings_.maxMatchDistDefault;
                bool lostPoint = trackedPoint.trackingStatus == tracking_status::lost;
                if (lostPoint && trackedPoint.pointType == tracked_point_type::active_point)
                {
                    maxDist = settings_.maxMatchDistLostActive;
                }
                if (dist < maxDist)
                {
                    trackedPoint.inactiveFrameCount = 0;
                    if (lostPoint)
                    {
                        //Recover a lost point -- move it to the recovery position
                        trackedPoint.position = targetPoint;
                        trackedPoint.referenceAreaSqrt = matrices.areaSqrt.at(trackedPoint.position);

                        trackedPoint.worldPosition = worldPosition;
                        trackedPoint.worldDeltaPosition = Vector3f();

                        LOG_TRACE("point_processor", "createCycle: Recovered #%d",
                                  trackedPoint.trackingId);

                        //it could be faulty recovery, so start out in probation just like a new point
                        start_probation(trackedPoint);
                    }
                    trackedPoint.trackingStatus = tracking_status::tracking;
                    existingPoint = true;
                    break;
                }
            }
        }
        if (!existingPoint)
        {
            LOG_TRACE("point_processor", "createCycle: Created new point #%d",
                      nextTrackingId_);

            tracked_point newPoint(targetPoint, worldPosition, nextTrackingId_);
            newPoint.pointType = tracked_point_type::candidate_point;
            newPoint.trackingStatus = tracking_status::tracking;
            ++nextTrackingId_;
            trackedPoints_.push_back(newPoint);
            start_probation(newPoint);
        }
    }
}}
