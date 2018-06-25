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
#ifndef HND_SEGMENTATION_H
#define HND_SEGMENTATION_H

#include "hnd_scaling_coordinate_mapper.hpp"
#include "hnd_point.hpp"

#ifndef MIN
#define MIN(a,b) (((a)<(b))?(a):(b))
#endif

#ifndef MAX
#define MAX(a,b) (((a)>(b))?(a):(b))
#endif

namespace astra { namespace hand {

    struct tracking_data;

    enum class foreground_status
    {
        empty = 0,
        has_points = 1
    };

    namespace segmentation {

        static Point2i INVALID_POINT(-1, -1);

        float get_point_area(tracking_matrices& matrices,
                             area_test_settings& settings,
                             const Point2i& point);

        float get_point_area_integral(tracking_matrices& matrices,
                                      BitmapF& integralArea,
                                      area_test_settings& settings,
                                      const Point2i& point);

        bool test_point_in_range(tracking_matrices& matrices,
                                 const Point2i& targetPoint,
                                 test_behavior outputLog);

        bool test_point_area(tracking_matrices& matrices,
                             area_test_settings& settings,
                             const Point2i& targetPoint,
                             test_phase phase,
                             test_behavior outputLog);

        bool test_point_area_integral(tracking_matrices& matrices,
                                      BitmapF& integralArea,
                                      area_test_settings& settings,
                                      const Point2i& targetPoint,
                                      test_phase phase,
                                      test_behavior outputLog);

        bool test_foreground_radius_percentage(tracking_matrices& matrices,
                                               circumference_test_settings& settings,
                                               const Point2i& targetPoint,
                                               test_phase phase,
                                               test_behavior outputLog);

        foreground_status create_test_pass_from_foreground(tracking_data& data);

        bool find_next_velocity_seed_pixel(BitmapMask& foregroundMatrix,
                                           BitmapMask& searchedMatrix,
                                           Point2i& foregroundPosition,
                                           Point2i& nextSearchStart);

        void calculate_edge_distance(BitmapMask& segmentationMatrix,
                                     BitmapF& areaSqrtMatrix,
                                     BitmapF& edgeDistanceMatrix,
                                     const float maxEdgeDistance);

        float count_neighborhood_area(BitmapMask& matSegmentation,
                                      BitmapF& matDepth,
                                      BitmapF& matArea,
                                      const Point2i& center,
                                      const float bandwidth,
                                      const float bandwidthDepth,
                                      const scaling_coordinate_mapper& mapper);

        BitmapF& calculate_integral_area(tracking_matrices& matrices);

        float count_neighborhood_area_integral(BitmapF& matDepth,
                                               BitmapF& matAreaIntegral,
                                               const Point2i& center,
                                               const float bandwidth,
                                               const scaling_coordinate_mapper& mapper);

        Point2i track_point_impl(tracking_data& data);

        Point2i track_point_from_seed(tracking_data& data);

        void get_circumference_points(BitmapF& matDepth,
                                      const Point2i& center,
                                      const float& radius,
                                      const scaling_coordinate_mapper& mapper,
                                      std::vector<astra::Vector2i>& points);

        float get_max_sequential_circumference_percentage(BitmapF& matDepth,
                                                          BitmapMask& matSegmentation,
                                                          const Point2i& center,
                                                          const float& radius,
                                                          const scaling_coordinate_mapper& mapper,
                                                          std::vector<astra::Vector2i>& points);

        float get_percent_natural_edges(BitmapF& matDepth,
                                        BitmapMask& matSegmentation,
                                        const Point2i& center,
                                        const float bandwidth,
                                        const scaling_coordinate_mapper& mapper);

        bool test_natural_edges(tracking_matrices& matrices,
                                natural_edge_test_settings& settings,
                                const Point2i& targetPoint,
                                test_phase phase,
                                test_behavior outputLog);
    }
}}

#endif // HND_SEGMENTATION_H
