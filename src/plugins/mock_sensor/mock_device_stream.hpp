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
#ifndef MOCK_DEVICE_STREAM_H
#define MOCK_DEVICE_STREAM_H

#include <astra/streams/Image.hpp>
#include <astra/capi/streams/image_parameters.h>
#include <astra/capi/streams/image_types.h>
#include <astra/capi/streams/image_capi.h>
#include <astra/capi/streams/stream_types.h>

#include <chrono>
#include <memory>
#include <cstring>
#include <vector>
#include <sstream>

#include "mock_stream.hpp"
#include "mock_flag_set.hpp"
#include "astra_sensor_stream.hpp"
#include "astra_stream_listener.hpp"

namespace orbbec { namespace mocks {

    template<typename TFrameWrapper>
    class device_stream : public mock_stream,
                          public astra::devices::stream_listener
    {
    public:
        using wrapper_type = TFrameWrapper;

        device_stream(astra::PluginServiceProxy& pluginService,
                      astra_streamset_t streamSet,
                      astra::StreamDescription desc,
                      orbbec::mocks::stream_listener& listener,
                      astra::devices::sensor_stream::shared_ptr stream);

        virtual ~device_stream();

        bool is_streaming() const { return is_open() && is_started(); }

        virtual void new_frame_available(astra::devices::sensor_stream::shared_ptr sensor) override;

        virtual void property_changed(astra::devices::sensor_stream::shared_ptr stream,
                                      astra::devices::sensor_stream::property_id id) override;

    protected:
        virtual void on_get_parameter(astra_streamconnection_t connection,
                                      astra_parameter_id id,
                                      astra_parameter_bin_t& parameterBin) override;

        virtual void on_set_parameter(astra_streamconnection_t connection,
                                      astra_parameter_id id,
                                      size_t inByteLength,
                                      astra_parameter_data_t inData) override;

        virtual void on_connection_added(astra_streamconnection_t connection) override;
        virtual void on_connection_removed(astra_bin_t bin,
                                           astra_streamconnection_t connection) override;

        virtual void on_connection_started(astra_streamconnection_t connection) override;
        virtual void on_connection_stopped(astra_streamconnection_t connection) override;

        virtual astra_status_t on_read(astra_frame_index_t frameIndex) override final;

        virtual astra_status_t on_open() override;
        virtual astra_status_t on_close() override;
        virtual astra_status_t on_start() override;
        virtual astra_status_t on_stop() override;

        virtual void on_new_buffer(wrapper_type* wrapper);

        void change_mode(const astra::ImageStreamMode& mode);
        void set_mode(const astra::ImageStreamMode& mode);

        astra::devices::sensor_stream::shared_ptr deviceStream_;

        using bin_type = astra::plugins::stream_bin<wrapper_type>;
        std::unique_ptr<bin_type> bin_;

        size_t bufferLength_{0};

    private:
        bool frameAvailable_{false};
        astra_stream_t streamHandle_{nullptr};
        std::vector<astra_streamconnection_t> connections_;
    };
}}

namespace orbbec { namespace mocks {

    template<typename TFrameWrapper>
    device_stream<TFrameWrapper>::device_stream(astra::PluginServiceProxy& pluginService,
                                                astra_streamset_t streamSet,
                                                astra::StreamDescription desc,
                                                orbbec::mocks::stream_listener& listener,
                                                astra::devices::sensor_stream::shared_ptr stream)
        : orbbec::mocks::mock_stream(pluginService,
                                     streamSet,
                                     desc,
                                     listener),
          deviceStream_(stream)
    {
        deviceStream_->add_listener(this);
    }

    template<typename TFrameWrapper>
    device_stream<TFrameWrapper>::~device_stream()
    {
        close();
    }

    template<typename TFrameWrapper>
    void device_stream<TFrameWrapper>::on_connection_added(astra_streamconnection_t connection)
    {
        auto it = std::find(connections_.begin(), connections_.end(), connection);

        if (it == connections_.end())
        {
            if (bin_ != nullptr)
                bin_->link_connection(connection);
            connections_.push_back(connection);
        }
    }

    template<typename TFrameWrapper>
    void device_stream<TFrameWrapper>::on_connection_removed(astra_bin_t bin,
                                                             astra_streamconnection_t connection)
    {
        auto it = std::find(connections_.begin(), connections_.end(), connection);

        if (it != connections_.end())
        {
            bin_->unlink_connection(connection);
            connections_.erase(it);
        }

        if (!bin_->has_connections())
        {
            bin_ = nullptr;
        }
    }

    template<typename TFrameWrapper>
    astra_status_t device_stream<TFrameWrapper>::on_read(astra_frame_index_t frameIndex)
    {
        if (frameAvailable_)
        {
            wrapper_type* wrapper = bin_->begin_write(frameIndex);

            wrapper->frame.frame = nullptr;
            wrapper->frame.data = &(wrapper->frame_data);

            on_new_buffer(wrapper);

            std::uint8_t* dest = static_cast<std::uint8_t*>(wrapper->frame.data);

            std::uint32_t width = wrapper->frame.metadata.width;
            std::uint32_t height = wrapper->frame.metadata.height;

            uint8_t bpp;
            astra_pixelformat_get_bytes_per_pixel(wrapper->frame.metadata.pixelFormat, &bpp);
            std::size_t bufferSize = width * height * bpp;

            deviceStream_->read_into(dest, bufferSize, 0);

            bin_->end_write();

            frameAvailable_ = false;
        }

        return ASTRA_STATUS_SUCCESS;
    }

    template<typename TFrameWrapper>
    void device_stream<TFrameWrapper>::on_connection_started(astra_streamconnection_t connection)
    {
        LOG_INFO("orbbec.mocks.device_stream", "client started: %u", description().type());
        if (is_open() && !is_started())
        {
            start();
        }
    }

