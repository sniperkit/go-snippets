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
#ifndef ASTRA_STREAM_H
#define ASTRA_STREAM_H

#include <atomic>
#include <memory>
#include <vector>

#include "astra_signal.hpp"
#include "astra_stream_connection.hpp"
#include "astra_stream_backend.hpp"

namespace astra {

    class stream_listener;

    class stream : public stream_backend
    {
    public:
        stream(astra_stream_desc_t description);
        virtual ~stream();

        stream& operator=(const stream& stream) = delete;
        stream(const stream& stream) = delete;

        stream_connection* create_connection();
        void destroy_connection(stream_connection* connection);

        void start_connection(stream_connection* connection);
        void stop_connection(stream_connection* connection);

        void set_parameter(stream_connection* connection,
                           astra_parameter_id id,
                           size_t inByteLength,
                           astra_parameter_data_t inData);

        void get_parameter(stream_connection* connection,
                           astra_parameter_id id,
                           astra_parameter_bin_t& parameterBin);

        void invoke(stream_connection* connection,
                    astra_command_id commandId,
                    size_t inByteLength,
                    astra_parameter_data_t inData,
                    astra_parameter_bin_t& parameterBin);

        astra_stream_t get_handle()
        {
            return reinterpret_cast<astra_stream_t>(this);
        }

        static stream* get_ptr(astra_stream_t stream)
        {
            return reinterpret_cast<astra::stream*>(stream);
        }

        virtual void on_availability_changed() override;

        virtual void on_destroying_bin(stream_bin* bin) override
        {
            disconnect_connections(bin);
        }

        bool has_connections() { return connections_.size() > 0; }

        void set_listener(stream_listener* listener)
        {
            listener_ = listener;
        }

    private:
        void disconnect_connections(stream_bin* bin);

        using connection_ptr = std::unique_ptr<stream_connection>;
        using connection_vector = std::vector<connection_ptr>;

        connection_vector connections_;
        stream_listener* listener_{nullptr};
    };
}

#endif /* ASTRA_STREAM_H */
