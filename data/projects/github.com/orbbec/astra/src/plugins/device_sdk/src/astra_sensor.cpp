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
#include "astra_sensor.hpp"
#include "astra_sensor_listener.hpp"
#include "astra_sensor_stream.hpp"

namespace astra { namespace devices {

    sensor::sensor(const sensor_info& info)
        : info_(info)
    {
    }

    sensor::~sensor()
    {
    }

    device_status sensor::initialize()
    {
        if (state() == sensor_state_value::ready)
            return device_status_value::ok;

        if (state() == sensor_state_value::faulted)
            return device_status_value::device_error;

        auto rc = on_initialize();

        if (rc)
        {
            set_state(sensor_state_value::ready);

            for(auto stream : streams_)
            {
                stream->initialize();
            }
        }

        return rc;
    }

    void sensor::update()
    {
        if (state())
        {
            on_update();

            for(auto& stream : streams_)
            {
                stream->update();
            }
        }
    }

    void sensor::add_listener(const listener_ptr listener)
    {
        auto it = std::find_if(listeners_.begin(), listeners_.end(),
                               [&listener](const listener_ptr& colItem) -> bool
                               {
                                   return colItem == listener;
                               });

        if (it == listeners_.end())
            listeners_.push_back(listener);
    }

    void sensor::remove_listener(const listener_ptr listener)
    {
        auto it = std::find_if(listeners_.begin(), listeners_.end(),
                               [&listener](const listener_ptr& colItem) -> bool
                               {
                                   return colItem == listener;
                               });

        if (it != listeners_.end())
            listeners_.erase(it);
    }

    sensor::stream_const_iterator sensor::streams_begin()
    {
        return streams_.begin();
    }

    sensor::stream_const_iterator sensor::streams_end()
    {
        return streams_.end();
    }

    sensor::stream_const_iterator sensor::find_stream(sensor_stream::stream_type type)
    {
        return std::find_if(streams_.begin(),
                            streams_.end(),
                            [&type](sensor_stream::shared_ptr s) -> bool
                            {
                                return s->type() == type;
                            });
    }

    bool sensor::is_available() const
    {
        return state() == sensor_state_value::ready;
    }

    sensor_state sensor::state() const
    {
        return state_;
    }

    const sensor_info& sensor::info() const
    {
        return info_;
    }

    device_status sensor::get_property(property_id id, void* value, std::size_t size)
    {
        return on_get_property(id, value, size);
    }

    device_status sensor::get_property_size(property_id id, std::size_t* size)
    {
        return on_get_property_size(id, size);
    }

    device_status sensor::set_property(property_id id, const void* value, std::size_t size)
    {
        auto rc = on_set_property(id, value, size);

        if (rc)
        {
            raise_property_changed(id);
        }

        return rc;
    }

    void sensor::set_state(const sensor_state& state)
    {
        if (state_ == state)
            return;

        state_ = state;
        raise_state_changed();
    }

    void sensor::add_stream(stream_shared_ptr stream)
    {
        auto it = std::find(streams_.begin(), streams_.end(), stream);

        if (it == streams_.end())
        {
            stream->initialize();
            streams_.push_back(stream);
        }
    }

    void sensor::remove_stream(stream_shared_ptr stream)
    {
        auto it = std::find(streams_.begin(), streams_.end(), stream);

        if (it != streams_.end())
            streams_.erase(it);
    }

    void sensor::raise_property_changed(property_id id)
    {
        listener_collection listenersCopy = listeners_;

        for(auto listener : listenersCopy)
        {
            listener->property_changed(shared_from_this(), id);
        }
    }

    void sensor::raise_state_changed()
    {
        listener_collection listenersCopy = listeners_;

        for(auto listener : listenersCopy)
        {
            listener->state_changed(shared_from_this());
        }
    }
}}
