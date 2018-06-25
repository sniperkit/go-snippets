/* THIS FILE AUTO-GENERATED FROM astra_context_impl.hpp.lpp. DO NOT EDIT. */
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
#ifndef ASTRA_CONTEXT_IMPL_H
#define ASTRA_CONTEXT_IMPL_H

#include <astra_core/capi/astra_types.h>
#include <atomic>
#include <memory>
#include <unordered_map>
#include <string>

#include "astra_plugin_manager.hpp"
#include "astra_streamset.hpp"
#include "astra_stream_reader.hpp"
#include "astra_shared_library.hpp"
#include "astra_logger.hpp"
#include "astra_streamset_catalog.hpp"

struct astra_streamservice_proxy_t;

namespace astra {

    class context_impl
    {
    public:
        context_impl() = default;
        context_impl& operator=(const context_impl& rhs) = delete;
        context_impl(const context_impl& context) = delete;

        astra_status_t initialize();
        astra_status_t terminate();

        astra_status_t streamset_open(const char* connectionString,
                                      astra_streamsetconnection_t& streamSet);

        astra_status_t streamset_close(astra_streamsetconnection_t& streamSet);

        astra_status_t reader_create(astra_streamsetconnection_t streamSet,
                                     astra_reader_t& reader);

        astra_status_t reader_destroy(astra_reader_t& reader);

        astra_status_t reader_get_stream(astra_reader_t reader,
                                         astra_stream_type_t type,
                                         astra_stream_subtype_t subtype,
                                         astra_streamconnection_t& connection);

        astra_status_t stream_get_description(astra_streamconnection_t connection,
                                              astra_stream_desc_t* description);

        astra_status_t stream_start(astra_streamconnection_t connection);

        astra_status_t stream_stop(astra_streamconnection_t connection);

        astra_status_t reader_open_frame(astra_reader_t reader,
                                         int timeoutMillis,
                                         astra_reader_frame_t& frame);

        astra_status_t reader_close_frame(astra_reader_frame_t& frame);

        astra_status_t reader_register_frame_ready_callback(astra_reader_t reader,
                                                            astra_frame_ready_callback_t callback,
                                                            void* clientTag,
                                                            astra_reader_callback_id_t& callbackId);

        astra_status_t reader_unregister_frame_ready_callback(astra_reader_callback_id_t& callbackId);

        astra_status_t reader_get_frame(astra_reader_frame_t frame,
                                        astra_stream_type_t type,
                                        astra_stream_subtype_t subtype,
                                        astra_frame_t*& subFrame);

        astra_status_t stream_set_parameter(astra_streamconnection_t connection,
                                            astra_parameter_id parameterId,
                                            size_t inByteLength,
                                            astra_parameter_data_t inData);

        astra_status_t stream_get_parameter(astra_streamconnection_t connection,
                                            astra_parameter_id parameterId,
                                            size_t& resultByteLength,
                                            astra_result_token_t& token);

        astra_status_t stream_get_result(astra_streamconnection_t connection,
                                         astra_result_token_t token,
                                         size_t dataByteLength,
                                         astra_parameter_data_t dataDestination);

        astra_status_t stream_invoke(astra_streamconnection_t connection,
                                     astra_command_id commandId,
                                     size_t inByteLength,
                                     astra_parameter_data_t inData,
                                     size_t& resultByteLength,
                                     astra_result_token_t& token);

        astra_status_t temp_update();

        astra_status_t notify_host_event(astra_event_id id, const void* data, size_t dataSize);

    private:
        bool initialized_{false};

        using plugin_manager_ptr = std::unique_ptr<plugin_manager>;
        plugin_manager_ptr pluginManager_;

        using ReaderList = std::vector<stream_reader*>;

        ReaderList activeReaders_;
        streamset_catalog setCatalog_;
    };
}

#endif /* ASTRA_CONTEXT_IMPL_H */
