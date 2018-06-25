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
#ifndef ASTRA_SENSOR_H
#define ASTRA_SENSOR_H

#include <vector>
#include <memory>

#include "astra_device_status.hpp"
#include "astra_sensor_info.hpp"
#include "astra_sensor_state.hpp"
#include "astra_sensor_stream.hpp"

namespace astra { namespace devices {

//    class sensor_stream;
    class sensor_listener;

    class sensor : public std::enable_shared_from_this<sensor>
    {
    public:
        using stream_shared_ptr = std::shared_ptr<sensor_stream>;
        using stream_collection = std::vector<stream_shared_ptr>;
        using shared_ptr = std::shared_ptr<sensor>;
        using stream_const_iterator = stream_collection::const_iterator;
        using listener_ptr = sensor_listener*;
        using property_id = std::uint32_t;

        sensor(const sensor_info& info);

        sensor(const sensor&) = delete;
        const sensor& operator=(const sensor&) = delete;

        virtual ~sensor();

        device_status initialize();

        void add_listener(const listener_ptr listener);
        void remove_listener(const listener_ptr listener);

        stream_const_iterator streams_begin();
        stream_const_iterator streams_end();
        stream_const_iterator find_stream(sensor_stream::stream_type type);

        bool is_available() const;

        sensor_state state() const;
        const sensor_info& info() const;

        void update();

        template<typename T>
        device_status set_property(property_id id, const T& value);

        template<typename T>
        device_status get_property(property_id id, T* value);

        device_status get_property(property_id id, void* value, std::size_t size);
        device_status get_property_size(property_id id, std::size_t* size);
        device_status set_property(property_id id, const void* value, std::size_t size);

        bool is_property_supported(property_id id) const;

    protected:
        void set_state(const sensor_state& state);

        void add_stream(stream_shared_ptr stream);
        void remove_stream(stream_shared_ptr stream);

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

        void raise_property_changed(property_id id);

    private:
        using listener_collection = std::vector<listener_ptr>;

        virtual device_status on_initialize() = 0;
        virtual void on_update() {}

        void raise_state_changed();

        sensor_state state_{sensor_state_value::not_ready};
        sensor_info info_;

        stream_collection streams_;
        listener_collection listeners_;
    };

    template<typename T>
    device_status sensor::set_property(property_id id, const T& value)
    {
        return set_property(id, &value, sizeof(T));
    }

    template<typename T>
    device_status sensor::get_property(property_id id, T* value)
    {
        return get_property(id, value, sizeof(T));
    }
}}

#endif /* ASTRA_SENSOR_H */
