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
#ifndef ASTRA_SENSOR_STREAM_H
#define ASTRA_SENSOR_STREAM_H

#include <vector>
#include <memory>

#include <astra/streams/Image.hpp>
#include "astra_device_status.hpp"
#include "astra_stream_state.hpp"

namespace astra { namespace devices {

    enum class stream_types
    {
        unknown = 0,
        depth = 1,
        color = 2,
        infrared = 3
    };

    enum property_keys
    {
        ACTIVE_STREAM_MODE = 1
    };

    template<typename T>
    std::uint32_t streamtype_cast(T p);

    template<>
    inline std::uint32_t streamtype_cast<stream_types>(stream_types type)
    {
        return static_cast<std::uint32_t>(type);
    }

    class stream_listener;

    class sensor_stream : public std::enable_shared_from_this<sensor_stream>
    {
    public:
        using shared_ptr = std::shared_ptr<sensor_stream>;
        using stream_mode = astra::ImageStreamMode;
        using mode_collection =  std::vector<stream_mode>;
        using mode_const_iterator = mode_collection::const_iterator;
        using listener_ptr = stream_listener*;
        using property_id = std::uint32_t;
        using stream_type = std::uint32_t;

        sensor_stream(const stream_type& type);

        sensor_stream(const sensor_stream&) = delete;
        const sensor_stream& operator=(const sensor_stream&) = delete;

        virtual ~sensor_stream();

        void add_listener(const listener_ptr listener);
        void remove_listener(const listener_ptr listener);

        device_status initialize();
        device_status start();
        device_status stop();
        device_status read_into(void* dest, std::size_t size, std::int32_t timeout);

        mode_const_iterator modes_begin();
        mode_const_iterator modes_end();
        std::size_t modes_size() const;

        const stream_mode& active_mode() const;

        bool is_available() const;
        bool is_streaming() const;

        stream_state state() const;
        const stream_type& type() const;

        template<typename T>
        device_status set_property(property_id id, const T& value);

        template<typename T>
        device_status get_property(property_id id, T* value);

        device_status get_property(property_id id, void* value, std::size_t size);
        device_status get_property_size(property_id id, std::size_t* size);
        device_status set_property(property_id id, const void* value, std::size_t size);

        bool is_property_supported(property_id id) const;

        void update();

    protected:
        void set_isStreaming(bool isStreaming);
        void set_state(const stream_state& state);

        void add_available_mode(const stream_mode& mode);
        void remove_available_mode(const stream_mode& mode);

        device_status set_active_mode(const stream_mode& mode);
        virtual void on_active_mode_changed() {}

        void raise_property_changed(property_id id);
        void raise_new_frame_available();

    private:
        using listener_collection = std::vector<listener_ptr>;

        virtual device_status on_initialize() = 0;
        virtual device_status on_start() = 0;
        virtual device_status on_stop() = 0;
        virtual void on_update() {}

        virtual device_status on_read_into(void* dest, std::size_t size, std::int32_t timeout) = 0;

        virtual device_status on_get_property(property_id id, void* value, std::size_t size)
        {
            return device_status_value::invalid_parameter;
        }

        virtual device_status on_get_property_size(property_id id, std::size_t* size)
        {
            return device_status_value::invalid_parameter;
        }

        virtual device_status on_set_property(property_id id, const void* value, std::size_t size)
        {
            return device_status_value::invalid_parameter;
        }

        virtual bool on_property_supported(property_id id) const
        {
            return false;
        }

        void raise_state_changed();
        void raise_streaming_changed();

        stream_state state_{stream_state_value::not_ready};
        bool isStreaming_{false};

        stream_type type_;
        std::unique_ptr<stream_mode> activeMode_;
        mode_collection availableModes_;
        listener_collection listeners_;
    };

    template<typename T>
    device_status sensor_stream::set_property(property_id id, const T& value)
    {
        return set_property(id, &value, sizeof(T));
    }

    template<typename T>
    device_status sensor_stream::get_property(property_id id, T* value)
    {
        return get_property(id, value, sizeof(T));
    }
}}

#endif /* ASTRA_SENSOR_STREAM_H */
