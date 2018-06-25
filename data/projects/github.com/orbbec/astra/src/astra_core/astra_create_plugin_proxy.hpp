/* THIS FILE AUTO-GENERATED FROM astra_create_plugin_proxy.hpp.lpp. DO NOT EDIT. */
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
#ifndef ASTRA_CREATE_PLUGIN_PROXY_H
#define ASTRA_CREATE_PLUGIN_PROXY_H

#include "astra_plugin_service.hpp"
#include "astra_plugin_service_delegate.hpp"
#include <astra_core/capi/plugins/astra_pluginservice_proxy.h>

namespace astra {
astra_pluginservice_proxy_t* create_plugin_proxy(plugin_service* service)
    {
        astra_pluginservice_proxy_t* proxy = new astra_pluginservice_proxy_t;

        proxy->register_stream_registered_callback = &plugin_service_delegate::register_stream_registered_callback;
        proxy->register_stream_unregistering_callback = &plugin_service_delegate::register_stream_unregistering_callback;
        proxy->register_host_event_callback = &plugin_service_delegate::register_host_event_callback;
        proxy->unregister_host_event_callback = &plugin_service_delegate::unregister_host_event_callback;
        proxy->unregister_stream_registered_callback = &plugin_service_delegate::unregister_stream_registered_callback;
        proxy->unregister_stream_unregistering_callback = &plugin_service_delegate::unregister_stream_unregistering_callback;
        proxy->create_stream_set = &plugin_service_delegate::create_stream_set;
        proxy->destroy_stream_set = &plugin_service_delegate::destroy_stream_set;
        proxy->get_streamset_uri = &plugin_service_delegate::get_streamset_uri;
        proxy->create_stream = &plugin_service_delegate::create_stream;
        proxy->register_stream = &plugin_service_delegate::register_stream;
        proxy->unregister_stream = &plugin_service_delegate::unregister_stream;
        proxy->destroy_stream = &plugin_service_delegate::destroy_stream;
        proxy->create_stream_bin = &plugin_service_delegate::create_stream_bin;
        proxy->destroy_stream_bin = &plugin_service_delegate::destroy_stream_bin;
        proxy->bin_has_connections = &plugin_service_delegate::bin_has_connections;
        proxy->cycle_bin_buffers = &plugin_service_delegate::cycle_bin_buffers;
        proxy->link_connection_to_bin = &plugin_service_delegate::link_connection_to_bin;
        proxy->get_parameter_bin = &plugin_service_delegate::get_parameter_bin;
        proxy->log = &plugin_service_delegate::log;
        proxy->pluginService = service;

        return proxy;
    }
}

#endif /* ASTRA_CREATE_PLUGIN_PROXY_H */
