/* THIS FILE AUTO-GENERATED FROM astra_create_stream_proxy.hpp.lpp. DO NOT EDIT. */
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
#ifndef ASTRA_CREATE_STREAM_PROXY_H
#define ASTRA_CREATE_STREAM_PROXY_H

#include "astra_context.hpp"
#include "astra_stream_service_delegate.hpp"
#include <astra_core/capi/astra_streamservice_proxy.h>

namespace astra {

    static inline astra_streamservice_proxy_t* create_stream_proxy(context* context)
    {
        astra_streamservice_proxy_t* proxy = new astra_streamservice_proxy_t;

        proxy->streamset_open = &stream_service_delegate::streamset_open;
        proxy->streamset_close = &stream_service_delegate::streamset_close;
        proxy->reader_create = &stream_service_delegate::reader_create;
        proxy->reader_destroy = &stream_service_delegate::reader_destroy;
        proxy->reader_get_stream = &stream_service_delegate::reader_get_stream;
        proxy->stream_get_description = &stream_service_delegate::stream_get_description;
        proxy->stream_start = &stream_service_delegate::stream_start;
        proxy->stream_stop = &stream_service_delegate::stream_stop;
        proxy->reader_open_frame = &stream_service_delegate::reader_open_frame;
        proxy->reader_close_frame = &stream_service_delegate::reader_close_frame;
        proxy->reader_register_frame_ready_callback = &stream_service_delegate::reader_register_frame_ready_callback;
        proxy->reader_unregister_frame_ready_callback = &stream_service_delegate::reader_unregister_frame_ready_callback;
        proxy->reader_get_frame = &stream_service_delegate::reader_get_frame;
        proxy->stream_set_parameter = &stream_service_delegate::stream_set_parameter;
        proxy->stream_get_parameter = &stream_service_delegate::stream_get_parameter;
        proxy->stream_get_result = &stream_service_delegate::stream_get_result;
        proxy->stream_invoke = &stream_service_delegate::stream_invoke;
        proxy->temp_update = &stream_service_delegate::temp_update;
        proxy->streamService = context;

        return proxy;
    }
}

#endif /* ASTRA_CREATE_STREAM_PROXY_H */
