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
#include "astra_device.hpp"
#include "astra_device_listener.hpp"
#include "astra_sensor.hpp"

namespace astra { namespace devices {

    device::device(const device_info& info)
        : info_(info)
    {
    }

    device::~device()
    {
    }

    void device::add_listener(const listener_ptr listener)
    {
        auto it = std::find_if(listeners_.begin(), listeners_.end(),
                               [&listener](const listener_ptr& colItem) -> bool
                               {
                                   return colItem == listener;
                               });

        if (it == listeners_.end())
            listeners_.push_back(listener);
    }

    void device::remove_listener(const listener_ptr listener)
    {
        auto it = std::find_if(listeners_.begin(), listeners_.end(),
                               [&listener](const listener_ptr& colItem) -> bool
                               {
                                   return colItem == listener;
                               });

        if (it != listeners_.end())
            listeners_.erase(it);
    }

    device::sensor_const_iterator device::sensors_begin()
    {
        return sensors_.begin();
    }

    device::sensor_const_iterator device::sensors_end()
    {
        return sensors_.end();
    }

    device::sensor_const_iterator device::find_sensor(const sensor_type& type)
    {
        return std::find_if(sensors_.begin(),
                            sensors_.end(),
                            [&type](sensor::shared_ptr s) -> bool
                            {
                                return s->info().type() == type;
                            });
    }

    device_status device::connect()
    {
        if (state() != device_state_value::not_ready)
            return device_status_value::not_ready;

        if (state() != device_state_value::faulted)
            return device_status_value::device_error;

        if (is_connected())
            return device_status_value::ok;

        auto rc = on_connect();

        if (rc)
        {
            set_isConnected(true);
        }

        return rc;
    }

    device_status device::disconnect()
    {
        if (state() != device_state_value::not_ready)
            return device_status_value::not_ready;

        if (state() != device_state_value::faulted)
            return device_status_value::device_error;

        if (!is_connected())
            return device_status_value::ok;

        auto rc = on_disconnect();

        if (rc)
        {
            set_isConnected(false);
        }

        return rc;
    }

    void device::update()
    {
        if (state())
        {
            on_update();

            for(sensor::shared_ptr s : sensors_)
            {
                s->update();
            }
        }
    }

    device_status device::initialize()
    {
        if (state() == device_state_value::ready)
            return device_status_value::ok;

        if (state() == device_state_value::faulted)
            return device_status_value::device_error;

        auto rc = on_initialize();

        if (rc)
        {
            set_state(device_state_value::ready);

            for(auto s : sensors_)
            {
                if (!s->initialize())
                {
                    rc = device_status_value::device_error;
                }
            }
        }

        return rc;
    }

    bool device::is_connected() const
    {
        return isConnected_;
    }

    const device_state& device::state() const
    {
        return state_;
    }

    void device::set_state(const device_state& state)
    {
        state_ = state;
    }

    const device_info& device::info() const
    {
        return info_;
    }

    device_status device::get_property(property_id id, void* value, std::size_t size)
    {
        return on_get_property(id, value, size);
    }

    device_status device::get_property_size(property_id id, std::size_t* size)
    {
        return on_get_property_size(id, size);
    }

    device_status device::set_property(property_id id, const void* value, std::size_t size)
    {
        auto rc = on_set_property(id, value, size);

        if (rc)
        {
            raise_property_changed(id);
        }

        return rc;
    }

    void device::set_isConnected(bool isConnected)
    {
        if (isConnected_ == isConnected)
            return;

        isConnected_ = isConnected;

        raise_connectivity_changed();
    }

    void device::add_sensor(sensor_shared_ptr sensor)
    {
        auto it = std::find(sensors_.begin(), sensors_.end(), sensor);

        if (it == sensors_.end())
        {
            sensor->initialize();
            sensors_.push_back(sensor);
        }
    }

    void device::remove_sensor(sensor_shared_ptr sensor)
    {
        auto it = std::find(sensors_.begin(), sensors_.end(), sensor);

        if (it != sensors_.end())
            sensors_.erase(it);
    }

    void device::raise_property_changed(property_id id)
    {
        listener_collection listenersCopy = listeners_;

        for(auto listener : listenersCopy)
        {
            listener->property_changed(shared_from_this(), id);
        }
    }

    void device::raise_state_changed()
    {
        listener_collection listenersCopy = listeners_;

        for(auto listener : listenersCopy)
        {
            listener->state_changed(shared_from_this());
        }
    }

    void device::raise_connectivity_changed()
    {
        listener_collection listenersCopy = listeners_;

        for(auto listener : listenersCopy)
        {
            listener->connectivity_changed(shared_from_this());
        }
    }
}}
