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
#include <astra_core/capi/astra_types.h>
#include "astra_generic_stream_api.hpp"
#include <memory.h>
#include <astra/capi/astra_ctypes.h>
#include <astra/capi/streams/stream_types.h>
#include <astra/capi/streams/hand_capi.h>
#include <string.h>
#include <astra/capi/streams/image_capi.h>
#include <astra/capi/streams/hand_parameters.h>

ASTRA_BEGIN_DECLS

ASTRA_API_EX astra_status_t astra_reader_get_handstream(astra_reader_t reader,
                                                           astra_handstream_t* handStream)

{
    return astra_reader_get_stream(reader,
                                      ASTRA_STREAM_HAND,
                                      DEFAULT_SUBTYPE,
                                      handStream);
}

ASTRA_API_EX astra_status_t astra_frame_get_handframe(astra_reader_frame_t readerFrame,
                                                          astra_handframe_t* handFrame)
{
    return astra_generic_frame_get<astra_handframe_wrapper_t>(readerFrame,
                                                                    ASTRA_STREAM_HAND,
                                                                    DEFAULT_SUBTYPE,
                                                                    handFrame);
}

ASTRA_API_EX astra_status_t astra_frame_get_handframe_with_subtype(astra_reader_frame_t readerFrame,
                                                                            astra_stream_subtype_t subtype,
                                                                            astra_handframe_t* handFrame)
{
    return astra_generic_frame_get<astra_handframe_wrapper_t>(readerFrame,
                                                                    ASTRA_STREAM_HAND,
                                                                    subtype,
                                                                    handFrame);
}

ASTRA_API_EX astra_status_t astra_handframe_get_frameindex(astra_handframe_t handFrame,
                                                                     astra_frame_index_t* index)
{
    return astra_generic_frame_get_frameindex(handFrame, index);
}

ASTRA_API_EX astra_status_t astra_handframe_get_hand_count(astra_handframe_t handFrame,
                                                                   size_t* handCount)
{
    *handCount = handFrame->handCount;

    return ASTRA_STATUS_SUCCESS;
}

ASTRA_API_EX astra_status_t astra_handframe_copy_hands(astra_handframe_t handFrame,
                                                                astra_handpoint_t* handPointsDestination)
{
    size_t size = handFrame->handCount * sizeof(astra_handpoint_t);

    memcpy(handPointsDestination, handFrame->handpoints, size);

    return ASTRA_STATUS_SUCCESS;
}

ASTRA_API_EX astra_status_t astra_handframe_get_shared_hand_array(astra_handframe_t handFrame,
                                                                           astra_handpoint_t** handpoints,
                                                                           size_t* handCount)
{
    *handpoints = handFrame->handpoints;
    astra_handframe_get_hand_count(handFrame, handCount);

    return ASTRA_STATUS_SUCCESS;
}

ASTRA_API_EX astra_status_t astra_handstream_get_include_candidate_points(astra_handstream_t handStream,
                                                                                   bool* includeCandidatePoints)
{
    return astra_stream_get_parameter_fixed(handStream,
                                               ASTRA_PARAMETER_HAND_INCLUDE_CANDIDATE_POINTS,
                                               sizeof(bool),
                                               reinterpret_cast<astra_parameter_data_t*>(includeCandidatePoints));

}

ASTRA_API_EX astra_status_t astra_handstream_set_include_candidate_points(astra_handstream_t handStream,
                                                                                   bool includeCandidatePoints)
{
    return astra_stream_set_parameter(handStream,
                                         ASTRA_PARAMETER_HAND_INCLUDE_CANDIDATE_POINTS,
                                         sizeof(bool),
                                         reinterpret_cast<astra_parameter_data_t>(&includeCandidatePoints));

}

ASTRA_API_EX astra_status_t astra_reader_get_debug_handstream(astra_reader_t reader,
                                                                       astra_debug_handstream_t* debugHandStream)

{
    return astra_reader_get_stream(reader,
                                      ASTRA_STREAM_DEBUG_HAND,
                                      DEFAULT_SUBTYPE,
                                      debugHandStream);
}

ASTRA_API_EX astra_status_t astra_frame_get_debug_handframe(astra_reader_frame_t readerFrame,
                                                                     astra_debug_handframe_t* debugHandFrame)
{
    return astra_reader_get_imageframe(readerFrame,
                                          ASTRA_STREAM_DEBUG_HAND,
                                          DEFAULT_SUBTYPE,
                                          debugHandFrame);
}

ASTRA_API_EX astra_status_t astra_debug_handstream_get_view_type(astra_debug_handstream_t debugHandStream,
                                                                          astra_debug_hand_view_type_t* viewType)
{
    return astra_stream_get_parameter_fixed(debugHandStream,
                                               ASTRA_PARAMETER_DEBUG_HAND_VIEW_TYPE,
                                               sizeof(astra_debug_hand_view_type_t),
                                               reinterpret_cast<astra_parameter_data_t*>(viewType));
}

ASTRA_API_EX astra_status_t astra_debug_handstream_set_view_type(astra_debug_handstream_t debugHandStream,
                                                                          astra_debug_hand_view_type_t viewType)
{
    return astra_stream_set_parameter(debugHandStream,
                                         ASTRA_PARAMETER_DEBUG_HAND_VIEW_TYPE,
                                         sizeof(astra_debug_hand_view_type_t),
                                         reinterpret_cast<astra_parameter_data_t>(&viewType));
}

ASTRA_API_EX astra_status_t astra_debug_handstream_set_mouse_position(astra_debug_handstream_t debugHandStream,
                                                                          astra_vector2f_t normPosition)
{
    return astra_stream_set_parameter(debugHandStream,
                                         ASTRA_PARAMETER_DEBUG_HAND_MOUSE_NORM_POSITION,
                                         sizeof(astra_vector2f_t),
                                         reinterpret_cast<astra_parameter_data_t>(&normPosition));
}

ASTRA_API_EX astra_status_t astra_debug_handstream_set_use_mouse_probe(astra_debug_handstream_t debugHandStream,
                                                                          bool useMouseProbe)
{
    return astra_stream_set_parameter(debugHandStream,
                                         ASTRA_PARAMETER_DEBUG_HAND_USE_MOUSE_PROBE,
                                         sizeof(bool),
                                         reinterpret_cast<astra_parameter_data_t>(&useMouseProbe));
}

ASTRA_API_EX astra_status_t astra_debug_handstream_set_pause_input(astra_debug_handstream_t debugHandStream,
                                                                            bool pauseInput)
{
    return astra_stream_set_parameter(debugHandStream,
                                         ASTRA_PARAMETER_DEBUG_HAND_PAUSE_INPUT,
                                         sizeof(bool),
                                         reinterpret_cast<astra_parameter_data_t>(&pauseInput));
}

ASTRA_API_EX astra_status_t astra_debug_handstream_set_lock_spawn_point(astra_debug_handstream_t debugHandStream,
                                                                                 bool lockSpawnPoint)
{
    return astra_stream_set_parameter(debugHandStream,
        ASTRA_PARAMETER_DEBUG_HAND_LOCK_SPAWN_POINT,
                                         sizeof(bool),
                                         reinterpret_cast<astra_parameter_data_t>(&lockSpawnPoint));
}

ASTRA_END_DECLS
