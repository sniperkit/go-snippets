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
#ifndef HND_DEBUG_VISUALIZER_H
#define HND_DEBUG_VISUALIZER_H

#include <vector>
#include "hnd_tracked_point.hpp"
#include <astra/capi/streams/stream_types.h>

namespace astra { namespace hand {

    using namespace std;

    class debug_visualizer
    {
    public:
        void overlay_crosshairs(const vector<tracked_point>& points,
                                _astra_imageframe& imageFrame)
        {
            int width = imageFrame.metadata.width;
            int height = imageFrame.metadata.height;

            uint8_t* colorData = static_cast<uint8_t*>(imageFrame.data);

            uint8_t bytesPerPixel;
            astra_pixelformat_get_bytes_per_pixel(imageFrame.metadata.pixelFormat, &bytesPerPixel);

            for (auto iter = points.begin(); iter != points.end(); ++iter)
            {
                tracked_point tracked = *iter;
                Point2i position = tracked.position;

                bool isActivePoint = tracked.pointType == tracked_point_type::active_point;
                bool isLostTrackingPoint = isActivePoint && tracked.trackingStatus == tracking_status::lost;

                int y0 = MAX(0, position.y - 1);
                int y1 = MIN(height - 1, position.y + 1);
                int x0 = MAX(0, position.x - 1);
                int x1 = MIN(width - 1, position.x + 1);

                for (int y = y0; y <= y1; y++)
                {
                    for (int x = x0; x <= x1; x++)
                    {
                        uint8_t r = 0;
                        uint8_t g = 0;
                        uint8_t b = 0;

                        if ((y == position.y || x == position.x))
                        {
                            if (isLostTrackingPoint)
                            {
                                r = 255;
                                g = 0;
                                b = 0;
                            }
                            else if (isActivePoint)
                            {
                                r = 0;
                                g = 139;
                                b = 69;
                            }
                            else
                            {
                                r = 255;
                                g = 255;
                                b = 0;
                            }
                        }
                        else
                        {
                            r = 0;
                            g = 0;
                            b = 0;
                        }

                        int index = x + y * width;
                        uint8_t* pixel = colorData + index * bytesPerPixel;

                        *(pixel) = r;
                        *(pixel + 1) = g;
                        *(pixel + 2) = b;
                    }
                }
            }
        }

        void overlay_mask(const BitmapMask& matMask,
                          _astra_imageframe& imageFrame,
                          const RgbPixel& maskColor,
                          const pixel_type targetValue)
        {
            assert(matMask.width() == imageFrame.metadata.width);
            assert(matMask.height() == imageFrame.metadata.height);

            int width = matMask.width();
            int height = matMask.height();

            RgbPixel* colorData = static_cast<RgbPixel*>(imageFrame.data);

            for (int y = 0; y < height; ++y)
            {
                const auto* maskRow = matMask.data(y);

                for (int x = 0; x < width; ++x, ++maskRow, ++colorData)
                {
                    pixel_type maskValue = static_cast<pixel_type>(*maskRow);

                    if (maskValue == targetValue)
                    {
                        *colorData = maskColor;
                    }
                }
            }
        }

        void show_depth_matrix(const BitmapF& matDepth,
                               _astra_imageframe& imageFrame)
        {
            assert(matDepth.width() == imageFrame.metadata.width);
            assert(matDepth.height() == imageFrame.metadata.height);

            int width = matDepth.width();
            int height = matDepth.height();

            RgbPixel* colorData = static_cast<RgbPixel*>(imageFrame.data);

            for (int y = 0; y < height; ++y)
            {
                const float* depthRow = matDepth.data(y);

                for (int x = 0; x < width; ++x, ++depthRow, ++colorData)
                {
                    float depth = *depthRow;
                    float normDepth = std::min(1.0f, std::max(0.0f, (depth - 400.0f) / 5600.0f));
                    uint8_t value = 255 * (1 - normDepth);
                    if (depth == 0)
                    {
                        value = 0;
                    }
                    RgbPixel color(0, value, value);

                    *colorData = color;
                }
            }
        }

        void show_velocity_matrix(const BitmapF& matVelocity,
                                  float maxScale,
                                  _astra_imageframe& imageFrame)
        {
            assert(matVelocity.width() == imageFrame.metadata.width);
            assert(matVelocity.height() == imageFrame.metadata.height);

            if (maxScale == 0)
            {
                throw new std::invalid_argument("maxScale cannot be 0");
            }

            const int width = matVelocity.width();
            const int height = matVelocity.height();

            RgbPixel* colorData = static_cast<RgbPixel*>(imageFrame.data);

            for (int y = 0; y < height; ++y)
            {
                const float* velocityRow = matVelocity.data(y);
                for (int x = 0; x < width; ++x, ++velocityRow, ++colorData)
                {
                    float velocity = *velocityRow;

                    uint8_t velocityValue = static_cast<uint8_t>(255 * sqrt(min(1.0f, abs(velocity / maxScale))));

                    RgbPixel color(0, velocityValue, 0);
                    if (velocity < 0)
                    {
                        color.r = velocityValue;
                        color.g = velocityValue;
                    }

                    *colorData = color;
                }
            }
        }

        template<typename T>
        void show_norm_array(const Bitmap<T>& mat,
                             const BitmapMask& mask,
                             _astra_imageframe& imageFrame)
        {
            assert(mat.width() == imageFrame.metadata.width);
            assert(mat.height() == imageFrame.metadata.height);

            int width = mat.width();
            int height = mat.height();

            MinMaxLoc<T> minMaxLoc;

            bool emptyMask = all_zero(mask);

            if (!emptyMask)
            {
                assert(mat.size() == mask.size());
                minMaxLoc = find_min_max_loc(mat, mask);
            }
            else
            {
                minMaxLoc = find_min_max_loc(mat);
            }

            double range = minMaxLoc.max - minMaxLoc.min;
            bool rangeZero = abs(range) < 0.00001;

            uint8_t* colorData = static_cast<uint8_t*>(imageFrame.data);

            uint8_t bytesPerPixel;
            astra_pixelformat_get_bytes_per_pixel(imageFrame.metadata.pixelFormat, &bytesPerPixel);

            for (int y = 0; y < height; ++y)
            {
                const T* dataRow = mat.data(y);
                const uint8_t* maskRow = nullptr;
                if (!emptyMask)
                {
                    maskRow = mask.data(y);
                }

                for (int x = 0; x < width; ++x, ++dataRow, colorData += bytesPerPixel)
                {
                    float data = *dataRow;
                    uint8_t maskValue = 1;
                    if (!emptyMask)
                    {
                        maskValue = *maskRow;
                        ++maskRow;
                    }
                    float value;
                    if (0 == data || 0 == maskValue)
                    {
                        value = 0;
                    }
                    else if (rangeZero)
                    {
                        value = 255;
                    }
                    else
                    {
                        value = 55 + 200 * ((*dataRow - minMaxLoc.min) / range);
                    }
                    *(colorData) = value;
                    *(colorData + 1) = value;
                    *(colorData + 2) = value;
                }
            }
        }
    };
}}

#endif // HND_DEBUGVISUALIZER_H
