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
#include "astra_sensor_stream.hpp"
#include "astra_stream_listener.hpp"

namespace astra { namespace devices {

    sensor_stream::sensor_stream(const stream_type& type)
        : type_(type)
    {
    }

    sensor_stream::~sensor_stream()
    {
    }

    const sensor_stream::stream_type& sensor_stream::type() const
    {
        return type_;
    }

    void sensor_stream::add_listener(const listener_ptr listener)
    {
        auto it = std::find_if(listeners_.begin(), listeners_.end(),
                               [&listener](const listener_ptr& colItem) -> bool
                               {
                                   return colItem == listener;
                               });

        if (it == listeners_.end())
            listeners_.push_back(listener);
    }

    void sensor_stream::remove_listener(const listener_ptr listener)
    {
        auto it = std::find_if(listeners_.begin(), listeners_.end(),
                               [&listener](const listener_ptr& colItem) -> bool
                               {
                                   return colItem == listener;
                               });

        if (it != listeners_.end())
            listeners_.erase(it);
    }

    device_status sensor_stream::start()
    {
        if (state() == stream_state_value::faulted)
            return device_status_value::device_error;

        if (state() == stream_state_value::not_ready)
            return device_status_value::not_ready;

        auto rc = on_start();

        if (rc)
        {
            set_isStreaming(true);
        }

        return rc;
    }

    device_status sensor_stream::stop()
    {
        if (!state())
            return device_status_value::device_error;

        auto rc = on_stop();

        if (rc)
        {
            set_isStreaming(false);
        }

        return rc;
    }

    device_status sensor_stream::initialize()
    {
        if (state() == stream_state_value::ready)
            return device_status_value::ok;

        if (state() == stream_state_value::faulted)
            return device_status_value::device_error;

        auto rc = on_initialize();

        if (rc)
        {
            set_state(stream_state_value::ready);
        }

        return rc;
    }

    void sensor_stream::update()
    {
        if (state())
        {
            on_update();
        }
    }

    device_status sensor_stream::read_into(void* dest, std::size_t size, std::int32_t timeout)
    {
        if (state() == stream_state_value::faulted)
            return device_status_value::device_error;

        if (!isStreaming_ || state() == stream_state_value::not_ready)
            return device_status_value::not_ready;

        return on_read_into(dest, size, timeout);
    }

    sensor_stream::mode_const_iterator sensor_stream::modes_begin()
    {
        return availableModes_.begin();
    }

    sensor_stream::mode_const_iterator sensor_stream::modes_end()
    {
        return availableModes_.end();
    }

    std::size_t sensor_stream::modes_size() const
    {
        return availableModes_.size();
    }

    const sensor_stream::stream_mode& sensor_stream::active_mode() const
    {
        return *activeMode_;
    }

    bool sensor_stream::is_available() const
    {
        return state() == stream_state_value::ready;
    }

    bool sensor_stream::is_streaming() const
    {
        return state() && isStreaming_;
    }

    stream_state sensor_stream::state() const
    {
        return state_;
    }

    bool sensor_stream::is_property_supported(property_id id) const
    {
        return on_property_supported(id);
    }

    void sensor_stream::set_isStreaming(bool isStreaming)
    {
        if (isStreaming_ == isStreaming)
            return;

        isStreaming_ = isStreaming;
        raise_streaming_changed();
    }

    void sensor_stream::set_state(const stream_state& state)
    {
        if (state_ == state)
            return;

        state_ = state;
        raise_state_changed();
    }

    void sensor_stream::add_available_mode(const stream_mode& mode)
    {
        auto it = std::find(availableModes_.begin(), availableModes_.end(), mode);

        if (it == availableModes_.end())
            availableModes_.push_back(mode);
    }

    void sensor_stream::remove_available_mode(const stream_mode& mode)
    {
        auto it = std::find(availableModes_.begin(), availableModes_.end(), mode);

        if (it != availableModes_.end())
            availableModes_.erase(it);
    }

    device_status sensor_stream::set_active_mode(const stream_mode& mode)
    {
        auto it = std::find(availableModes_.begin(), availableModes_.end(), mode);

        if (it != availableModes_.end())
        {
            activeMode_ = astra::make_unique<stream_mode>(*it);
            on_active_mode_changed();

            return device_status_value::ok;
        }

        return device_status_value::invalid_parameter;
    }

    void sensor_stream::raise_property_changed(property_id id)
    {
        listener_collection listenersCopy = listeners_;

        for(auto listener : listenersCopy)
        {
            listener->property_changed(shared_from_this(), id);
        }
    }

    void sensor_stream::raise_state_changed()
    {
        listener_collection listenersCopy = listeners_;

        for(auto listener : listenersCopy)
        {
            listener->state_changed(shared_from_this());
        }
    }

    void sensor_stream::raise_streaming_changed()
    {
        listener_collection listenersCopy = listeners_;

        for(auto listener : listenersCopy)
        {
            listener->streaming_changed(shared_from_this());
        }
    }

    void sensor_stream::raise_new_frame_available()
    {
        listener_collection listenersCopy = listeners_;

        for(auto listener : listenersCopy)
        {
            listener->new_frame_available(shared_from_this());
        }
    }

    device_status sensor_stream::get_property(property_id id, void* value, std::size_t size)
    {
        return on_get_property(id, value, size);
    }

    device_status sensor_stream::get_property_size(property_id id, std::size_t* size)
    {
        return on_get_property_size(id, size);
    }

    device_status sensor_stream::set_property(property_id id, const void* value, std::size_t size)
    {
        auto rc = on_set_property(id, value, size);

        if (rc)
        {
            raise_property_changed(id);
        }

        return rc;
    }
}}
