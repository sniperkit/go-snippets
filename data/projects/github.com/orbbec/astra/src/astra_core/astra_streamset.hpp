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
#ifndef ASTRA_STREAMSET_H
#define ASTRA_STREAMSET_H

#include "astra_stream.hpp"
#include "astra_logger.hpp"
#include "astra_signal.hpp"
#include "astra_stream_registered_event_args.hpp"
#include "astra_stream_unregistering_event_args.hpp"
#include <unordered_set>
#include <functional>
#include <vector>
#include <memory>
#include <string>
#include <cassert>
#include "astra_streamset_connection.hpp"
#include "astra_stream_listener.hpp"

namespace astra {

    class streamset : public stream_listener
    {
    public:
        using VisitorFunc = std::function<void(stream*)>;

        streamset(std::string uri);

        streamset& operator=(const streamset& rhs) = delete;
        streamset(const streamset& streamSet) = delete;

        streamset_connection* add_new_connection();
        void disconnect_streamset_connection(streamset_connection* connection);

        stream_connection* create_stream_connection(const astra_stream_desc_t& desc);
        bool destroy_stream_connection(stream_connection* connection);

        void claim();
        bool is_claimed() const;

        //plugins only below

        astra::stream* register_stream(astra_stream_desc_t desc);
        void destroy_stream(stream* stream);

        void claim_stream(stream* stream, stream_callbacks_t callbacks);

        bool is_member(astra_stream_t stream) const;

        astra_stream_t find_stream_by_type_subtype(astra_stream_type_t type,
                                                      astra_stream_subtype_t subtype) const;

        const std::string& get_uri() { return uri_; }

        astra_streamset_t get_handle() { return reinterpret_cast<astra_streamset_t>(this); }

        void visit_streams(VisitorFunc visitorMethod);

        static streamset* get_ptr(astra_streamset_t handle) { return reinterpret_cast<streamset*>(handle); }

        astra_callback_id_t register_for_stream_registered_event(StreamRegisteredCallback callback);
        astra_callback_id_t register_for_stream_unregistering_event(StreamUnregisteringCallback callback);
        void unregister_for_stream_registered_event(astra_callback_id_t callbackId);
        void unregister_for_stream_unregistering_event(astra_callback_id_t callbackId);

        virtual void on_stream_registered(stream* stream) override;
        virtual void on_stream_unregistering(stream* stream) override;

    private:
        stream* find_stream_by_type_subtype_impl(astra_stream_type_t type,
                                                 astra_stream_subtype_t subtype) const;

        using StreamCollection = std::unordered_set<stream*>;
        StreamCollection streamCollection_;

        using streamset_connectionPtr = std::unique_ptr<streamset_connection>;
        using streamset_connectionList = std::vector<streamset_connectionPtr>;
        streamset_connectionList connections_;

        bool isClaimed_{false};

        std::string uri_;

        signal<stream_registered_event_args> streamRegisteredSignal_;
        signal<stream_unregistering_event_args> streamUnregisteringSignal_;
    };
}

#endif /* ASTRA_STREAMSET_H */
