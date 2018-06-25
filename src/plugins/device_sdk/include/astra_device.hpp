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
#ifndef ASTRA_DEVICE_H
#define ASTRA_DEVICE_H

#include <vector>
#include <memory>

#include "astra_device_info.hpp"
#include "astra_device_state.hpp"
#include "astra_device_status.hpp"
#include "astra_sensor_info.hpp"

namespace astra { namespace devices {

    class device_listener;
    class sensor;

    class device : public std::enable_shared_from_this<device>
    {
    public:
        using sensor_shared_ptr = std::shared_ptr<sensor>;
        using listener_ptr = device_listener*;
        using shared_ptr = std::shared_ptr<device>;
        using sensor_collection = std::vector<sensor_shared_ptr>;
        using sensor_const_iterator = sensor_collection::const_iterator;
        using property_id = std::uint32_t;

        device(const device_info& info);

        device(const device&) = delete;
        const device& operator=(const device&) = delete;

        virtual ~device();

        void add_listener(const listener_ptr listener);
        void remove_listener(const listener_ptr listener);

        sensor_const_iterator sensors_begin();
        sensor_const_iterator sensors_end();
        sensor_const_iterator find_sensor(const sensor_type& type);

        device_status initialize();
        device_status connect();
        device_status disconnect();

        bool is_connected() const;

        const device_state& state() const;
        const device_info& info() const;

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
        void set_isConnected(bool isConnected);
        void add_sensor(sensor_shared_ptr sensor);
        void remove_sensor(sensor_shared_ptr sensor);
        void set_state(const device_state& state);

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
        virtual device_status on_connect() = 0;
        virtual device_status on_disconnect() = 0;
        virtual void on_update() {}

        void raise_state_changed();
        void raise_connectivity_changed();

        bool isConnected_{false};
        device_info info_;
        device_state state_{device_state_value::not_ready};
        sensor_collection sensors_;
        listener_collection listeners_;
    };

    template<typename T>
    device_status device::set_property(property_id id, const T& value)
    {
        return set_property(id, &value, sizeof(T));
    }

    template<typename T>
    device_status device::get_property(property_id id, T* value)
    {
        return get_property(id, value, sizeof(T));
    }
}}

#endif /* ASTRA_DEVICE_H */
