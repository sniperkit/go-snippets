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
#ifndef ONI_DEVICESTREAM_H
#define ONI_DEVICESTREAM_H

#include <Shiny.h>
#include <OpenNI.h>
#include <astra/streams/Image.hpp>
#include <astra/capi/streams/image_parameters.h>
#include <astra/capi/streams/image_types.h>
#include <astra/capi/streams/image_capi.h>
#include <memory>
#include <cstring>
#include <vector>
#include "oni_mappers.hpp"
#include "oni_stream.hpp"

#ifndef MIN
#define MIN(a,b) (((a)<(b))?(a):(b))
#endif

#ifndef MAX
#define MAX(a,b) (((a)>(b))?(a):(b))
#endif

namespace orbbec { namespace ni {

    template<typename TFrameWrapper>
    class devicestream : public stream
    {
    public:
        using wrapper_type = TFrameWrapper;

        devicestream(astra::PluginServiceProxy& pluginService,
                     astra_streamset_t streamSet,
                     astra::StreamDescription desc,
                     openni::Device& oniDevice,
                     openni::SensorType oniSensorType,
                     stream_listener& listener)
            : stream(pluginService,
                     streamSet,
                     desc,
                     listener),
              oniDevice_(oniDevice),
              oniSensorType_(oniSensorType)
        {
            PROFILE_FUNC();
        }

        virtual ~devicestream()
        {
            PROFILE_FUNC();
            close();
        }

        inline bool is_streaming() const { return is_open() && is_started(); }

        virtual void on_connection_started(astra_streamconnection_t connection) override
        {
            LOG_INFO("orbbec.ni.devicestream", "turn on stream %u", description().type());
            if (is_open() && !is_started())
            {
                start();
            }
        }

        virtual void on_connection_stopped(astra_streamconnection_t connection) override
        {
            if (is_open() && is_started())
            {
                stop();
            }
        }

        virtual void on_get_parameter(astra_streamconnection_t connection,
                                      astra_parameter_id id,
                                      astra_parameter_bin_t& parameterBin) override
        {
            PROFILE_FUNC();

            switch (id)
            {
            case ASTRA_PARAMETER_IMAGE_HFOV:
            {
                size_t resultByteLength = sizeof(float);

                astra_parameter_data_t parameterData;
                astra_status_t rc = pluginService().get_parameter_bin(resultByteLength,
                                                                      &parameterBin,
                                                                      &parameterData);
                if (rc == ASTRA_STATUS_SUCCESS)
                {
                    float* hFov = reinterpret_cast<float*>(parameterData);
                    *hFov = oniStream_.getHorizontalFieldOfView();
                }
                break;
            }
            case ASTRA_PARAMETER_IMAGE_VFOV:
            {
                size_t resultByteLength = sizeof(float);

                astra_parameter_data_t parameterData;
                astra_status_t rc = pluginService().get_parameter_bin(resultByteLength,
                                                                      &parameterBin,
                                                                      &parameterData);
                if (rc == ASTRA_STATUS_SUCCESS)
                {
                    float* vFov = reinterpret_cast<float*>(parameterData);
                    *vFov = oniStream_.getVerticalFieldOfView();
                }
                break;
            }
            case ASTRA_PARAMETER_IMAGE_MIRRORING:
            {
                size_t resultByteLength = sizeof(bool);

                astra_parameter_data_t parameterData;
                astra_status_t rc = pluginService().get_parameter_bin(resultByteLength,
                                                                      &parameterBin,
                                                                      &parameterData);
                if (rc == ASTRA_STATUS_SUCCESS)
                {
                    bool mirroring = oniStream_.getMirroringEnabled();
                    *reinterpret_cast<bool*>(parameterData) = mirroring;
                }

                break;
            }
            case ASTRA_PARAMETER_IMAGE_MODES:
            {
                std::size_t resultSize = sizeof(astra_imagestream_mode_t) * modes_.size();

                astra_parameter_data_t parameterData;
                astra_status_t rc = pluginService().get_parameter_bin(resultSize,
                                                                      &parameterBin,
                                                                      &parameterData);

                astra::ImageStreamMode* result = static_cast<astra::ImageStreamMode*>(parameterData);

                if (rc == ASTRA_STATUS_SUCCESS)
                {
                    for(int i = 0; i < modes_.size(); i++)
                    {
                        result[i] = modes_[i];
                    }
                }

                break;
            }
            }
        }

