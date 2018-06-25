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
#include "mock_device_streamset.hpp"
#include "mock_device.hpp"
#include "astra_sensor.hpp"
#include "mock_device_stream.hpp"
#include <chrono>

namespace orbbec { namespace mocks {

    device_streamset::device_streamset(const std::string& name,
                                       astra::PluginServiceProxy& pluginService,
                                       const char* uri)
        : pluginService_(pluginService),
          uri_(uri)
    {
        uri_ = uri;

        pluginService_.create_stream_set(name.c_str(), streamSetHandle_);

        device_ = std::make_shared<mock_device>();
        device_->initialize();
    }

    device_streamset::~device_streamset()
    {
        close();
        streams_.clear();

        LOG_DEBUG("orbbec.mocks.device_streamset", "destroying streamset for device %s", uri_.c_str());

        pluginService_.destroy_stream_set(streamSetHandle_);
    }

    astra_status_t device_streamset::open()
    {
        if (isOpen_)
            return astra_status_t::ASTRA_STATUS_SUCCESS;

        LOG_DEBUG("orbbec.mocks.device_streamset", "opening device: %s", uri_.c_str());

        open_sensor_streams();

        LOG_DEBUG("orbbec.mocks.device_streamset", "opened device: %s", uri_.c_str());

        isOpen_ = true;

        return astra_status_t::ASTRA_STATUS_SUCCESS;
    }

    astra_status_t device_streamset::close()
    {
        if (!isOpen_)
            return astra_status_t::ASTRA_STATUS_SUCCESS;

        device_->disconnect();

        LOG_INFO("orbbec.mocks.device_streamset", "closing device: %s", uri_.c_str());

        isOpen_ = false;

        return astra_status_t::ASTRA_STATUS_SUCCESS;
    }

    astra_status_t device_streamset::read()
    {
        if (!isOpen_)
            return astra_status_t::ASTRA_STATUS_SUCCESS;

        frameIndex_++;

        device_->update();

        for(auto& stream : streams_)
        {
            stream->read(frameIndex_);
        }

        return astra_status_t::ASTRA_STATUS_SUCCESS;
    }

    void device_streamset::add_stream(mock_stream* stream)
    {
        streams_.push_back(stream_ptr(stream));
    }

    astra_status_t device_streamset::open_sensor_streams()
    {
        bool enableColor = true;

        if (enableColor)
        {
            auto it = device_->find_sensor(astra::devices::sensor_type::color);
            if (it != device_->sensors_end())
            {
                auto it2 = (*it)->find_stream(
                    astra::devices::streamtype_cast<>(astra::devices::stream_types::color));

                if (it2 != (*it)->streams_end())
                {
                    image_stream* stream = astra::plugins::make_stream<image_stream>(pluginService_,
                                                                                     streamSetHandle_,
                                                                                     astra::StreamDescription(
                                                                                         ASTRA_STREAM_COLOR,
                                                                                         DEFAULT_SUBTYPE),
                                                                                     *this,
                                                                                     *it2);

                    astra_status_t rc = astra_status_t::ASTRA_STATUS_SUCCESS;
                    rc = stream->open();
                    add_stream(stream);

                    if (rc != astra_status_t::ASTRA_STATUS_SUCCESS)
                        LOG_WARN("orbbec.mocks.device_streamset", "unable to open mock color stream.");
                }
            }
        }

        bool enableDepth = true;

        if (enableDepth)
        {
            auto it = device_->find_sensor(astra::devices::sensor_type::depth);
            if (it != device_->sensors_end())
            {
                auto it2 = (*it)->find_stream(
                    astra::devices::streamtype_cast<>(astra::devices::stream_types::depth));

                if (it2 != (*it)->streams_end())
                {
                    image_stream* stream = astra::plugins::make_stream<image_stream>(pluginService_,
                                                                                     streamSetHandle_,
                                                                                     astra::StreamDescription(
                                                                                         ASTRA_STREAM_DEPTH,
                                                                                         DEFAULT_SUBTYPE),
                                                                                     *this,
                                                                                     *it2);

                    astra_status_t rc = astra_status_t::ASTRA_STATUS_SUCCESS;
                    rc = stream->open();
                    add_stream(stream);

                    if (rc != astra_status_t::ASTRA_STATUS_SUCCESS)
                        LOG_WARN("orbbec.mocks.device_streamset", "unable to open mock depth stream.");
                }
            }
        }

        bool enableInfrared = true;

        if (enableInfrared)
        {
            auto it = device_->find_sensor(astra::devices::sensor_type::depth);
            if (it != device_->sensors_end())
            {
                auto it2 = (*it)->find_stream(
                    astra::devices::streamtype_cast<>(astra::devices::stream_types::infrared));

                if (it2 != (*it)->streams_end())
                {
                    image_stream* stream = astra::plugins::make_stream<image_stream>(pluginService_,
                                                                                     streamSetHandle_,
                                                                                     astra::StreamDescription(
                                                                                         ASTRA_STREAM_INFRARED,
                                                                                         DEFAULT_SUBTYPE),
                                                                                     *this,
                                                                                     *it2);

                    astra_status_t rc = astra_status_t::ASTRA_STATUS_SUCCESS;
                    rc = stream->open();
                    add_stream(stream);

                    if (rc != astra_status_t::ASTRA_STATUS_SUCCESS)
                        LOG_WARN("orbbec.mocks.device_streamset", "unable to open mock infrared stream.");
                }
            }
        }

        return astra_status_t::ASTRA_STATUS_SUCCESS;
    }

    void device_streamset::on_started(mock_stream* stream)
    {
        LOG_INFO("orbbec.mocks.device_streamset",
                 "adding stream type %u to active streams",
                 stream->description().type());
    }

    void device_streamset::on_stopped(mock_stream* stream)
    {
        LOG_INFO("orbbec.mocks.device_streamset",
                 "removing stream type %u from active streams",
                 stream->description().type());
    }

    astra_status_t device_streamset::close_sensor_streams()
    {
        streams_.clear();

        return astra_status_t::ASTRA_STATUS_SUCCESS;
    }
}}
