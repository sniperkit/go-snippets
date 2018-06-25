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
#include "mock_base_stream.hpp"
#include <astra/capi/streams/image_parameters.h>
#include <cstring>

namespace orbbec { namespace mocks {

    base_stream::base_stream(const stream_type& type)
        : sensor_stream(type)
    {
        hFov_ = 60 * 3.141519 / 180.0f;
        vFov_ = 45 * 3.141519 / 180.0f;
    }

    base_stream::~base_stream() = default;

    float base_stream:: hFov() const
    {
        return hFov_;
    }

    float base_stream::vFov() const
    {
        return vFov_;
    }

    void base_stream::set_vFov(float vFov)
    {
        vFov_ = vFov;
    }

    void base_stream::set_hFov(float hFov)
    {
        hFov_ = hFov;
    }

    device_status base_stream::on_start()
    {
        timer_.set_interval(1.0 / (double)active_mode().fps());

        timer_.start();

        return device_status_value::ok;
    }

    device_status base_stream::on_stop()
    {
        timer_.stop();

        return device_status_value::ok;
    }

    void base_stream::on_update()
    {
        if (!is_streaming())
            return;

        if (timer_.is_triggered())
        {
            timer_.reset();
            raise_new_frame_available();
        }
    }

    device_status base_stream::on_get_property_size(property_id id, std::size_t* size)
    {
        switch(id)
        {
        case ASTRA_PARAMETER_IMAGE_MIRRORING:
            *size = sizeof(bool);
            break;
        case ASTRA_PARAMETER_IMAGE_HFOV:
            *size = sizeof(float);
            break;
        case ASTRA_PARAMETER_IMAGE_VFOV:
            *size =  sizeof(float);
            break;
        case ASTRA_PARAMETER_IMAGE_MODE:
            *size = sizeof(astra::ImageStreamMode);
            break;
        case ASTRA_PARAMETER_IMAGE_MODES:
            *size = sizeof(astra::ImageStreamMode) * modes_size();
            break;
        default:
            return device_status_value::invalid_parameter;
        }

        return device_status_value::ok;
    }

    device_status base_stream::on_get_property(property_id id, void* value, std::size_t size)
    {
        switch(id)
        {
        case ASTRA_PARAMETER_IMAGE_MIRRORING:
        {
            bool* registered = static_cast<bool*>(value);
            *registered = has_flag("mirroring");

            break;
        }
        case ASTRA_PARAMETER_IMAGE_HFOV:
        {
            float* hFov = reinterpret_cast<float*>(value);
            *hFov = this->hFov();

            break;
        }
        case ASTRA_PARAMETER_IMAGE_VFOV:
        {
            float* vFov = reinterpret_cast<float*>(value);
            *vFov = this->vFov();

            break;
        }
        case ASTRA_PARAMETER_IMAGE_MODE:
        {
            stream_mode* mode = reinterpret_cast<stream_mode*>(value);
            *mode = this->active_mode();

            break;
        }
        case ASTRA_PARAMETER_IMAGE_MODES:
        {
            assert(size >= sizeof(astra::ImageStreamMode) * modes_size());
            std::memcpy(value, &(*modes_begin()), size);

            break;
        }
        default:
            return device_status_value::invalid_parameter;
        }

        return device_status_value::ok;
    }

    device_status base_stream::on_set_property(property_id id, const void* value, std::size_t size)
    {
        switch(id)
        {
        case ASTRA_PARAMETER_IMAGE_MIRRORING:
        {
            bool mirroring = *static_cast<const bool*>(value);
            set_flag("mirroring", mirroring);

            break;
        }
        case ASTRA_PARAMETER_IMAGE_MODE:
        {
            const astra::ImageStreamMode mode = *static_cast<const astra::ImageStreamMode*>(value);
            timer_.set_interval(1.0 / mode.fps());
            set_active_mode(mode);

            break;
        }
        default:
            return device_status_value::invalid_parameter;
        }

        return device_status_value::ok;
    }

    bool base_stream::on_property_supported(property_id id) const
    {
        switch(id)
        {
        case ASTRA_PARAMETER_IMAGE_MIRRORING:
        case ASTRA_PARAMETER_IMAGE_MODES:
        case ASTRA_PARAMETER_IMAGE_MODE:
        case ASTRA_PARAMETER_IMAGE_HFOV:
        case ASTRA_PARAMETER_IMAGE_VFOV:
            break;
        default:
            return false;
        }

        return true;
    }

    bool base_stream::has_flag(const std::string& flag)
    {
        return flags_.has_flag(flag);
    }

    void base_stream::set_flag(const std::string& flag, bool value)
    {
        flags_.set_flag(flag, value);
    }

    const flag_set& base_stream::flags() const
    {
        return flags_;
    }
}}
