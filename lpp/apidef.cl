;; This file is part of the Orbbec Astra SDK [https://orbbec3d.com]
;; Copyright (c) 2015 Orbbec 3D
;;
;; Licensed under the Apache License, Version 2.0 (the "License");
;; you may not use this file except in compliance with the License.
;; You may obtain a copy of the License at
;;
;; http://www.apache.org/licenses/LICENSE-2.0
;;
;; Unless required by applicable law or agreed to in writing, software
;; distributed under the License is distributed on an "AS IS" BASIS,
;; WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
;; See the License for the specific language governing permissions and
;; limitations under the License.
;;
;; Be excellent to each other.

;; definitions manually adapted from astra_capi.h

(add-void-param :funcset "stream"
                :type "void*"
                :name "streamService")

(add-void-param :funcset "plugin"
                :type "void*"
                :name "pluginService")

(add-void-param :funcset "plugincallbacks"
                :type "void*"
                :name "instance")

;; void set_parameter(astra_streamconnection_t connection,
;;                    astra_parameter_id id,
;;                    size_t inByteLength,
;;                    astra_parameter_data_t inData)
(add-func       :funcset "plugincallbacks"
                :returntype "void"
                :funcname "set_parameter"
                :params (list (make-param :type "astra_streamconnection_t" :name "connection")
                              (make-param :type "astra_parameter_id" :name "id")
                              (make-param :type "size_t" :name "inByteLength")
                              (make-param :type "astra_parameter_data_t" :name "inData")))

;; void get_parameter(astra_streamconnection_t connection,
;;                    astra_parameter_id id,
;;                    astra_parameter_bin_t* parameterBin)
(add-func       :funcset "plugincallbacks"
                :returntype "void"
                :funcname "get_parameter"
                :params (list (make-param :type "astra_streamconnection_t" :name "connection")
                              (make-param :type "astra_parameter_id" :name "id")
                              (make-param :type "astra_parameter_bin_t*" :name "parameterBin" :deref T)))

;; void invoke(astra_streamconnection_t connection,
;;             astra_command_id commandId,
;;             size_t inByteLength,
;;             astra_parameter_data_t inData,
;;             astra_parameter_bin_t* parameterBin)
(add-func       :funcset "plugincallbacks"
                :returntype "void"
                :funcname "invoke"
                :params (list (make-param :type "astra_streamconnection_t" :name "connection")
                              (make-param :type "astra_command_id" :name "commandId")
                              (make-param :type "size_t" :name "inByteLength")
                              (make-param :type "astra_parameter_data_t" :name "inData")
                              (make-param :type "astra_parameter_bin_t*" :name "parameterBin" :deref T)))

;; void connection_added(astra_stream_t stream,
;;                       astra_streamconnection_t connection)
(add-func       :funcset "plugincallbacks"
                :returntype "void"
                :funcname "connection_added"
                :params (list (make-param :type "astra_stream_t" :name "stream")
                              (make-param :type "astra_streamconnection_t" :name "connection")))

;; void connection_removed(astra_stream_t stream,
;;                         astra_bin_t bin,
;;                         astra_streamconnection_t connection)
(add-func       :funcset "plugincallbacks"
                :returntype "void"
                :funcname "connection_removed"
                :params (list (make-param :type "astra_stream_t" :name "stream")
                              (make-param :type "astra_bin_t" :name "bin")
                              (make-param :type "astra_streamconnection_t" :name "connection")))

;; void connection_started(astra_stream_t stream,
;;                         astra_streamconnection_t connection)
(add-func       :funcset "plugincallbacks"
                :returntype "void"
                :funcname "connection_started"
                :params (list (make-param :type "astra_stream_t" :name "stream")
                              (make-param :type "astra_streamconnection_t" :name "connection")))

;; void connection_stopped(astra_stream_t stream,
;;                         astra_streamconnection_t connection)
(add-func       :funcset "plugincallbacks"
                :returntype "void"
                :funcname "connection_stopped"
                :params (list (make-param :type "astra_stream_t" :name "stream")
                              (make-param :type "astra_streamconnection_t" :name "connection")))

;; astra_status_t register_stream_added_callback(stream_added_callback_t callback,
;;                                                  void* clientTag,
;;                                                  astra_callback_id_t* callbackId)
(add-func       :funcset "plugin"
                :returntype "astra_status_t"
                :funcname "register_stream_registered_callback"
                :params (list (make-param :type "stream_registered_callback_t" :name "callback")
                              (make-param :type "void*" :name "clientTag")
                              (make-param :type "astra_callback_id_t*" :name "callbackId" :deref T)))

;; astra_status_t register_stream_removing_callback(stream_removing_callback_t callback, void* clientTag, astra_callback_id_t* callbackId)
(add-func       :funcset "plugin"
                :returntype "astra_status_t"
                :funcname "register_stream_unregistering_callback"
                :params (list (make-param :type "stream_unregistering_callback_t" :name "callback")
                              (make-param :type "void*" :name "clientTag")
                              (make-param :type "astra_callback_id_t*" :name "callbackId" :deref T)))

;; astra_status_t register_host_event_callback(host_event_callback_t callback,
;;                                                void* clientTag,
;;                                                astra_callback_id_t* callbackId)
(add-func       :funcset "plugin"
                :returntype "astra_status_t"
                :funcname "register_host_event_callback"
                :params (list (make-param :type "host_event_callback_t" :name "callback")
                              (make-param :type "void*" :name "clientTag")
                              (make-param :type "astra_callback_id_t*" :name "callbackId" :deref T)))

;; astra_status_t unregister_host_event_callback(astra_callback_id_t callbackId)
(add-func       :funcset "plugin"
                :returntype "astra_status_t"
                :funcname "unregister_host_event_callback"
                :params (list (make-param :type "astra_callback_id_t" :name "callback")))

;; astra_status_t unregister_stream_added_callback(astra_callback_id_t callbackId)
(add-func       :funcset "plugin"
                :returntype "astra_status_t"
                :funcname "unregister_stream_registered_callback"
                :params (list (make-param :type "astra_callback_id_t" :name "callback")))

;; astra_status_t unregister_stream_removing_callback(astra_callback_id_t callbackId)
(add-func       :funcset "plugin"
                :returntype "astra_status_t"
                :funcname "unregister_stream_unregistering_callback"
                :params (list (make-param :type "astra_callback_id_t" :name "callback")))

;; astra_status_t create_stream_set(astra_streamset_t& setHandle)
(add-func       :funcset "plugin"
                :returntype "astra_status_t"
                :funcname "create_stream_set"
                :params (list (make-param :type "const char*" :name "setUri")
                              (make-param :type "astra_streamset_t&" :name "setHandle")))

;; astra_status_t destroy_stream_set(astra_streamset_t& setHandle)
(add-func       :funcset "plugin"
                :returntype "astra_status_t"
                :funcname "destroy_stream_set"
                :params (list (make-param :type "astra_streamset_t&" :name "setHandle")))

;; astra_status_t get_streamset_uri(astra_streamset_t setHandle,
;;                                     const char* uri)
(add-func       :funcset "plugin"
                :returntype "astra_status_t"
                :funcname "get_streamset_uri"
                :params (list
                         (make-param :type "astra_streamset_t" :name "setHandle")
                         (make-param :type "const char**" :name "uri" :deref T)))

;; astra_status_t create_stream(astra_streamset_t setHandle,
;;                              astra_stream_desc_t desc,
;;                              stream_callbacks_t pluginCallbacks,
;;                              astra_stream_t* handle)
(add-func       :funcset "plugin"
                :returntype "astra_status_t"
                :funcname "create_stream"
                :params (list (make-param :type "astra_streamset_t" :name "setHandle")
                              (make-param :type "astra_stream_desc_t" :name "desc")
                              (make-param :type "astra_stream_t*" :name "handle" :deref T)))

;; astra_status_t register_stream(astra_stream_t& handle)
(add-func       :funcset "plugin"
                :returntype "astra_status_t"
                :funcname "register_stream"
                :params (list (make-param :type "astra_stream_t" :name "handle")
                              (make-param :type "stream_callbacks_t" :name "pluginCallbacks")))

;; astra_status_t unregister_stream(astra_stream_t handle)
(add-func       :funcset "plugin"
                :returntype "astra_status_t"
                :funcname "unregister_stream"
                :params (list (make-param :type "astra_stream_t" :name "handle")))

;; astra_status_t destroy_stream(astra_stream_t& handle)
(add-func       :funcset "plugin"
                :returntype "astra_status_t"
                :funcname "destroy_stream"
                :params (list (make-param :type "astra_stream_t&" :name "handle")))

;; astra_status_t create_stream_bin(astra_stream_t streamHandle,
;;                                     size_t lengthInBytes,
;;                                     astra_bin_t* binHandle,
;;                                     astra_frame_t** binBuffer)
(add-func       :funcset "plugin"
                :returntype "astra_status_t"
                :funcname "create_stream_bin"
                :params (list (make-param :type "astra_stream_t" :name "streamHandle")
                              (make-param :type "size_t" :name "lengthInBytes")
                              (make-param :type "astra_bin_t*" :name "binHandle" :deref t)
                              (make-param :type "astra_frame_t**" :name "binBuffer" :deref t)))

;; astra_status_t destroy_stream_bin(astra_stream_t handle,
;;                                      astra_bin_t* binHandle,
;;                                      astra_frame_t** binBuffer)
(add-func       :funcset "plugin"
                :returntype "astra_status_t"
                :funcname "destroy_stream_bin"
                :params (list (make-param :type "astra_stream_t" :name "streamHandle")
                              (make-param :type "astra_bin_t*" :name "binHandle" :deref t)
                              (make-param :type "astra_frame_t**" :name "binBuffer" :deref t)))

;; astra_status_t bin_has_connections(astra_bin_t binHandle,
;;                                             bool* hasConnections)
(add-func       :funcset "plugin"
                :returntype "astra_status_t"
                :funcname "bin_has_connections"
                :params (list (make-param :type "astra_bin_t" :name "binHandle")
                              (make-param :type "bool*" :name "hasConnections" :deref t)))

;; astra_status_t cycle_bin_buffers(astra_bin_t binHandle,
;;                                     astra_frame_t** binBuffer)
(add-func       :funcset "plugin"
                :returntype "astra_status_t"
                :funcname "cycle_bin_buffers"
                :params (list (make-param :type "astra_bin_t" :name "binHandle")
                              (make-param :type "astra_frame_t**" :name "binBuffer" :deref t)))

;; astra_status_t link_connection_to_bin(astra_streamconnection_t connection,
;;                                          astra_bin_t binHandle)
(add-func       :funcset "plugin"
                :returntype "astra_status_t"
                :funcname "link_connection_to_bin"
                :params (list (make-param :type "astra_streamconnection_t" :name "connection")
                              (make-param :type "astra_bin_t" :name "binHandle")))

;; astra_status_t get_parameter_bin(size_t byteSize,
;;                                     astra_parameter_bin_t* binHandle,
;;                                     astra_parameter_data_t** parameterData)
(add-func       :funcset "plugin"
                :returntype "astra_status_t"
                :funcname "get_parameter_bin"
                :params (list (make-param :type "size_t" :name "byteSize")
                              (make-param :type "astra_parameter_bin_t*" :name "binHandle" :deref t)
                              (make-param :type "astra_parameter_data_t*" :name "parameterData" :deref t)))

;; astra_status_t log(astra_log_severity_t logLeve,
;;                       const char* format,
;;                       va_list args)
(add-func       :funcset "plugin"
                :returntype "astra_status_t"
                :funcname "log"
                :params (list (make-param :type "const char*" :name "channel")
                              (make-param :type "astra_log_severity_t" :name "logLevel")
                              (make-param :type "const char*" :name "fileName")
                              (make-param :type "int" :name "lineNo")
                              (make-param :type "const char*" :name "func")
                              (make-param :type "const char*" :name "format")
                              (make-param :type "va_list" :name "args")))

;; ASTRA_API astra_status_t astra_initialize();
;; (add-func       :funcset "stream"
;;                 :returntype "astra_status_t"
;;                 :funcname "initialize"
;;                 :params '())

;; ASTRA_API astra_status_t astra_terminate();
;; (add-func       :funcset "stream"
;;                 :returntype "astra_status_t"
;;                 :funcname "terminate"
;;                 :params '())

;;ASTRA_API astra_status_t astra_streamset_open(const char* connectionString,
;;                                                       astra_streamsetconnection_t* streamSet);
(add-func       :funcset "stream"
                :returntype "astra_status_t"
                :funcname "streamset_open"
                :params (list (make-param :type "const char*" :name "connectionString")
                              (make-param :type "astra_streamsetconnection_t*" :name "streamSet" :deref T)))

;; ASTRA_API astra_status_t astra_streamset_close(astra_streamset_connection_t* streamSet);
(add-func       :funcset "stream"
                :returntype "astra_status_t"
                :funcname "streamset_close"
                :params (list (make-param :type "astra_streamsetconnection_t*" :name "streamSet" :deref T)))

;; ASTRA_API char* astra_get_status_string(astra_status_t status);
;; (add-func       :funcset "stream"
;;                 :returntype "astra_status_t"
;;                 :funcname "get_status_string"
;;                 :params (list (make-param :type "astra_status_t" :name "status")))

;; ASTRA_API astra_status_t astra_reader_create(astra_streamsetconnection_t streamSet,
;;                                                       astra_reader_t* reader);
(add-func       :funcset "stream"
                :returntype "astra_status_t"
                :funcname "reader_create"
                :params (list (make-param :type "astra_streamsetconnection_t" :name "streamSet")
                              (make-param :type "astra_reader_t*" :name "reader" :deref T)))

;; ASTRA_API astra_status_t astra_reader_destroy(astra_reader_t* reader);
(add-func       :funcset "stream"
                :returntype "astra_status_t"
                :funcname "reader_destroy"
                :params (list (make-param :type "astra_reader_t*" :name "reader" :deref T)))

;; ASTRA_API astra_status_t astra_reader_get_stream(astra_reader_t reader,
;;                                                           astra_stream_type_t type,
;;                                                           astra_stream_subtype_t subtype,
;;                                                           astra_streamconnection_t* connection);
(add-func       :funcset "stream"
                :returntype "astra_status_t"
                :funcname "reader_get_stream"
                :params (list (make-param :type "astra_reader_t" :name "reader")
                              (make-param :type "astra_stream_type_t" :name "type")
                              (make-param :type "astra_stream_subtype_t" :name "subtype")
                              (make-param :type "astra_streamconnection_t*" :name "connection" :deref T)))

;; ASTRA_API astra_status_t astra_stream_get_description(astra_streamconnection_t connection,
;;                                                                astra_stream_desc_t* description);
(add-func       :funcset "stream"
                :returntype "astra_status_t"
                :funcname "stream_get_description"
                :params (list (make-param :type "astra_streamconnection_t" :name "connection")
                              (make-param :type "astra_stream_desc_t*" :name "description")))

;; ASTRA_API astra_status_t astra_stream_start(astra_streamconnection_t connection);
(add-func       :funcset "stream"
                :returntype "astra_status_t"
                :funcname "stream_start"
                :params (list (make-param :type "astra_streamconnection_t" :name "connection")))

;; ASTRA_API astra_status_t astra_stream_stop(astra_streamconnection_t connection);
(add-func       :funcset "stream"
                :returntype "astra_status_t"
                :funcname "stream_stop"
                :params (list (make-param :type "astra_streamconnection_t" :name "connection")))

;; ASTRA_API astra_status_t astra_reader_open_frame(astra_reader_t reader,
;;                                                           int timeoutMillis,
;;                                                           astra_reader_frame_t* frame); //0 = return immediately
(add-func       :funcset "stream"
                :returntype "astra_status_t"
                :funcname "reader_open_frame"
                :params (list (make-param :type "astra_reader_t" :name "reader")
                              (make-param :type "int" :name "timeoutMillis")
                              (make-param :type "astra_reader_frame_t*" :name "frame" :deref T)))

;; ASTRA_API astra_status_t astra_reader_close_frame(astra_reader_frame_t* frame); //frame set to null
(add-func       :funcset "stream"
                :returntype "astra_status_t"
                :funcname "reader_close_frame"
                :params (list (make-param :type "astra_reader_frame_t*" :name "frame" :deref T)))

;; ASTRA_API astra_status_t reader_register_frame_ready_callback(astra_reader_t reader,
;;                                                                     astra_frame_ready_callback_t callback,
;;                                                                     void* clientTag,
;;                                                                     astra_reader_callback_id_t* callbackId)
(add-func       :funcset "stream"
                :returntype "astra_status_t"
                :funcname "reader_register_frame_ready_callback"
                :params (list (make-param :type "astra_reader_t" :name "reader")
                              (make-param :type "astra_frame_ready_callback_t" :name "callback")
                              (make-param :type "void*" :name "clientTag")
                              (make-param :type "astra_reader_callback_id_t*" :name "callbackId" :deref T)))

;; ASTRA_API astra_status_t reader_unregister_frame_ready_callback(astra_reader_callback_id_t* callbackId)
(add-func       :funcset "stream"
                :returntype "astra_status_t"
                :funcname "reader_unregister_frame_ready_callback"
                :params (list (make-param :type "astra_reader_callback_id_t*" :name "callbackId" :deref T)))

;; ASTRA_API astra_status_t astra_reader_get_frame(astra_reader_frame_t frame,
;;                                                          astra_stream_type_t type,
;;                                                          astra_stream_subtype_t subtype,
;;                                                          astra_frame_t** frame);
(add-func       :funcset "stream"
                :returntype "astra_status_t"
                :funcname "reader_get_frame"
                :params (list (make-param :type "astra_reader_frame_t" :name "frame")
                              (make-param :type "astra_stream_type_t" :name "type")
                              (make-param :type "astra_stream_subtype_t" :name "subtype")
                              (make-param :type "astra_frame_t**" :name "subFrame" :deref T)))

;; ASTRA_API astra_status_t astra_stream_set_parameter(astra_streamconnection_t connection,
;;                                                              astra_parameter_id parameterId,
;;                                                              size_t inByteLength,
;;                                                              astra_parameter_data_t inData);
(add-func       :funcset "stream"
                :returntype "astra_status_t"
                :funcname "stream_set_parameter"
                :params (list (make-param :type "astra_streamconnection_t" :name "connection")
                              (make-param :type "astra_parameter_id" :name "parameterId")
                              (make-param :type "size_t" :name "inByteLength")
                              (make-param :type "astra_parameter_data_t" :name "inData")))

;; ASTRA_API astra_status_t astra_stream_get_parameter(astra_streamconnection_t connection,
;;                                                              astra_parameter_id parameterId,
;;                                                              size_t* resultByteLength,
;;                                                              astra_result_token_t* token);
(add-func       :funcset "stream"
                :returntype "astra_status_t"
                :funcname "stream_get_parameter"
                :params (list (make-param :type "astra_streamconnection_t" :name "connection")
                              (make-param :type "astra_parameter_id" :name "parameterId")
                              (make-param :type "size_t*" :name "resultByteLength" :deref T)
                              (make-param :type "astra_result_token_t*" :name "token" :deref T)))

;; ASTRA_API astra_status_t astra_stream_get_result(astra_streamconnection_t connection,
;;                                                                   astra_result_token_t parameterId,
;;                                                                   size_t dataByteLength,
;;                                                                   astra_parameter_data_t dataDestination);
(add-func       :funcset "stream"
                :returntype "astra_status_t"
                :funcname "stream_get_result"
                :params (list (make-param :type "astra_streamconnection_t" :name "connection")
                              (make-param :type "astra_result_token_t" :name "token")
                              (make-param :type "size_t" :name "dataByteLength")
                              (make-param :type "astra_parameter_data_t" :name "dataDestination")))


;; ASTRA_API astra_status_t astra_stream_invoke(astra_streamconnection_t connection,
;;                                                       astra_command_id commandId,
;;                                                       size_t inByteLength,
;;                                                       astra_parameter_data_t inData,
;;                                                       astra_result_token_t* token);
(add-func       :funcset "stream"
                :returntype "astra_status_t"
                :funcname "stream_invoke"
                :params (list (make-param :type "astra_streamconnection_t" :name "connection")
                              (make-param :type "astra_command_id" :name "commandId")
                              (make-param :type "size_t" :name "inByteLength")
                              (make-param :type "astra_parameter_data_t" :name "inData")
                              (make-param :type "size_t*" :name "resultByteLength" :deref T)
                              (make-param :type "astra_result_token_t*" :name "token" :deref T)))

;; ASTRA_API astra_status_t astra_temp_update();
(add-func       :funcset "stream"
                :returntype "astra_status_t"
                :funcname "temp_update"
                :params '())
