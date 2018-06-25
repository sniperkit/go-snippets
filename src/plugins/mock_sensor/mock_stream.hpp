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
#ifndef MOCK_STREAM_H
#define MOCK_STREAM_H

#include <memory>

#include <astra_core/astra_core.hpp>
#include <astra_core/capi/plugins/astra_plugin.h>
#include <astra_core/plugins/PluginStream.hpp>
#include <astra_core/plugins/StreamBin.hpp>

#include "mock_stream_listener.hpp"

namespace orbbec { namespace mocks {

    class mock_stream : public astra::plugins::stream
    {
    public:
        inline mock_stream(astra::PluginServiceProxy& pluginService,
                           astra_streamset_t streamSet,
                           astra::StreamDescription desc,
                           stream_listener& listener);

        inline astra_status_t read(astra_frame_index_t frameIndex);
        inline astra_status_t open();
        inline astra_status_t close();
        inline astra_status_t start();
        inline astra_status_t stop();

        bool is_open() const { return isOpen_; }
        bool is_started() const { return isStarted_; }

    protected:
        virtual astra_status_t on_read(astra_frame_index_t frameIndex)
        {
            return ASTRA_STATUS_SUCCESS;
        }

        virtual astra_status_t on_open() { return ASTRA_STATUS_SUCCESS; }
        virtual astra_status_t on_close() { return ASTRA_STATUS_SUCCESS; }
        virtual astra_status_t on_start() { return ASTRA_STATUS_SUCCESS; }
        virtual astra_status_t on_stop() { return ASTRA_STATUS_SUCCESS; }

    private:
        bool isOpen_{false};
        bool isStarted_{false};
        stream_listener& listener_;
    };
}}

namespace orbbec { namespace mocks {

    mock_stream::mock_stream(astra::PluginServiceProxy& pluginService,
                             astra_streamset_t streamSet,
                             astra::StreamDescription desc,
                             orbbec::mocks::stream_listener& listener)
        : astra::plugins::stream(pluginService,
                                 streamSet,
                                 desc),
          listener_(listener)
    {}

    astra_status_t mock_stream::read(astra_frame_index_t frameIndex)
    {
        if (!is_open() || !is_started())
            return astra_status_t::ASTRA_STATUS_INVALID_OPERATION;

        return on_read(frameIndex);
    }

    astra_status_t mock_stream::open()
    {
        if (is_open())
            return astra_status_t::ASTRA_STATUS_SUCCESS;

        auto rc = on_open();

        if (rc == astra_status_t::ASTRA_STATUS_SUCCESS)
        {
            isOpen_ = true;
            listener_.on_opened(this);
        }

        return rc;
    }

    astra_status_t mock_stream::close()
    {
        if (!is_open())
            return astra_status_t::ASTRA_STATUS_SUCCESS;

        if (is_started())
        {
            auto rc = stop();
            if (rc != astra_status_t::ASTRA_STATUS_SUCCESS)
                return rc;
        }

        auto rc = on_close();

        if (rc == astra_status_t::ASTRA_STATUS_SUCCESS)
        {
            isOpen_ = false;
            listener_.on_closed(this);
        }

        return rc;
    }

    astra_status_t mock_stream::start()
    {
        if (isStarted_)
            return astra_status_t::ASTRA_STATUS_SUCCESS;

        auto rc = on_start();

        if (rc == astra_status_t::ASTRA_STATUS_SUCCESS)
        {
            isStarted_ = true;
            listener_.on_started(this);
        }

        return rc;
    }

    astra_status_t mock_stream::stop()
    {
        if (!isStarted_)
            return astra_status_t::ASTRA_STATUS_SUCCESS;

        auto rc = on_stop();

        if (rc == astra_status_t::ASTRA_STATUS_SUCCESS)
        {
            isStarted_ = false;
            listener_.on_stopped(this);
        }

        return rc;
    }
}}

#endif /* MOCK_STREAM_H */
