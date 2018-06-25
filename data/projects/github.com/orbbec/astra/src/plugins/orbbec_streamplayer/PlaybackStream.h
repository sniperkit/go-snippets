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
#ifndef PLAYBACKSTREAM_H
#define PLAYBACKSTREAM_H

#include <astra_core/plugins/PluginStream.hpp>
#include <astra_core/plugins/StreamBin.hpp>
#include <astra_core/capi/plugins/astra_plugin.h>
#include <astra/capi/streams/image_parameters.h>
#include <common/serialization/FrameStreamReader.h>
#include <cstring>

using namespace astra::serialization;

namespace astra { namespace plugins { namespace streamplayer {

    class PlaybackStreamBase : public stream
    {
    public:
        PlaybackStreamBase(PluginServiceProxy& pluginService,
                           astra_streamset_t streamSet,
                           StreamDescription desc)
            : stream(pluginService,
                     streamSet,
                     desc) { }

        virtual ~PlaybackStreamBase() { };

        virtual astra_status_t read() = 0;
        virtual astra_status_t open() = 0;
        virtual astra_status_t close() = 0;
        virtual astra_status_t start() = 0;
        virtual astra_status_t stop() = 0;
    };

    template<typename TFrameWrapper>
    class PlaybackStream : public PlaybackStreamBase
    {
    public:
        using wrapper_type = TFrameWrapper;

        PlaybackStream(FrameStreamReader& depthStreamParser,
                       PluginServiceProxy& pluginService,
                       astra_streamset_t streamSet,
                       StreamDescription desc);

        virtual ~PlaybackStream();

        virtual astra_status_t close() override final;
        virtual astra_status_t start() override { return ASTRA_STATUS_SUCCESS; }
        virtual astra_status_t stop() override { return ASTRA_STATUS_SUCCESS; }
        virtual astra_status_t open() override;
        virtual astra_status_t read() override;

        bool is_streaming();

        void on_connection_added(astra_streamconnection_t connection) override;
        void on_connection_removed(astra_bin_t bin,
                                    astra_streamconnection_t connection) override;

        virtual void on_get_parameter(astra_streamconnection_t connection,
                                        astra_parameter_id id,
                                        astra_parameter_bin_t& parameterBin) override;

    protected:
        virtual astra_status_t on_open()
        {
            return ASTRA_STATUS_SUCCESS;
        }

        FrameStreamReader& m_frameStreamReader;

    private:
        bool m_isOpen{ false };
        bool m_isStreaming{ true };

        int m_frameIndex{ 0 };

        using BinType = stream_bin<wrapper_type>;
        std::unique_ptr<BinType> m_bin;
    };

    template<typename TFrameWrapper>
    PlaybackStream<TFrameWrapper>::PlaybackStream(FrameStreamReader& frameStreamReader,
                                                  PluginServiceProxy& pluginService,
                                                  astra_streamset_t streamSet,
                                                  StreamDescription desc) :
        PlaybackStreamBase(pluginService,
                           streamSet,
                           desc),
        m_frameStreamReader(frameStreamReader)
    {

    }

    template<typename TFrameWrapper>
    PlaybackStream<TFrameWrapper>::~PlaybackStream()
    {
        close();
    }

    template<typename TFrameWrapper>
    astra_status_t PlaybackStream<TFrameWrapper>::open()
    {
        if (m_isOpen)
        {
            return ASTRA_STATUS_SUCCESS;
        }

        m_bin = astra::make_unique<stream_bin<wrapper_type> >(
            pluginService(),
            get_handle(),
            m_frameStreamReader.get_buffer_length() - sizeof(wrapper_type));

        on_open();

        m_isOpen = true;

        return ASTRA_STATUS_SUCCESS;
    }

    template<typename TFrameWrapper>
    astra_status_t PlaybackStream<TFrameWrapper>::close()
    {
        if (!m_isOpen)
            return ASTRA_STATUS_SUCCESS;

        m_isOpen = m_isStreaming = false;

        return ASTRA_STATUS_SUCCESS;
    }

    template<typename TFrameWrapper>
    bool PlaybackStream<TFrameWrapper>::is_streaming()
    {
        return m_isOpen && m_isStreaming;
    }

    template<typename TFrameWrapper>
    void PlaybackStream<TFrameWrapper>::on_connection_added(
        astra_streamconnection_t connection)
    {
        m_bin->link_connection(connection);
    }

    template<typename TFrameWrapper>
    void PlaybackStream<TFrameWrapper>::on_connection_removed(
        astra_bin_t bin, astra_streamconnection_t connection)
    {
        m_bin->unlink_connection(connection);

        if (!m_bin->has_connections())
        {
            m_bin = nullptr;
        }
    }

    template<typename TFrameWrapper>
    void PlaybackStream<TFrameWrapper>::on_get_parameter(astra_streamconnection_t connection,
                                                            astra_parameter_id id,
                                                            astra_parameter_bin_t& parameterBin)
    {
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
                    *hFov = 1.02259994;
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
                    *vFov = 0.796615660;
                }
                break;
            }
        }
    }

    template<typename TFrameWrapper>
    astra_status_t PlaybackStream<TFrameWrapper>::read()
    {
        if (!is_streaming()) return ASTRA_STATUS_SUCCESS;

        if (!m_frameStreamReader.read())
        {
            return ASTRA_STATUS_SUCCESS;
        }

        astra::serialization::Frame& decodedFrame = m_frameStreamReader.peek();

        auto framePair = m_bin->begin_write_ex(m_frameIndex);

        wrapper_type* frameWrapper = framePair.second;
        astra_frame_t* frame = framePair.first;

	std::memcpy(frame->data, decodedFrame.rawFrameWrapper, decodedFrame.byteLength);
        frameWrapper->frame.data = &(frameWrapper->frame_data);
        frameWrapper->frame.frame = frame;

        m_bin->end_write();

        ++m_frameIndex;

        return ASTRA_STATUS_SUCCESS;
    }

}}}

#endif /* PLAYBACKSTREAM_H */
