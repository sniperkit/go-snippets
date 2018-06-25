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
#include "astra_context.hpp"
#include "astra_context_impl.hpp"
#include "astra_create_stream_proxy.hpp"
#include "astra_cxx_compatibility.hpp"
#include <astra_core/capi/astra_streamservice_proxy.h>
#include <astra_core_api.h>

namespace astra {

    context::context()
          : impl_(astra::make_unique<context_impl>()),
            proxy_(create_stream_proxy(this))
            {}

    context::~context() {}

    astra_status_t context::initialize()
    {
        astra_api_set_proxy(proxy());
        return impl_->initialize();
    }

    astra_status_t context::terminate()
    {
        return impl_->terminate();
    }

    astra_streamservice_proxy_t* context::proxy()
    {
        return proxy_.get();
    }

    astra_status_t context::streamset_open(const char* connectionString,
                                           astra_streamsetconnection_t& streamSet)
    {
        return impl_->streamset_open(connectionString, streamSet);
    }

    astra_status_t context::streamset_close(astra_streamsetconnection_t& streamSet)
    {
        return impl_->streamset_close(streamSet);
    }

    astra_status_t context::reader_create(astra_streamsetconnection_t streamSet,
                                          astra_reader_t& reader)
    {
        return impl_->reader_create(streamSet, reader);
    }

    astra_status_t context::reader_destroy(astra_reader_t& reader)
    {
        return impl_->reader_destroy(reader);
    }

    astra_status_t context::reader_get_stream(astra_reader_t reader,
                                              astra_stream_type_t type,
                                              astra_stream_subtype_t subtype,
                                              astra_streamconnection_t& connection)
    {
        return impl_->reader_get_stream(reader, type, subtype, connection);
    }

    astra_status_t context::stream_get_description(astra_streamconnection_t connection,
                                                   astra_stream_desc_t* description)
    {
        return impl_->stream_get_description(connection, description);
    }

    astra_status_t context::stream_start(astra_streamconnection_t connection)
    {
        return impl_->stream_start(connection);
    }

    astra_status_t context::stream_stop(astra_streamconnection_t connection)
    {
        return impl_->stream_stop(connection);
    }

    astra_status_t context::reader_open_frame(astra_reader_t reader,
                                              int timeoutMillis,
                                              astra_reader_frame_t& frame)
    {
        return impl_->reader_open_frame(reader, timeoutMillis, frame);
    }

    astra_status_t context::reader_close_frame(astra_reader_frame_t& frame)
    {
        return impl_->reader_close_frame(frame);
    }

    astra_status_t context::reader_register_frame_ready_callback(astra_reader_t reader,
                                                                 astra_frame_ready_callback_t callback,
                                                                 void* clientTag,
                                                                 astra_reader_callback_id_t& callbackId)
    {
        return impl_->reader_register_frame_ready_callback(reader, callback, clientTag, callbackId);
    }

    astra_status_t context::reader_unregister_frame_ready_callback(astra_reader_callback_id_t& callbackId)
    {
        return impl_->reader_unregister_frame_ready_callback(callbackId);
    }

    astra_status_t context::reader_get_frame(astra_reader_frame_t frame,
                                             astra_stream_type_t type,
                                             astra_stream_subtype_t subtype,
                                             astra_frame_t*& subFrame)
    {
        return impl_->reader_get_frame(frame, type, subtype, subFrame);
    }

    astra_status_t context::stream_set_parameter(astra_streamconnection_t connection,
                                                 astra_parameter_id parameterId,
                                                 size_t inByteLength,
                                                 astra_parameter_data_t inData)
    {
        return impl_->stream_set_parameter(connection, parameterId, inByteLength, inData);
    }

    astra_status_t context::stream_get_parameter(astra_streamconnection_t connection,
                                                 astra_parameter_id parameterId,
                                                 size_t& resultByteLength,
                                                 astra_result_token_t& token)
    {
        return impl_->stream_get_parameter(connection, parameterId, resultByteLength, token);
    }

    astra_status_t context::stream_get_result(astra_streamconnection_t connection,
                                              astra_result_token_t token,
                                              size_t dataByteLength,
                                              astra_parameter_data_t dataDestination)
    {
        return impl_->stream_get_result(connection, token, dataByteLength, dataDestination);
    }

    astra_status_t context::stream_invoke(astra_streamconnection_t connection,
                                          astra_command_id commandId,
                                          size_t inByteLength,
                                          astra_parameter_data_t inData,
                                          size_t& resultByteLength,
                                          astra_result_token_t& token)
    {
        return impl_->stream_invoke(connection, commandId, inByteLength, inData, resultByteLength, token);
    }

    astra_status_t context::temp_update()
    {
        return impl_->temp_update();
    }


    astra_status_t context::notify_host_event(astra_event_id id, const void* data, size_t dataSize)
    {
        return impl_->notify_host_event(id, data, dataSize);
    }
}