        void on_set_parameter(astra_streamconnection_t connection,
                              astra_parameter_id id,
                              size_t inByteLength,
                              astra_parameter_data_t inData) override
        {
            switch (id)
            {
            case ASTRA_PARAMETER_IMAGE_MIRRORING:
            {
                bool enable = *static_cast<bool*>(inData);
                oniStream_.setMirroringEnabled(enable);
                break;
            }
            case ASTRA_PARAMETER_IMAGE_MODE:
                astra::ImageStreamMode mode = *static_cast<astra::ImageStreamMode*>(inData);
                auto oniMode = convert_mode(mode);

                LOG_INFO("orbbec.ni.devicestream", "mode change requested: %ux%ux%u@%u pf:%u",
                         mode.width(),
                         mode.height(),
                         mode.bytes_per_pixel(),
                         mode.fps(),
                         mode.pixel_format());

                const bool wasStarted = is_started();
                if (wasStarted)
                {
                    stop();
                }
                auto rc = oniStream_.setVideoMode(oniMode);
                if (rc == openni::Status::STATUS_OK)
                {
                    LOG_INFO("orbbec.ni.devicestream", "stream mode changed");
                }
                else
                {
                    LOG_WARN("orbbec.ni.devicestream", "failed to change stream mode.");
                }

                if (wasStarted)
                {
                    start();
                }
                break;
            }
        }

        void change_mode(const astra::ImageStreamMode& mode)
        {
            mode_ = mode;
            oniMode_ = convert_mode(mode);
            assert(mode_.pixel_format() != 0);

            bufferLength_ =
                mode_.width() *
                mode_.height() *
                mode_.bytes_per_pixel();

            LOG_INFO("orbbec.ni.devicestream", "bin swap: %ux%ux%u len: %u",
                     mode_.width(),
                     mode_.height(),
                     mode_.bytes_per_pixel(),
                     bufferLength_);

            bin_ = astra::make_unique<bin_type>(pluginService(),
                                              get_handle(),
                                              bufferLength_);

            for(auto& conn : connections_)
            {
                bin_->link_connection(conn);
            };
        }

        virtual void on_new_buffer(wrapper_type* wrapper)
        {
            PROFILE_FUNC();

            if (!wrapper)
                return;

            auto& md = wrapper->frame.metadata;

            md.pixelFormat = mode_.pixel_format();
            md.width = mode_.width();
            md.height = mode_.height();
        }

        virtual astra_status_t on_read(astra_frame_index_t frameIndex) override;

        virtual openni::VideoStream* get_stream() override { return &oniStream_; }

        virtual void on_connection_added(astra_streamconnection_t connection) override;
        virtual void on_connection_removed(astra_bin_t bin,
                                           astra_streamconnection_t connection) override;

    protected:
        openni::Device& oniDevice_;
        openni::SensorType oniSensorType_;
        openni::VideoStream oniStream_;
        openni::VideoMode oniMode_;
        astra::ImageStreamMode mode_;

        virtual astra_status_t on_open() override
        {
            PROFILE_FUNC();
            if (is_open())
                return ASTRA_STATUS_SUCCESS;

            LOG_INFO("orbbec.ni.devicestream", "creating oni stream of type: %d", description().type());
            openni::Status rc = oniStream_.create(oniDevice_, oniSensorType_);

            if (rc != openni::STATUS_OK)
            {
                return ASTRA_STATUS_DEVICE_ERROR;
            }

            LOG_INFO("orbbec.ni.devicestream", "created oni stream of type: %d", description().type());

            const openni::SensorInfo& pInfo = oniStream_.getSensorInfo();
            auto& modes = pInfo.getSupportedVideoModes();

            LOG_INFO("orbbec.ni.devicestream", "stream type %d supports modes:", description().type());

            for(int i = 0; i < modes.getSize(); i++)
            {
                const openni::VideoMode& oniMode = modes[i];

                if (std::get<0>(convert_format(oniMode.getPixelFormat())) != 0)
                {
                    astra::ImageStreamMode mode = convert_mode(oniMode);
                    modes_.push_back(mode);
                }

                LOG_INFO("orbbec.ni.devicestream", "- w: %d h: %d fps: %d pf: %d",
                         oniMode.getResolutionX(),
                         oniMode.getResolutionY(),
                         oniMode.getFps(),
                         oniMode.getPixelFormat());
            }

            auto oniVideoMode = oniStream_.getVideoMode();
            change_mode(convert_mode(oniVideoMode));

            return ASTRA_STATUS_SUCCESS;
        }

