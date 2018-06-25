/* THIS FILE AUTO-GENERATED FROM astra_core_api.cpp.lpp. DO NOT EDIT. */
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
#include "astra_core_api.h"
#include <astra_core/StreamServiceProxy.hpp>

ASTRA_BEGIN_DECLS

ASTRA_LOCAL astra_streamservice_proxy_t* __astra_api_proxy_ptr = nullptr;

static astra::StreamServiceProxy* get_api_proxy()
{
    return reinterpret_cast<astra::StreamServiceProxy*>(__astra_api_proxy_ptr);
}

ASTRA_API_PROXY astra_streamservice_proxy_t* astra_api_get_proxy()
{
    return __astra_api_proxy_ptr;
}

ASTRA_API_PROXY void astra_api_set_proxy(astra_streamservice_proxy_t* proxy)
{
    __astra_api_proxy_ptr = proxy;
}

ASTRA_API astra_status_t astra_streamset_open(const char* connectionString,
                                              astra_streamsetconnection_t* streamSet)
{
    return get_api_proxy()->streamset_open(connectionString, streamSet);
}

ASTRA_API astra_status_t astra_streamset_close(astra_streamsetconnection_t* streamSet)
{
    return get_api_proxy()->streamset_close(streamSet);
}

ASTRA_API astra_status_t astra_reader_create(astra_streamsetconnection_t streamSet,
                                             astra_reader_t* reader)
{
    return get_api_proxy()->reader_create(streamSet, reader);
}

ASTRA_API astra_status_t astra_reader_destroy(astra_reader_t* reader)
{
    return get_api_proxy()->reader_destroy(reader);
}

ASTRA_API astra_status_t astra_reader_get_stream(astra_reader_t reader,
                                                 astra_stream_type_t type,
                                                 astra_stream_subtype_t subtype,
                                                 astra_streamconnection_t* connection)
{
    return get_api_proxy()->reader_get_stream(reader, type, subtype, connection);
}

ASTRA_API astra_status_t astra_stream_get_description(astra_streamconnection_t connection,
                                                      astra_stream_desc_t* description)
{
    return get_api_proxy()->stream_get_description(connection, description);
}

ASTRA_API astra_status_t astra_stream_start(astra_streamconnection_t connection)
{
    return get_api_proxy()->stream_start(connection);
}

ASTRA_API astra_status_t astra_stream_stop(astra_streamconnection_t connection)
{
    return get_api_proxy()->stream_stop(connection);
}

ASTRA_API astra_status_t astra_reader_open_frame(astra_reader_t reader,
                                                 int timeoutMillis,
                                                 astra_reader_frame_t* frame)
{
    return get_api_proxy()->reader_open_frame(reader, timeoutMillis, frame);
}

ASTRA_API astra_status_t astra_reader_close_frame(astra_reader_frame_t* frame)
{
    return get_api_proxy()->reader_close_frame(frame);
}

ASTRA_API astra_status_t astra_reader_register_frame_ready_callback(astra_reader_t reader,
                                                                    astra_frame_ready_callback_t callback,
                                                                    void* clientTag,
                                                                    astra_reader_callback_id_t* callbackId)
{
    return get_api_proxy()->reader_register_frame_ready_callback(reader, callback, clientTag, callbackId);
}

ASTRA_API astra_status_t astra_reader_unregister_frame_ready_callback(astra_reader_callback_id_t* callbackId)
{
    return get_api_proxy()->reader_unregister_frame_ready_callback(callbackId);
}

ASTRA_API astra_status_t astra_reader_get_frame(astra_reader_frame_t frame,
                                                astra_stream_type_t type,
                                                astra_stream_subtype_t subtype,
                                                astra_frame_t** subFrame)
{
    return get_api_proxy()->reader_get_frame(frame, type, subtype, subFrame);
}

ASTRA_API astra_status_t astra_stream_set_parameter(astra_streamconnection_t connection,
                                                    astra_parameter_id parameterId,
                                                    size_t inByteLength,
                                                    astra_parameter_data_t inData)
{
    return get_api_proxy()->stream_set_parameter(connection, parameterId, inByteLength, inData);
}

ASTRA_API astra_status_t astra_stream_get_parameter(astra_streamconnection_t connection,
                                                    astra_parameter_id parameterId,
                                                    size_t* resultByteLength,
                                                    astra_result_token_t* token)
{
    return get_api_proxy()->stream_get_parameter(connection, parameterId, resultByteLength, token);
}

ASTRA_API astra_status_t astra_stream_get_result(astra_streamconnection_t connection,
                                                 astra_result_token_t token,
                                                 size_t dataByteLength,
                                                 astra_parameter_data_t dataDestination)
{
    return get_api_proxy()->stream_get_result(connection, token, dataByteLength, dataDestination);
}

ASTRA_API astra_status_t astra_stream_invoke(astra_streamconnection_t connection,
                                             astra_command_id commandId,
                                             size_t inByteLength,
                                             astra_parameter_data_t inData,
                                             size_t* resultByteLength,
                                             astra_result_token_t* token)
{
    return get_api_proxy()->stream_invoke(connection, commandId, inByteLength, inData, resultByteLength, token);
}

ASTRA_API astra_status_t astra_temp_update()
{
    return get_api_proxy()->temp_update();
}

ASTRA_END_DECLS
