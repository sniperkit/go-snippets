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
#ifndef ASTRA_STREAMSET_CONNECTION_H
#define ASTRA_STREAMSET_CONNECTION_H

#include "astra_registry.hpp"
#include <astra_core/capi/astra_types.h>
#include "astra_stream_reader.hpp"
#include <vector>
#include <memory>
#include "astra_logger.hpp"

namespace astra {

    class streamset;

    class streamset_connection : public tracked_instance<streamset_connection>
    {
    public:
        streamset_connection(streamset* streamSet)
            : streamSet_(streamSet)
        {}

        ~streamset_connection()
        {
            LOG_TRACE("streamset_connection", "destroying streamset_connection: %p", this);
        }

        streamset_connection& operator=(const streamset_connection& rhs) = delete;
        streamset_connection(const streamset_connection& conn) = delete;

        streamset* get_streamSet() { return streamSet_; }

        stream_reader* create_reader();
        bool destroy_reader(stream_reader* reader);

        bool is_connected() { return streamSet_ != nullptr; }

        static streamset_connection* get_ptr(astra_streamsetconnection_t conn)
        {
            return registry::get<streamset_connection>(conn);
        }

        astra_streamsetconnection_t get_handle()
        {
            return reinterpret_cast<astra_streamsetconnection_t>(this);
        }

    private:
        streamset* streamSet_{nullptr};

        using ReaderPtr = std::unique_ptr<stream_reader>;
        using ReaderList = std::vector<ReaderPtr>;

        ReaderList readers_;
    };
}


#endif /* ASTRA_STREAMSET_CONNECTION_H */
