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
#include "hnd_debug_handstream.hpp"
#include <astra/capi/streams/hand_parameters.h>
#include <cstring>

namespace astra { namespace hand {

    void debug_handstream::on_set_parameter(astra_streamconnection_t connection,
                                            astra_parameter_id id,
                                            std::size_t inByteLength,
                                            astra_parameter_data_t inData)
    {
        switch (id)
        {
        case ASTRA_PARAMETER_DEBUG_HAND_VIEW_TYPE:
            set_view_parameter(inByteLength, inData);
            break;
        case ASTRA_PARAMETER_DEBUG_HAND_USE_MOUSE_PROBE:
            set_use_mouse_probe(inByteLength, inData);
            break;
        case ASTRA_PARAMETER_DEBUG_HAND_MOUSE_NORM_POSITION:
            set_mouse_norm_position(inByteLength, inData);
            break;
        case ASTRA_PARAMETER_DEBUG_HAND_PAUSE_INPUT:
            set_pause_input(inByteLength, inData);
            break;
        case ASTRA_PARAMETER_DEBUG_HAND_LOCK_SPAWN_POINT:
            set_lock_spawn_point(inByteLength, inData);
            break;
        }
    }

    void debug_handstream::on_get_parameter(astra_streamconnection_t connection,
                                            astra_parameter_id id,
                                            astra_parameter_bin_t& parameterBin)
    {
        switch (id)
        {
        case ASTRA_PARAMETER_DEBUG_HAND_VIEW_TYPE:
            get_view_parameter(parameterBin);
            break;
        }
    }

    void debug_handstream::on_invoke(astra_streamconnection_t connection,
                                     astra_command_id commandId,
                                     std::size_t inByteLength,
                                     astra_parameter_data_t inData,
                                     astra_parameter_bin_t& parameterBin)
    {
    }

    void debug_handstream::get_view_parameter(astra_parameter_bin_t& parameterBin)
    {
        std::size_t resultByteLength = sizeof(debug_handview_type);

        astra_parameter_data_t parameterData;
        astra_status_t rc = pluginService().get_parameter_bin(resultByteLength,
                                                              &parameterBin,
                                                              &parameterData);
        if (rc == ASTRA_STATUS_SUCCESS)
        {
            std::memcpy(parameterData, &viewType_, resultByteLength);
        }
    }

    void debug_handstream::set_view_parameter(std::size_t inByteLength, astra_parameter_data_t& inData)
    {
        if (inByteLength >= sizeof(debug_handview_type))
        {
            debug_handview_type newViewType;
            std::memcpy(&newViewType, inData, sizeof(debug_handview_type));

            set_view_type(newViewType);
        }
    }

    void debug_handstream::set_use_mouse_probe(std::size_t inByteLength, astra_parameter_data_t& inData)
    {
        if (inByteLength >= sizeof(bool))
        {
            bool newUseMouseProbe;
            std::memcpy(&newUseMouseProbe, inData, sizeof(bool));

            useMouseProbe_ = newUseMouseProbe;
        }
    }

    void debug_handstream::set_mouse_norm_position(std::size_t inByteLength, astra_parameter_data_t& inData)
    {
        if (inByteLength >= sizeof(astra_vector2f_t))
        {
            astra_vector2f_t newMousePosition;
            std::memcpy(&newMousePosition, inData, sizeof(astra_vector2f_t));

            mouseNormPosition_ = newMousePosition;
        }
    }

    void debug_handstream::set_pause_input(std::size_t inByteLength, astra_parameter_data_t& inData)
    {
        if (inByteLength >= sizeof(bool))
        {
            bool newPauseInput;
            std::memcpy(&newPauseInput, inData, sizeof(bool));

            pauseInput_ = newPauseInput;
        }
    }

    void debug_handstream::set_lock_spawn_point(std::size_t inByteLength, astra_parameter_data_t& inData)
    {
        if (inByteLength >= sizeof(bool))
        {
            bool newLockSpawnPoint;
            std::memcpy(&newLockSpawnPoint, inData, sizeof(bool));

            if (newLockSpawnPoint && !lockSpawnPoint_)
            {
                spawnNormPosition_ = mouseNormPosition_;
            }
            lockSpawnPoint_ = newLockSpawnPoint;
        }
    }
}}
