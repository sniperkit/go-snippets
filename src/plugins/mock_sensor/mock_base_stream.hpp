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
#ifndef MOCK_BASE_STREAM_H
#define MOCK_BASE_STREAM_H

#include "astra_simple_timer.hpp"
#include "astra_sensor_stream.hpp"
#include "mock_flag_set.hpp"

#include <string>

using astra::devices::device_status;
using astra::devices::device_status_value;

namespace orbbec { namespace mocks {

    class base_stream : public astra::devices::sensor_stream
    {
    public:
        base_stream(const stream_type& type);
        virtual ~base_stream();

        virtual device_status on_get_property(property_id id, void* value, std::size_t size) override;
        virtual device_status on_get_property_size(property_id id, std::size_t* size) override;
        virtual device_status on_set_property(property_id id, const void* value, std::size_t size) override;
        virtual bool on_property_supported(property_id id) const override;

        float hFov() const;
        float vFov() const;

    protected:
        virtual device_status on_start() override;
        virtual device_status on_stop() override;

        bool has_flag(const std::string& flag);
        void set_flag(const std::string& flag, bool value);
        const flag_set& flags() const;

        void set_hFov(float hFov);
        void set_vFov(float vFov);

    private:
        virtual void on_update() override;

        flag_set flags_;
        float hFov_{0};
        float vFov_{0};

        astra::utility::simple_timer<std::chrono::high_resolution_clock> timer_;
    };
}}

#endif /* MOCK_BASE_STREAM_H */
