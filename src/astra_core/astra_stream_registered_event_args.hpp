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
#ifndef ASTRA_STREAM_REGISTERED_EVENT_ARGS_H
#define ASTRA_STREAM_REGISTERED_EVENT_ARGS_H

#include <astra_core/capi/astra_types.h>

namespace astra {

    class streamset;
    class stream;

    struct stream_registered_event_args
    {
        streamset* streamSet;
        astra::stream* stream;
        astra_stream_desc_t description;

        stream_registered_event_args(streamset* streamSet, astra::stream* strm, astra_stream_desc_t description)
            : streamSet(streamSet), stream(strm), description(description)
        { }
    };

    using StreamRegisteredCallback = std::function<void(stream_registered_event_args)>;
}

#endif /* ASTRA_STREAM_REGISTERED_EVENT_ARGS_H */
