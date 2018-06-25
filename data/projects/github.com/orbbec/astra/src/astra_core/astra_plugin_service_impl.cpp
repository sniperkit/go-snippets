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
#include "astra_plugin_service_impl.hpp"
#include "astra_streamset_catalog.hpp"
#include <astra_core/capi/astra_types.h>
#include "astra_stream_registered_event_args.hpp"
#include "astra_stream_unregistering_event_args.hpp"
#include "astra_parameter_bin.hpp"
#include "astra_logging.hpp"
#include <cstdio>
#include <memory>

namespace astra
{
    void plugin_service_impl::notify_host_event(astra_event_id id, const void* data, size_t dataSize)
    {
        hostEventSignal_.raise(id, data, dataSize);
    }

    astra_status_t plugin_service_impl::create_stream_set(const char* streamUri, astra_streamset_t& streamSet)
    {
        streamset& set = setCatalog_.get_or_add(streamUri, true);
        streamSet = set.get_handle();

        LOG_INFO("astra.plugin_service", "creating streamset: %s %x", streamUri, streamSet);

        return ASTRA_STATUS_SUCCESS;
    }

    astra_status_t plugin_service_impl::destroy_stream_set(astra_streamset_t& streamSet)
    {
        streamset* actualSet = streamset::get_ptr(streamSet);

        LOG_INFO("astra.plugin_service", "destroying streamset: %s %x", actualSet->get_uri().c_str(), streamSet);
        setCatalog_.destroy_set(actualSet);

        streamSet = nullptr;

        return ASTRA_STATUS_SUCCESS;
    }

    astra_status_t plugin_service_impl::register_stream_registered_callback(stream_registered_callback_t callback,
                                                                    void* clientTag,
                                                                    CallbackId& callbackId)
    {
        auto thunk = [clientTag, callback](stream_registered_event_args args)
            {
                callback(clientTag,
                         args.streamSet->get_handle(),
                         args.stream->get_handle(),
                         args.description);
            };

        callbackId = setCatalog_.register_for_stream_registered_event(thunk);

        return ASTRA_STATUS_SUCCESS;
    }

    astra_status_t plugin_service_impl::register_stream_unregistering_callback(stream_unregistering_callback_t callback,
                                                                       void* clientTag,
                                                                       CallbackId& callbackId)
    {
        auto thunk = [clientTag, callback](stream_unregistering_event_args args)
            {
                callback(clientTag,
                         args.streamSet->get_handle(),
                         args.stream->get_handle(),
                         args.description);
            };

        callbackId = setCatalog_.register_for_stream_unregistering_event(thunk);

        return ASTRA_STATUS_SUCCESS;
    }

    astra_status_t plugin_service_impl::unregister_stream_registered_callback(CallbackId callbackId)
    {
        setCatalog_.unregister_for_stream_registered_event(callbackId);

        return ASTRA_STATUS_SUCCESS;
    }

    astra_status_t plugin_service_impl::unregister_stream_unregistering_callback(CallbackId callbackId)
    {
        setCatalog_.unregister_form_stream_unregistering_event(callbackId);

        return ASTRA_STATUS_SUCCESS;
    }

    astra_status_t plugin_service_impl::create_stream(astra_streamset_t setHandle,
                                                   astra_stream_desc_t desc,
                                                   astra_stream_t& handle)
    {
        streamset* set = streamset::get_ptr(setHandle);
        stream* stream = set->register_stream(desc);
        handle = stream->get_handle();

        LOG_INFO("astra.plugin_service", "created stream -- handle %p type: %d", handle, desc.type);

        return ASTRA_STATUS_SUCCESS;
    }

    astra_status_t plugin_service_impl::destroy_stream(astra_stream_t& streamHandle)
    {
        if (streamHandle == nullptr)
            return ASTRA_STATUS_INVALID_PARAMETER;

        stream* stream = stream::get_ptr(streamHandle);

        assert(stream != nullptr);

        streamset* set =
            setCatalog_.find_streamset_for_stream(stream);

        assert(set != nullptr);

        const astra_stream_desc_t& desc = stream->get_description();

        LOG_INFO("astra.plugin_service", "destroying stream -- handle: %x type: %d", stream->get_handle(), desc.type);

        set->destroy_stream(stream);

        streamHandle = nullptr;

        return ASTRA_STATUS_SUCCESS;
    }

    astra_status_t plugin_service_impl::register_stream(astra_stream_t handle, stream_callbacks_t pluginCallbacks)
    {
        if (handle == nullptr)
            return ASTRA_STATUS_INVALID_PARAMETER;

        stream* stream = stream::get_ptr(handle);
        stream->set_callbacks(pluginCallbacks);

        return ASTRA_STATUS_SUCCESS;
    }

    astra_status_t plugin_service_impl::unregister_stream(astra_stream_t handle)
    {
        if (handle == nullptr)
            return ASTRA_STATUS_INVALID_PARAMETER;

        stream* stream = stream::get_ptr(handle);

        stream->clear_callbacks();

        return ASTRA_STATUS_SUCCESS;
    }

