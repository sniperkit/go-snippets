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
#ifndef ASTRA_DEVICE_STATE_H
#define ASTRA_DEVICE_STATE_H

namespace astra { namespace devices {

    enum class device_state_value
    {
        not_ready = 0,
        ready,
        faulted
    };

    class device_state
    {
    public:
        device_state();

        device_state(const device_state_value& value)
            : value_(value) {};

        const device_state_value& value() const;
        operator bool() const;

    private:
        device_state_value value_{device_state_value::not_ready};
    };

    bool operator==(const device_state& lhs, const device_state& rhs);
    bool operator!=(const device_state& lhs, const device_state& rhs);
}}


#endif /* ASTRA_DEVICE_STATE_H */
