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
#include "mock_device_depth_stream.hpp"

#include <astra/capi/streams/depth_parameters.h>
#include <sstream>
#include <cstring>

namespace orbbec { namespace mocks {

    void populate_conversion_cache(float horizontalFov,
                                   float verticalFov,
                                   int resolutionX,
                                   int resolutionY,
                                   conversion_cache_t& cache)
    {
        cache.xzFactor = std::tan(horizontalFov / 2) * 2;
        cache.yzFactor = std::tan(verticalFov / 2) * 2;
        cache.resolutionX = resolutionX;
        cache.resolutionY = resolutionY;
        cache.halfResX = cache.resolutionX / 2;
        cache.halfResY = cache.resolutionY / 2;
        cache.coeffX = cache.resolutionX / cache.xzFactor;
        cache.coeffY = cache.resolutionY / cache.yzFactor;
    }

    mock_depth_stream::mock_depth_stream()
        : base_stream(astra::devices::streamtype_cast(astra::devices::stream_types::depth))
    {
    }

    mock_depth_stream::~mock_depth_stream()
    {
    }

    device_status mock_depth_stream::on_initialize()
    {
        // Populate modes
        add_available_mode(stream_mode(320, 240, 30, ASTRA_PIXEL_FORMAT_DEPTH_MM));
        add_available_mode(stream_mode(320, 240, 60, ASTRA_PIXEL_FORMAT_DEPTH_MM));

        add_available_mode(stream_mode(640, 480, 30, ASTRA_PIXEL_FORMAT_DEPTH_MM));
        add_available_mode(stream_mode(640, 480, 60, ASTRA_PIXEL_FORMAT_DEPTH_MM));

        set_active_mode(*modes_begin());

        generator_ = astra::make_unique<depth_generator>(active_mode().width(), active_mode().height());

        return device_status_value::ok;
    }

    void mock_depth_stream::refresh_conversion_cache()
    {
        populate_conversion_cache(hFov(), vFov(), active_mode().width(), active_mode().height(), conversionCache_);
    }

    void mock_depth_stream::on_active_mode_changed()
    {
        refresh_conversion_cache();
    }

    device_status mock_depth_stream::on_read_into(void* dest, std::size_t size, std::int32_t timeout)
    {
        int width = active_mode().width();
        int height = active_mode().height();

        std::stringstream ss;
        ss << "depth" << std::endl << active_mode() << std::endl;

        if (flags().has_flags())
        {
            ss << "(" << flags() << ")";
        }

        generator_->set_size(width, height);
        generator_->set_overlayText(ss.str());
        generator_->set_overlayColor(255, 0, 0, 255);
        generator_->generate();

        const std::uint8_t* src = generator_->pixels();
        std::uint16_t* pixelDest = static_cast<std::uint16_t*>(dest);

        for(int y = 0; y < height; y++)
        {
            for(int x = 0; x < width ; x++)
            {
                std::size_t offset = (y * width + x) * 4;
                float r = src[offset];

                std::uint16_t mm = r > 128 ? 2000 : 50000;

                pixelDest[y * width + x] = mm;
            }
        }

        return device_status_value::ok;
    }

    device_status mock_depth_stream::on_get_property_size(property_id id, std::size_t* size)
    {
        switch(id)
        {
        case ASTRA_PARAMETER_DEPTH_REGISTRATION:
            *size = sizeof(bool);
            break;
        case ASTRA_PARAMETER_DEPTH_CONVERSION_CACHE:
            *size = sizeof(conversion_cache_t);
            break;
        default:
            return base_stream::on_get_property_size(id, size);
        }

        return device_status_value::ok;
    }

    device_status mock_depth_stream::on_get_property(property_id id, void* value, std::size_t size)
    {
        switch(id)
        {
        case ASTRA_PARAMETER_DEPTH_REGISTRATION:
        {
            bool* registered = static_cast<bool*>(value);
            *registered = has_flag("registered");
            break;
        }
        case ASTRA_PARAMETER_DEPTH_CONVERSION_CACHE:
            std::memcpy(value, &conversionCache_, size);
            break;

        default:
            return base_stream::on_get_property(id, value, size);
        }

        return device_status_value::ok;
    }

    device_status mock_depth_stream::on_set_property(property_id id, const void* value, std::size_t size)
    {
        switch(id)
        {
        case ASTRA_PARAMETER_DEPTH_REGISTRATION:
        {
            bool registered = *static_cast<const bool*>(value);
            set_flag("registered", registered);

            return device_status_value::ok;
        }
        default:
            return base_stream::on_set_property(id, value, size);
        }
    }

    bool mock_depth_stream::on_property_supported(property_id id) const
    {
        switch(id)
        {
        case ASTRA_PARAMETER_DEPTH_REGISTRATION:
            return true;
        case ASTRA_PARAMETER_DEPTH_CONVERSION_CACHE:
            return true;
        }

        return base_stream::on_property_supported(id);
    }
}}
