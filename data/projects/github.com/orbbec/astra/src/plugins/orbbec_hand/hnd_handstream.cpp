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
#include "hnd_handstream.hpp"
#include <astra/capi/streams/hand_parameters.h>

namespace astra { namespace hand {

    void handstream::on_set_parameter(astra_streamconnection_t connection,
                                      astra_parameter_id id,
                                      size_t inByteLength,
                                      astra_parameter_data_t inData)
    {
        switch (id)
        {
        case ASTRA_PARAMETER_HAND_INCLUDE_CANDIDATE_POINTS:
            set_include_candidates(inByteLength, inData);
            break;
        }
    }

    void handstream::on_get_parameter(astra_streamconnection_t connection,
                                      astra_parameter_id id,
                                      astra_parameter_bin_t& parameterBin)
    {
        switch (id)
        {
        case ASTRA_PARAMETER_HAND_INCLUDE_CANDIDATE_POINTS:
            get_include_candidates(parameterBin);
            break;
        }
    }

    void handstream::get_include_candidates(astra_parameter_bin_t& parameterBin)
    {
        size_t resultByteLength = sizeof(bool);

        astra_parameter_data_t parameterData;
        astra_status_t rc = pluginService().get_parameter_bin(resultByteLength,
                                                              &parameterBin,
                                                              &parameterData);
        if (rc == ASTRA_STATUS_SUCCESS)
        {
            memcpy(parameterData, &includeCandidatePoints_, resultByteLength);
        }
    }

    void handstream::set_include_candidates(size_t inByteLength, astra_parameter_data_t& inData)
    {
        if (inByteLength >= sizeof(bool))
        {
            bool newIncludeCandidatePoints;
            memcpy(&newIncludeCandidatePoints, inData, sizeof(bool));

            set_include_candidate_points(newIncludeCandidatePoints);
        }
    }
}}