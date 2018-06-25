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
#ifndef PLAYBACKSTREAMSET_H
#define PLAYBACKSTREAMSET_H

#include <astra_core/capi/plugins/astra_plugin.h>
#include <memory>
#include <vector>
#include <string>
#include "PlaybackStream.h"
#include "DepthStream.h"
#include <common/serialization/FrameStreamReader.h>
#include <common/serialization/FrameInputStream.h>

namespace astra { namespace plugins { namespace streamplayer {

    const char STREAMPLAYERPLUGIN_FILE_PATH[] = "plugins/test.df";

    class PlaybackStreamSetBase
    {
    public:
        PlaybackStreamSetBase(PluginServiceProxy& pluginService, std::string uri)
            : m_pluginService(pluginService),
              m_uri(uri)
        {
            m_pluginService.create_stream_set(uri.c_str(), m_streamSetHandle);
        }
        virtual ~PlaybackStreamSetBase()
        {
            m_pluginService.destroy_stream_set(m_streamSetHandle);
        }

        virtual astra_status_t open() = 0;
        virtual astra_status_t close() = 0;
        virtual astra_status_t read() = 0;
    protected:
        PluginServiceProxy& m_pluginService;
        astra_streamset_t m_streamSetHandle;
        std::string m_uri;

        using StreamPtr = std::unique_ptr<PlaybackStreamBase>;
        using StreamPtrList = std::vector<StreamPtr>;

        StreamPtrList m_streams;
    };

    class PlaybackStreamSet : public PlaybackStreamSetBase
    {
    public:
        PlaybackStreamSet(PluginServiceProxy& pluginService, std::string uri)
            : PlaybackStreamSetBase(pluginService, uri)
        {}

        virtual ~PlaybackStreamSet()
        {
            close();
        }

        virtual astra_status_t open() override
        {
            if (m_isOpen)
            {
                return ASTRA_STATUS_SUCCESS;
            }

            try
            {
                m_frameStream = std::unique_ptr<FrameInputStream>(open_frame_input_stream(STREAMPLAYERPLUGIN_FILE_PATH));
                m_frameStreamReader = astra::make_unique<FrameStreamReader>(m_frameStream.get());

                open_stream();

                m_isOpen = true;
            }
            catch (ResourceNotFoundException)
            {
                return ASTRA_STATUS_DEVICE_ERROR;
            }

            return ASTRA_STATUS_SUCCESS;
        }

        virtual astra_status_t close() override final
        {
            if (!m_isOpen)
            {
                return ASTRA_STATUS_SUCCESS;
            }

            for (StreamPtr& s : m_streams)
            {
                s->close();
            }

            m_streams.clear();

            m_isOpen = false;

            return ASTRA_STATUS_SUCCESS;
        }

        virtual astra_status_t read() override
        {
            if (!m_isOpen || m_streams.size() == 0)
            {
                return ASTRA_STATUS_SUCCESS;
            }

            for (StreamPtr& s : m_streams)
            {
                s->read();
            }

            return ASTRA_STATUS_SUCCESS;
        }

    private:
        astra_status_t open_stream()
        {
            switch (m_frameStreamReader->get_stream_type())
            {
                case ASTRA_STREAM_DEPTH:
                {
                    DepthStream* stream = make_stream<DepthStream>(*m_frameStreamReader, m_pluginService, m_streamSetHandle);
                    stream->open();
                    m_streams.push_back(StreamPtr(stream));
                    break;
                }
                default:
                {
                    return ASTRA_STATUS_DEVICE_ERROR;
                }
            }

            return ASTRA_STATUS_SUCCESS;
        }


        bool m_isOpen{ false };

        std::unique_ptr<FrameStreamReader> m_frameStreamReader;
        std::unique_ptr<FrameInputStream> m_frameStream;
    };
}}}

#endif /* PLAYBACKSTREAMSET_H */