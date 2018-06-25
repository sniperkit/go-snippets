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
#ifndef ASTRA_STREAM_BACKEND_H
#define ASTRA_STREAM_BACKEND_H

#include <astra_core/capi/astra_types.h>
#include <astra_core/capi/plugins/astra_plugin.h>
#include "astra_stream_connection.hpp"
#include <vector>
#include <memory>

namespace astra {

    class stream_bin;

    class stream_backend
    {
    public:
        stream_backend(astra_stream_desc_t description)
            : description_(description) {}

        virtual ~stream_backend()
        {
            bins_.clear();
        }

        stream_bin* create_bin(size_t byteLength);
        void destroy_bin(stream_bin* bin);

        const astra_stream_desc_t& get_description() const { return description_; }

        bool is_available() { return callbacks_ != nullptr; }

        void set_callbacks(const stream_callbacks_t& callbacks);
        void clear_callbacks();

    protected:
        void on_connection_created(stream_connection* connection, astra_stream_t stream);
        void on_connection_destroyed(stream_connection* connection, astra_stream_t stream);

        void on_connection_started(stream_connection* connection, astra_stream_t stream);
        void on_connection_stopped(stream_connection* connection, astra_stream_t stream);

        void on_set_parameter(stream_connection* connection,
                              astra_parameter_id id,
                              size_t inByteLength,
                              astra_parameter_data_t inData);

        void on_get_parameter(stream_connection* connection,
                              astra_parameter_id id,
                              astra_parameter_bin_t& parameterBin);

        void on_invoke(stream_connection* connection,
                       astra_command_id commandId,
                       size_t inByteLength,
                       astra_parameter_data_t inData,
                       astra_parameter_bin_t& parameterBin);

        virtual void on_destroying_bin(stream_bin* bin) {};

        virtual void on_availability_changed() = 0;

    private:
        using bin_ptr = std::unique_ptr<stream_bin>;
        using bin_list = std::vector<bin_ptr>;

        astra_stream_desc_t description_;
        std::unique_ptr<stream_callbacks_t> callbacks_{nullptr};

        bin_list bins_;
    };
}

#endif /* ASTRA_STREAM_BACKEND_H */