        virtual astra_status_t on_close() override
        {
            PROFILE_FUNC();

            stop();

            LOG_INFO("orbbec.ni.devicestream", "destroying oni stream of type: %d", description().type());
            oniStream_.destroy();

            return ASTRA_STATUS_SUCCESS;
        }

        virtual astra_status_t on_start() override
        {
            PROFILE_FUNC();

            LOG_INFO("orbbec.ni.devicestream", "starting oni stream of type: %d", description().type());

            auto rc = oniStream_.start();

            if (rc == openni::Status::STATUS_OK)
            {
                LOG_INFO("orbbec.ni.devicestream",
                         "started oni stream of type: %d",
                         description().type());

                return ASTRA_STATUS_SUCCESS;
            }
            else
            {
                LOG_INFO("orbbec.ni.devicestream",
                         "unable to start oni stream of type: %d",
                         description().type());

                return ASTRA_STATUS_DEVICE_ERROR;
            }
        }

        virtual astra_status_t on_stop() override
        {
            PROFILE_FUNC();

            LOG_INFO("orbbec.ni.devicestream", "stopping oni stream of type: %d", description().type());
            oniStream_.stop();
            LOG_INFO("orbbec.ni.devicestream", "stopped oni stream of type: %d", description().type());

            return ASTRA_STATUS_SUCCESS;
        }

    private:
        using bin_type = astra::plugins::stream_bin<wrapper_type>;
        std::unique_ptr<bin_type> bin_;

        std::vector<astra_streamconnection_t> connections_;

        size_t bufferLength_{0};
        astra_stream_t streamHandle_{nullptr};

        std::vector<astra::ImageStreamMode> modes_;
    };

    template<typename TFrameWrapper>
    void devicestream<TFrameWrapper>::on_connection_added(astra_streamconnection_t connection)
    {
        PROFILE_FUNC();

        auto it = std::find(connections_.begin(), connections_.end(), connection);

        if (it == connections_.end())
        {
            if (bin_ != nullptr)
                bin_->link_connection(connection);
            connections_.push_back(connection);
        }
    }

    template<typename TFrameWrapper>
    void devicestream<TFrameWrapper>::on_connection_removed(astra_bin_t bin,
                                                            astra_streamconnection_t connection)
    {
        PROFILE_FUNC();

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

    inline bool operator==(const openni::VideoMode& lhs, const openni::VideoMode& rhs)
    {
        return lhs.getResolutionX() == rhs.getResolutionX()
            && lhs.getResolutionY() == rhs.getResolutionY()
            && lhs.getPixelFormat() == rhs.getPixelFormat()
            && lhs.getFps() == lhs.getFps();
    }

    inline bool operator!=(const openni::VideoMode& lhs, const openni::VideoMode& rhs)
    {
        return !(lhs == rhs);
    }

    template<typename TFrameWrapper>
    astra_status_t devicestream<TFrameWrapper>::on_read(astra_frame_index_t frameIndex)
    {
        PROFILE_FUNC();

        if (!is_streaming()) return ASTRA_STATUS_SUCCESS;

        openni::VideoFrameRef ref;
        PROFILE_BEGIN(oni_stream_readFrame);
        auto status = oniStream_.readFrame(&ref);
        PROFILE_END();

        if (status == ::openni::STATUS_OK)
        {
            const auto* oniFrameData = ref.getData();

            auto oniVideoMode = ref.getVideoMode();
            // LOG_INFO("orbbec.ni.devicestream", "type: %u oniM: %ux%u",
            //          description().type(),
            //          oniVideoMode.getResolutionX(),
            //          oniVideoMode.getResolutionY());

            if (bin_ == nullptr || oniVideoMode != oniMode_)
            {
                change_mode(convert_mode(oniVideoMode));
            }

            assert(ref.getDataSize() == bufferLength_);

            size_t byteSize = MIN(ref.getDataSize(), bufferLength_);

            wrapper_type* wrapper = bin_->begin_write(frameIndex);

            wrapper->frame.frame = nullptr;
            wrapper->frame.data = &(wrapper->frame_data);

            on_new_buffer(wrapper);

            std::memcpy(wrapper->frame.data, oniFrameData, byteSize);

            PROFILE_BEGIN(oni_stream_end_write);
            bin_->end_write();
            PROFILE_END();
        }

        return ASTRA_STATUS_SUCCESS;
    }
}}

#endif /* ONI_DEVICESTREAM_H */
