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
#ifndef ASTRA_STREAM_CONNECTION_H
#define ASTRA_STREAM_CONNECTION_H

#include <astra_core/capi/astra_types.h>
#include <astra_core/capi/plugins/astra_plugin.h>
#include "astra_parameter_bin.hpp"
#include "astra_stream_bin.hpp"
#include "astra_logger.hpp"
#include "astra_registry.hpp"

namespace astra {

    class stream;

    class stream_connection : public tracked_instance<stream_connection>
    {
    public:
        using FrameReadyCallback = std::function<void(stream_connection*, astra_frame_index_t)>;

        stream_connection(stream* stream);
        ~stream_connection();

        stream_connection& operator=(const stream_connection& rhs) = delete;
        stream_connection(const stream_connection& connection) = delete;

        void start();
        void stop();

        bool is_started() const { return started_; }

        astra_frame_t* lock();
        void unlock();

        void set_bin(stream_bin* bin);
        stream_bin* get_bin() const { return bin_; }

        stream* get_stream() const { return stream_; }

        astra_streamconnection_t get_handle() { return &connection_; }

        astra_bin_t get_bin_handle();

        static stream_connection* get_ptr(astra_streamconnection_t conn)
        {
            return registry::get<stream_connection>(conn->handle);
        }

        const astra_stream_desc_t& get_description() const;

        astra_callback_id_t register_frame_ready_callback(FrameReadyCallback callback);
        void unregister_frame_ready_callback(astra_callback_id_t& callbackId);

        void set_parameter(astra_parameter_id id,
                           size_t inByteLength,
                           astra_parameter_data_t inData);

        void get_parameter(astra_parameter_id id,
                           size_t& resultByteLength,
                           astra_result_token_t& token);

        astra_status_t get_result(astra_result_token_t token,
                                     size_t dataByteLength,
                                     astra_parameter_data_t dataDestination);

        void invoke(astra_command_id commandId,
                    size_t inByteLength,
                    astra_parameter_data_t inData,
                    size_t& resultByteLength,
                    astra_result_token_t& token);

    private:
        void on_bin_front_buffer_ready(stream_bin* bin, astra_frame_index_t frameIndex);
        void clear_pending_parameter_result();
        void cache_parameter_bin_token(astra_parameter_bin_t parameterBinHandle,
                                       size_t& resultByteLength,
                                       astra_result_token_t& token);

        _astra_streamconnection connection_;
        astra_frame_t* currentFrame_{nullptr};

        bool locked_{false};
        bool started_{false};

        stream* stream_{nullptr};
        stream_bin* bin_{nullptr};
        parameter_bin* pendingParameterResult_{nullptr};

        stream_bin::FrontBufferReadyCallback binFrontBufferReadyCallback_;
        astra_callback_id_t binFrontBufferReadyCallbackId_;

        signal<stream_connection*, astra_frame_index_t> frameReadySignal_;
    };
}

#endif /* ASTRA_STREAM_CONNECTION_H */