    template<typename TFrameWrapper>
    void device_stream<TFrameWrapper>::on_connection_stopped(astra_streamconnection_t connection)
    {
        LOG_INFO("orbbec.mocks.device_stream", "client stopped: %u", description().type());
        if (is_open() && is_started())
        {
            stop();
        }
    }

    template<typename TFrameWrapper>
    void device_stream<TFrameWrapper>::on_get_parameter(astra_streamconnection_t connection,
                                                        astra_parameter_id id,
                                                        astra_parameter_bin_t& parameterBin)
    {
        size_t resultByteLength;
        deviceStream_->get_property_size(id, &resultByteLength);

        astra_parameter_data_t parameterData;
        astra_status_t rc = pluginService().get_parameter_bin(resultByteLength,
                                                              &parameterBin,
                                                              &parameterData);
        if (rc == astra_status_t::ASTRA_STATUS_SUCCESS)
        {
            deviceStream_->get_property(id,
                                        parameterData, resultByteLength);
        }
    }

    template<typename TFrameWrapper>
    void device_stream<TFrameWrapper>::on_set_parameter(astra_streamconnection_t connection,
                                                        astra_parameter_id id,
                                                        size_t inByteLength,
                                                        astra_parameter_data_t inData)
    {
        deviceStream_->set_property(id, inData, inByteLength);
    }

    template<typename TFrameWrapper>
    void device_stream<TFrameWrapper>::set_mode(const astra::ImageStreamMode& mode)
    {
        assert(mode.pixel_format() != 0);

        bufferLength_ =
            mode.width() *
            mode.height() *
            mode.bytes_per_pixel();

        LOG_INFO("orbbec.mocks.device_stream", "bin change: %ux%ux%u len: %u",
                 mode.width(),
                 mode.height(),
                 mode.bytes_per_pixel(),
                 bufferLength_);

        bin_ = astra::make_unique<bin_type>(pluginService(),
                                          get_handle(),
                                          bufferLength_);

        for(auto& conn : connections_)
        {
            bin_->link_connection(conn);
        };
    }

    template<typename TFrameWrapper>
    void device_stream<TFrameWrapper>::on_new_buffer(wrapper_type* wrapper)
    {
        if (!wrapper)
            return;

        auto& md = wrapper->frame.metadata;

        const auto& mode = deviceStream_->active_mode();

        md.pixelFormat = mode.pixel_format();
        md.width = mode.width();
        md.height = mode.height();
    }

    template<typename TFrameWrapper>
    astra_status_t device_stream<TFrameWrapper>::on_open()
    {
        if (is_open())
            return astra_status_t::ASTRA_STATUS_SUCCESS;

        LOG_INFO("orbbec.mocks.device_stream",
                 "creating mock stream of type: %d",
                 description().type());

        LOG_INFO("orbbec.mocks.device_stream",
                 "created mock stream of type: %d",
                 description().type());

        LOG_INFO("orbbec.mocks.device_stream",
                 "stream type %d supports modes:",
                 description().type());

        for(auto it = deviceStream_->modes_begin();
            it != deviceStream_->modes_end(); ++it)
        {
            auto mode = *it;

            LOG_INFO("orbbec.mocks.device_stream", "mode: %ux%ux%u@%u pf:%u",
                     mode.width(),
                     mode.height(),
                     mode.bytes_per_pixel(),
                     mode.fps(),
                     mode.pixel_format());
        }

        set_mode(deviceStream_->active_mode());

        return astra_status_t::ASTRA_STATUS_SUCCESS;
    }

    template<typename TFrameWrapper>
    astra_status_t device_stream<TFrameWrapper>::on_close()
    {
        stop();

        LOG_INFO("orbbec.mocks.device_stream",
                 "destroying mock stream of type: %d",
                 description().type());

        return astra_status_t::ASTRA_STATUS_SUCCESS;
    }

    template<typename TFrameWrapper>
    astra_status_t device_stream<TFrameWrapper>::on_start()
    {
        LOG_INFO("orbbec.mocks.device_stream",
                 "starting mock stream of type: %d",
                 description().type());

        deviceStream_->start();

        return astra_status_t::ASTRA_STATUS_SUCCESS;
    }

    template<typename TFrameWrapper>
    astra_status_t device_stream<TFrameWrapper>::on_stop()
    {
        LOG_INFO("orbbec.mocks.device_stream",
                 "stopping mock stream of type: %d",
                 description().type());

        deviceStream_->stop();

        LOG_INFO("orbbec.mocks.device_stream",
                 "stopped mock stream of type: %d",
                 description().type());

        return astra_status_t::ASTRA_STATUS_SUCCESS;
    }

    template<typename TFrameWrapper>
    void device_stream<TFrameWrapper>::new_frame_available(astra::devices::sensor_stream::shared_ptr stream)
    {

        frameAvailable_ = true;
    }

    template<typename TFrameWrapper>
    void device_stream<TFrameWrapper>::property_changed(astra::devices::sensor_stream::shared_ptr stream,
                                                        astra::devices::sensor_stream::property_id id)
    {
        if (id == ASTRA_PARAMETER_IMAGE_MODE)
        {
            set_mode(stream->active_mode());
        }
    }

    class image_stream : public device_stream<::astra_imageframe_wrapper_t>
    {
    public:
        image_stream(astra::PluginServiceProxy& pluginService,
                     astra_streamset_t streamSet,
                     astra::StreamDescription desc,
                     orbbec::mocks::stream_listener& listener,
                     astra::devices::sensor_stream::shared_ptr deviceStream)
            : device_stream(pluginService, streamSet, desc, listener, deviceStream)
        {
        }
    };

}}

#endif /* MOCK_DEVICE_STREAM_H */
