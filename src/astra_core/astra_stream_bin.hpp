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
#ifndef ASTRA_STREAM_BIN_H
#define ASTRA_STREAM_BIN_H

#include <exception>
#include <array>
#include <astra_core/capi/astra_types.h>
#include "astra_signal.hpp"
#include <astra_core/capi/plugins/astra_plugin.h>
#include "astra_logger.hpp"

namespace astra {

    class stream_bin
    {
    public:
        using FrontBufferReadyCallback = std::function<void(stream_bin*,astra_frame_index_t)>;

        stream_bin(size_t bufferSizeInBytes);
        ~stream_bin();

        stream_bin(const stream_bin& bin) = delete;
        stream_bin& operator=(const stream_bin& rhs) = delete;

        //exposed to plugins
        astra_frame_t* get_backBuffer()
            {
                return &buffers_[backBufferIndex_];
            }

        astra_frame_t* cycle_buffers();

        astra_frame_t* lock_front_buffer();
        void unlock_front_buffer();

        astra_callback_id_t register_front_buffer_ready_callback(FrontBufferReadyCallback callback);
        void unregister_front_buffer_ready_callback(astra_callback_id_t& callbackId);

        void inc_active() { activeCount_++; }
        void dec_active()
            {
                activeCount_--;
                if (activeCount_ < 0)
                {
                    throw std::exception();
                }
            }

        void inc_connected() { connectedCount_++; }
        void dec_connected()
            {
                connectedCount_--;
                if (connectedCount_ < 0)
                {
                    throw std::exception();
                }
            }

        bool is_active() { return activeCount_ > 0; }
        bool has_clients_connected() { return connectedCount_ > 0; }
        size_t bufferSize() { return bufferSize_; }

        astra_bin_t get_handle() { return reinterpret_cast<astra_bin_t>(this); }

        static stream_bin* get_ptr(astra_bin_t bin)
            { return reinterpret_cast<stream_bin*>(bin); }

    private:
        inline bool is_front_buffer_locked() { return frontBufferLockCount_ > 0; }
        void init_buffers(size_t bufferLengthInBytes);
        void deinit_buffers();
        void init_buffer(astra_frame_t& frame, size_t bufferLengthInBytes);
        void deinit_buffer(astra_frame_t& frame);
        astra_frame_t* get_frontBuffer();
        astra_frame_t* get_middleBuffer();
        void raiseFrameReadySignal();

        size_t bufferSize_{0};

        size_t frontBufferIndex_{0};
        size_t middleBufferIndex_{1};
        size_t backBufferIndex_{2};

        uint32_t frontBufferLockCount_{0};

        const static size_t BUFFER_COUNT = 3;
        astra_frame_t buffers_[BUFFER_COUNT];

        int connectedCount_{0};
        int activeCount_{0};

        signal<stream_bin*, astra_frame_index_t> frontBufferReadySignal_;
    };
}

#endif /* ASTRA_STREAM_BIN_H */