    astra_status_t plugin_service_impl::get_streamset_uri(astra_streamset_t setHandle,
                                                       const char*& uri)
    {
        assert(setHandle != nullptr);

        streamset* actualSet = streamset::get_ptr(setHandle);
        uri = actualSet->get_uri().c_str();

        return ASTRA_STATUS_SUCCESS;
    }

    astra_status_t plugin_service_impl::create_stream_bin(astra_stream_t streamHandle,
                                                       size_t lengthInBytes,
                                                       astra_bin_t& binHandle,
                                                       astra_frame_t*& binBuffer)
    {
        stream* actualStream = stream::get_ptr(streamHandle);
        stream_bin* bin = actualStream->create_bin(lengthInBytes);

        binHandle = bin->get_handle();
        binBuffer = bin->get_backBuffer();

        LOG_INFO("astra.plugin_service", "creating bin -- handle: %x stream: %x type: %d size: %u",
                      binHandle,
                      streamHandle,
                      actualStream->get_description().type,
                      lengthInBytes);

        return ASTRA_STATUS_SUCCESS;
    }

    astra_status_t plugin_service_impl::destroy_stream_bin(astra_stream_t streamHandle,
                                                        astra_bin_t& binHandle,
                                                        astra_frame_t*& binBuffer)
    {
        stream* actualStream = stream::get_ptr(streamHandle);
        stream_bin* bin = stream_bin::get_ptr(binHandle);

        LOG_INFO("astra.plugin_service", "destroying bin -- %x stream: %x type: %d size: %u",
                      binHandle,
                      streamHandle,
                      actualStream->get_description().type,
                      bin->bufferSize());

        actualStream->destroy_bin(bin);

        binHandle = nullptr;
        binBuffer = nullptr;

        return ASTRA_STATUS_SUCCESS;
    }

    astra_status_t plugin_service_impl::bin_has_connections(astra_bin_t binHandle, bool& hasConnections)
    {
        stream_bin* bin = stream_bin::get_ptr(binHandle);
        hasConnections = bin->has_clients_connected();

        return ASTRA_STATUS_SUCCESS;
    }

    astra_status_t plugin_service_impl::cycle_bin_buffers(astra_bin_t binHandle,
                                                       astra_frame_t*& binBuffer)
    {
        assert(binHandle != nullptr);

        stream_bin* bin = stream_bin::get_ptr(binHandle);
        binBuffer = bin->cycle_buffers();

        return ASTRA_STATUS_SUCCESS;
    }

    astra_status_t plugin_service_impl::link_connection_to_bin(astra_streamconnection_t connection,
                                                            astra_bin_t binHandle)
    {
        stream_connection* underlyingConnection = stream_connection::get_ptr(connection);
        stream_bin* bin = stream_bin::get_ptr(binHandle);

        stream* stream = underlyingConnection->get_stream();
        if (binHandle != nullptr)
        {
            LOG_INFO("astra.plugin_service", "linking connection to bin -- stream: %x type: %d conn: %x bin: %x",
                          stream->get_handle(),
                          stream->get_description().type,
                          connection,
                          bin);
        }
        else
        {
            LOG_INFO("astra.plugin_service", "unlinking connection to bin -- stream: %x type: %d conn: %x",
                          stream->get_handle(),
                          stream->get_description().type,
                          connection);
        }

        underlyingConnection->set_bin(bin);

        return ASTRA_STATUS_SUCCESS;
    }

    astra_status_t plugin_service_impl::get_parameter_bin(size_t byteSize,
                                                       astra_parameter_bin_t& binHandle,
                                                       astra_parameter_data_t& parameterData)
    {
        //TODO pooling
        parameter_bin* parameterBin = new parameter_bin(byteSize);

        binHandle = parameterBin->get_handle();
        parameterData = parameterBin->data();

        return ASTRA_STATUS_SUCCESS;
    }

    astra_status_t plugin_service_impl::log(const char* channel,
                              astra_log_severity_t logLevel,
                              const char* fileName,
                              int lineNo,
                              const char* func,
                              const char* format,
                              va_list args)
    {
        astra::log_vargs(channel, logLevel, fileName, lineNo, func, format, args);
        return ASTRA_STATUS_SUCCESS;
    }

    astra_status_t plugin_service_impl::register_host_event_callback(host_event_callback_t callback,
                                                                  void* clientTag,
                                                                  CallbackId& callbackId)
    {
        auto thunk = [clientTag, callback](astra_event_id id, const void* data, size_t dataSize)
            {
                callback(clientTag, id, data, dataSize);
            };

        callbackId = hostEventSignal_ += thunk;

        //context_.raise_existing_device_unavailable(callback, clientTag);

        return ASTRA_STATUS_SUCCESS;
    }

    astra_status_t plugin_service_impl::unregister_host_event_callback(CallbackId callbackId)
    {
        hostEventSignal_ -= callbackId;

        return ASTRA_STATUS_SUCCESS;
    }
}
