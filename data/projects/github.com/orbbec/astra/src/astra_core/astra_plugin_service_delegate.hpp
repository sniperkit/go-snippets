/* THIS FILE AUTO-GENERATED FROM astra_plugin_service_delegate.hpp.lpp. DO NOT EDIT. */
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
#ifndef ASTRA_PLUGIN_SERVICE_DELEGATE_H
#define ASTRA_PLUGIN_SERVICE_DELEGATE_H

#include <astra_core/capi/astra_types.h>
#include <stdarg.h>
#include "astra_context.hpp"

namespace astra {

    class plugin_service_delegate
    {
    public:

        static astra_status_t register_stream_registered_callback(void* pluginService,
                                                                  stream_registered_callback_t callback,
                                                                  void* clientTag,
                                                                  astra_callback_id_t* callbackId)
        {
            return static_cast<plugin_service*>(pluginService)->register_stream_registered_callback(callback, clientTag, *callbackId);
        }

        static astra_status_t register_stream_unregistering_callback(void* pluginService,
                                                                     stream_unregistering_callback_t callback,
                                                                     void* clientTag,
                                                                     astra_callback_id_t* callbackId)
        {
            return static_cast<plugin_service*>(pluginService)->register_stream_unregistering_callback(callback, clientTag, *callbackId);
        }

        static astra_status_t register_host_event_callback(void* pluginService,
                                                           host_event_callback_t callback,
                                                           void* clientTag,
                                                           astra_callback_id_t* callbackId)
        {
            return static_cast<plugin_service*>(pluginService)->register_host_event_callback(callback, clientTag, *callbackId);
        }

        static astra_status_t unregister_host_event_callback(void* pluginService,
                                                             astra_callback_id_t callback)
        {
            return static_cast<plugin_service*>(pluginService)->unregister_host_event_callback(callback);
        }

        static astra_status_t unregister_stream_registered_callback(void* pluginService,
                                                                    astra_callback_id_t callback)
        {
            return static_cast<plugin_service*>(pluginService)->unregister_stream_registered_callback(callback);
        }

        static astra_status_t unregister_stream_unregistering_callback(void* pluginService,
                                                                       astra_callback_id_t callback)
        {
            return static_cast<plugin_service*>(pluginService)->unregister_stream_unregistering_callback(callback);
        }

        static astra_status_t create_stream_set(void* pluginService,
                                                const char* setUri,
                                                astra_streamset_t& setHandle)
        {
            return static_cast<plugin_service*>(pluginService)->create_stream_set(setUri, setHandle);
        }

        static astra_status_t destroy_stream_set(void* pluginService,
                                                 astra_streamset_t& setHandle)
        {
            return static_cast<plugin_service*>(pluginService)->destroy_stream_set(setHandle);
        }

        static astra_status_t get_streamset_uri(void* pluginService,
                                                astra_streamset_t setHandle,
                                                const char** uri)
        {
            return static_cast<plugin_service*>(pluginService)->get_streamset_uri(setHandle, *uri);
        }

        static astra_status_t create_stream(void* pluginService,
                                            astra_streamset_t setHandle,
                                            astra_stream_desc_t desc,
                                            astra_stream_t* handle)
        {
            return static_cast<plugin_service*>(pluginService)->create_stream(setHandle, desc, *handle);
        }

        static astra_status_t register_stream(void* pluginService,
                                              astra_stream_t handle,
                                              stream_callbacks_t pluginCallbacks)
        {
            return static_cast<plugin_service*>(pluginService)->register_stream(handle, pluginCallbacks);
        }

        static astra_status_t unregister_stream(void* pluginService,
                                                astra_stream_t handle)
        {
            return static_cast<plugin_service*>(pluginService)->unregister_stream(handle);
        }

        static astra_status_t destroy_stream(void* pluginService,
                                             astra_stream_t& handle)
        {
            return static_cast<plugin_service*>(pluginService)->destroy_stream(handle);
        }

        static astra_status_t create_stream_bin(void* pluginService,
                                                astra_stream_t streamHandle,
                                                size_t lengthInBytes,
                                                astra_bin_t* binHandle,
                                                astra_frame_t** binBuffer)
        {
            return static_cast<plugin_service*>(pluginService)->create_stream_bin(streamHandle, lengthInBytes, *binHandle, *binBuffer);
        }

        static astra_status_t destroy_stream_bin(void* pluginService,
                                                 astra_stream_t streamHandle,
                                                 astra_bin_t* binHandle,
                                                 astra_frame_t** binBuffer)
        {
            return static_cast<plugin_service*>(pluginService)->destroy_stream_bin(streamHandle, *binHandle, *binBuffer);
        }

        static astra_status_t bin_has_connections(void* pluginService,
                                                  astra_bin_t binHandle,
                                                  bool* hasConnections)
        {
            return static_cast<plugin_service*>(pluginService)->bin_has_connections(binHandle, *hasConnections);
        }

        static astra_status_t cycle_bin_buffers(void* pluginService,
                                                astra_bin_t binHandle,
                                                astra_frame_t** binBuffer)
        {
            return static_cast<plugin_service*>(pluginService)->cycle_bin_buffers(binHandle, *binBuffer);
        }

        static astra_status_t link_connection_to_bin(void* pluginService,
                                                     astra_streamconnection_t connection,
                                                     astra_bin_t binHandle)
        {
            return static_cast<plugin_service*>(pluginService)->link_connection_to_bin(connection, binHandle);
        }

        static astra_status_t get_parameter_bin(void* pluginService,
                                                size_t byteSize,
                                                astra_parameter_bin_t* binHandle,
                                                astra_parameter_data_t* parameterData)
        {
            return static_cast<plugin_service*>(pluginService)->get_parameter_bin(byteSize, *binHandle, *parameterData);
        }

        static astra_status_t log(void* pluginService,
                                  const char* channel,
                                  astra_log_severity_t logLevel,
                                  const char* fileName,
                                  int lineNo,
                                  const char* func,
                                  const char* format,
                                  va_list args)
        {
            return static_cast<plugin_service*>(pluginService)->log(channel, logLevel, fileName, lineNo, func, format, args);
        }
    };
}

#endif /* ASTRA_PLUGIN_SERVICE_DELEGATE_H */
