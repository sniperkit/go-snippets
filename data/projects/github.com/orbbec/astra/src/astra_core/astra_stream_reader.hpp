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
#ifndef ASTRA_STREAM_READER_H
#define ASTRA_STREAM_READER_H

#include <astra_core/capi/astra_types.h>
#include "astra_registry.hpp"
#include <unordered_map>
#include <vector>
#include <cassert>
#include "astra_signal.hpp"
#include "astra_private.h"
#include "astra_stream_connection.hpp"

namespace astra {

    class streamset_connection;
    //class stream_connection;

    class stream_desc_hash
    {
    public:
        std::size_t operator()(const astra_stream_desc_t desc) const
            {
                std::size_t h1 = std::hash<astra_stream_type_t>()(desc.type);
                std::size_t h2 = std::hash<astra_stream_subtype_t>()(desc.subtype);

                return h1 ^ (h2 << 1);
            }
    };

    class stream_desc_equal_to
    {
    public:
        std::size_t operator()(const astra_stream_desc_t& lhs,
                               const astra_stream_desc_t& rhs) const
            {
                return lhs.type == rhs.type && lhs.subtype == rhs.subtype;
            }
    };

    struct reader_connection_data
    {
        stream_connection* connection;
        astra_callback_id_t scFrameReadyCallbackId;
        bool isNewFrameReady;
        astra_frame_index_t currentFrameIndex;
    };

    class stream_reader : public tracked_instance<stream_reader>
    {
    public:
        stream_reader(streamset_connection& connection);
        ~stream_reader();

        stream_reader& operator=(const stream_reader& rhs) = delete;
        stream_reader(const stream_reader& reader) = delete;

        inline streamset_connection& get_connection() const { return connection_; }

        inline astra_reader_t get_handle() { return reinterpret_cast<astra_reader_t>(this); }

        stream_connection* get_stream(astra_stream_desc_t& desc);
        astra_frame_t* get_subframe(astra_stream_desc_t& desc);

        astra_callback_id_t register_frame_ready_callback(astra_frame_ready_callback_t callback, void* clientTag);
        void unregister_frame_ready_callback(astra_callback_id_t& callbackId);

        //TODO: locking currently not threadsafe

        astra_status_t lock(int timeoutMillis, astra_reader_frame_t& readerFrame);
        astra_status_t unlock(astra_reader_frame_t& readerFrame);

        static inline stream_reader* get_ptr(astra_reader_t reader) { return registry::get<stream_reader>(reader); }
        static inline stream_reader* from_frame(astra_reader_frame_t& frame)
        {
            if (frame == nullptr)
            {
                return nullptr;
            }

            return get_ptr(frame->reader);
        }

    private:
        enum class block_result
        {
            TIMEOUT,
            FRAMEREADY
        };

        block_result block_until_frame_ready_or_timeout(int timeoutMillis);

        astra_reader_frame_t lock_frame_for_event_callback();
        astra_reader_frame_t lock_frame_for_poll();
        astra_reader_frame_t acquire_available_reader_frame();

        astra_status_t unlock_frame_and_check_connections(astra_reader_frame_t& readerFrame);
        astra_status_t return_locked_frame(astra_reader_frame_t& readerFrame);

        void ensure_connections_locked();
        astra_status_t unlock_connections_if_able();

        stream_connection* find_stream_of_type(astra_stream_desc_t& desc);
        stream_connection::FrameReadyCallback get_sc_frame_ready_callback();
        void on_connection_frame_ready(stream_connection* connection, astra_frame_index_t frameIndex);
        void check_for_all_frames_ready();
        void raise_frame_ready();

        bool locked_{false};
        bool isFrameReadyForLock_{false};
        astra_frame_index_t lastFrameIndex_{-1};
        streamset_connection& connection_;

        using ConnectionMap = std::unordered_map<astra_stream_desc_t,
                                                 reader_connection_data*,
                                                 stream_desc_hash,
                                                 stream_desc_equal_to>;
        ConnectionMap streamMap_;

        using FramePtr  = std::unique_ptr<_astra_reader_frame>;
        using FrameList = std::vector<FramePtr>;

        FrameList frameList_;
        int32_t lockedFrameCount_{0};

        signal<astra_reader_t, astra_reader_frame_t> frameReadySignal_;

        stream_connection::FrameReadyCallback scFrameReadyCallback_;
    };
}

#endif /* ASTRA_STREAM_READER_H */
