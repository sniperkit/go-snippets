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
#include <common/clock/Pulser.h>

namespace astra { namespace clock {

    Pulser::Pulser() :
        m_swatchName("Pulser")
    {
        m_swatch.set_mode(REAL_TIME);
        m_swatch.start(m_swatchName);

        reset();
    }

    Pulser::~Pulser()
    {

    }

    void Pulser::set_period(double period)
    {
        m_period = period;
    }

    double Pulser::get_period()
    {
        return m_period;
    }

    void Pulser::start()
    {
        m_swatch.start(m_swatchName);
    }

    void Pulser::stop()
    {
        m_swatch.stop(m_swatchName);
    }

    void Pulser::pause()
    {
        m_swatch.pause(m_swatchName);
    }

    bool Pulser::is_pulse()
    {
        bool isPulse = m_swatch.get_time_so_far(m_swatchName) >= m_period;

        if (isPulse)
        {
            reset();
        }

        return isPulse;
    }

    void Pulser::reset()
    {
        m_swatch.reset(m_swatchName);
    }
}}
