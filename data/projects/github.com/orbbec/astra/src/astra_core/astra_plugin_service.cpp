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
#include "astra_plugin_service.hpp"
#include "astra_plugin_service_impl.hpp"
#include "astra_create_plugin_proxy.hpp"
#include "astra_cxx_compatibility.hpp"

namespace astra
{
    plugin_service::plugin_service(streamset_catalog& catalog)
        : impl_(astra::make_unique<plugin_service_impl>(catalog)),
          proxy_(create_plugin_proxy(this))
    {}

    plugin_service::~plugin_service()
    {}

    void plugin_service::notify_host_event(astra_event_id id, const void* data, size_t dataSize)
    {
        impl_->notify_host_event(id, data, dataSize);
    }

    astra_pluginservice_proxy_t* plugin_service::proxy()
    {
        return proxy_.get();
    }

   astra_status_t plugin_service::register_stream_registered_callback(stream_registered_callback_t callback,
                                                                      void* clientTag,
                                                                      astra_callback_id_t& callbackId)
   {
       return impl_->register_stream_registered_callback(callback, clientTag, callbackId);
   }

   astra_status_t plugin_service::register_stream_unregistering_callback(stream_unregistering_callback_t callback,
                                                                         void* clientTag,
                                                                         astra_callback_id_t& callbackId)
   {
       return impl_->register_stream_unregistering_callback(callback, clientTag, callbackId);
   }

   astra_status_t plugin_service::register_host_event_callback(host_event_callback_t callback,
                                                               void* clientTag,
                                                               astra_callback_id_t& callbackId)
   {
       return impl_->register_host_event_callback(callback, clientTag, callbackId);
   }

   astra_status_t plugin_service::unregister_host_event_callback(astra_callback_id_t callback)
   {
       return impl_->unregister_host_event_callback(callback);
   }

   astra_status_t plugin_service::unregister_stream_registered_callback(astra_callback_id_t callback)
   {
       return impl_->unregister_stream_registered_callback(callback);
   }

   astra_status_t plugin_service::unregister_stream_unregistering_callback(astra_callback_id_t callback)
   {
       return impl_->unregister_stream_unregistering_callback(callback);
   }

   astra_status_t plugin_service::create_stream_set(const char* setUri,
                                                    astra_streamset_t& setHandle)
   {
       return impl_->create_stream_set(setUri, setHandle);
   }

   astra_status_t plugin_service::destroy_stream_set(astra_streamset_t& setHandle)
   {
       return impl_->destroy_stream_set(setHandle);
   }

   astra_status_t plugin_service::get_streamset_uri(astra_streamset_t setHandle,
                                                    const char*& uri)
   {
       return impl_->get_streamset_uri(setHandle, uri);
   }

   astra_status_t plugin_service::create_stream(astra_streamset_t setHandle,
                                                astra_stream_desc_t desc,
                                                astra_stream_t& handle)
   {
       return impl_->create_stream(setHandle, desc, handle);
   }

   astra_status_t plugin_service::register_stream(astra_stream_t handle,
                                                  stream_callbacks_t pluginCallbacks)
   {
       return impl_->register_stream(handle, pluginCallbacks);
   }

   astra_status_t plugin_service::unregister_stream(astra_stream_t handle)
   {
       return impl_->unregister_stream(handle);
   }

   astra_status_t plugin_service::destroy_stream(astra_stream_t& handle)
   {
       return impl_->destroy_stream(handle);
   }

   astra_status_t plugin_service::create_stream_bin(astra_stream_t streamHandle,
                                                    size_t lengthInBytes,
                                                    astra_bin_t& binHandle,
                                                    astra_frame_t*& binBuffer)
   {
       return impl_->create_stream_bin(streamHandle, lengthInBytes, binHandle, binBuffer);
   }

   astra_status_t plugin_service::destroy_stream_bin(astra_stream_t streamHandle,
                                                     astra_bin_t& binHandle,
                                                     astra_frame_t*& binBuffer)
   {
       return impl_->destroy_stream_bin(streamHandle, binHandle, binBuffer);
   }

   astra_status_t plugin_service::bin_has_connections(astra_bin_t binHandle,
                                                      bool& hasConnections)
   {
       return impl_->bin_has_connections(binHandle, hasConnections);
   }

   astra_status_t plugin_service::cycle_bin_buffers(astra_bin_t binHandle,
                                                    astra_frame_t*& binBuffer)
   {
       return impl_->cycle_bin_buffers(binHandle, binBuffer);
   }

   astra_status_t plugin_service::link_connection_to_bin(astra_streamconnection_t connection,
                                                         astra_bin_t binHandle)
   {
       return impl_->link_connection_to_bin(connection, binHandle);
   }

   astra_status_t plugin_service::get_parameter_bin(size_t byteSize,
                                                    astra_parameter_bin_t& binHandle,
                                                    astra_parameter_data_t& parameterData)
   {
       return impl_->get_parameter_bin(byteSize, binHandle, parameterData);
   }

   astra_status_t plugin_service::log(const char* channel,
                                      astra_log_severity_t logLevel,
                                      const char* fileName,
                                      int lineNo,
                                      const char* func,
                                      const char* format,
                                      va_list args)
   {
       return impl_->log(channel, logLevel, fileName, lineNo, func, format, args);
   }


}
