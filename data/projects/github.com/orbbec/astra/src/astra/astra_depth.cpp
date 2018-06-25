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
#include <astra/capi/streams/depth_types.h>
#include <memory.h>
#include <astra/capi/streams/stream_types.h>
#include <astra/capi/streams/depth_parameters.h>
#include <string.h>
#include <astra/capi/streams/image_capi.h>
#include <astra/capi/streams/image_parameters.h>
#include <unordered_map>
#include <Shiny.h>

using ConversionMap = std::unordered_map<astra_depthstream_t, conversion_cache_t>;

ConversionMap g_astra_conversion_map;

conversion_cache_t astra_depth_fetch_conversion_cache(astra_depthstream_t depthStream)
{
    PROFILE_FUNC();
    auto it = g_astra_conversion_map.find(depthStream);

    if (it != g_astra_conversion_map.end())
    {
        return it->second;
    }
    else
    {
        PROFILE_BEGIN(depth_cache_get);
        conversion_cache_t conversionCache;
        astra_stream_get_parameter_fixed(depthStream,
                                         ASTRA_PARAMETER_DEPTH_CONVERSION_CACHE,
                                         sizeof(conversion_cache_t),
                                         reinterpret_cast<astra_parameter_data_t*>(&conversionCache));
        g_astra_conversion_map.insert(std::make_pair(depthStream, conversionCache));
        PROFILE_END();
        return conversionCache;
    }
}

ASTRA_BEGIN_DECLS

ASTRA_API_EX astra_status_t astra_convert_depth_to_world(astra_depthstream_t depthStream,
                                                         float depthX, float depthY, float depthZ,
                                                         float* pWorldX, float* pWorldY, float* pWorldZ)
{
    PROFILE_FUNC();
    conversion_cache_t conversionCache = astra_depth_fetch_conversion_cache(depthStream);

    PROFILE_BEGIN(depth_to_world_math);
    float normalizedX = depthX / conversionCache.resolutionX - .5f;
    float normalizedY = .5f - depthY / conversionCache.resolutionY;

    *pWorldX = normalizedX * depthZ * conversionCache.xzFactor;
    *pWorldY = normalizedY * depthZ * conversionCache.yzFactor;
    *pWorldZ = depthZ;

    PROFILE_END();

    return ASTRA_STATUS_SUCCESS;
}

ASTRA_API_EX astra_status_t astra_convert_world_to_depth(astra_depthstream_t depthStream,
                                                         float worldX, float worldY, float worldZ,
                                                         float* pDepthX, float* pDepthY, float* pDepthZ)
{
    PROFILE_FUNC();
    conversion_cache_t conversionCache = astra_depth_fetch_conversion_cache(depthStream);

    *pDepthX = conversionCache.coeffX * worldX / worldZ + conversionCache.halfResX;
    *pDepthY = conversionCache.halfResY - conversionCache.coeffY * worldY / worldZ;
    *pDepthZ = worldZ;

    return ASTRA_STATUS_SUCCESS;
}


ASTRA_API_EX astra_status_t astra_reader_get_depthstream(astra_reader_t reader,
                                                         astra_depthstream_t* depthStream)

{
    return astra_reader_get_stream(reader,
                                   ASTRA_STREAM_DEPTH,
                                   DEFAULT_SUBTYPE,
                                   depthStream);
}

ASTRA_API_EX astra_status_t astra_depthstream_get_depth_to_world_data(astra_depthstream_t depthStream,
                                                                      conversion_cache_t* conversion_data)
{
    conversion_cache_t conversionCache = astra_depth_fetch_conversion_cache(depthStream);
    memcpy(conversion_data, &conversionCache, sizeof(conversion_cache_t));

    return ASTRA_STATUS_SUCCESS;
}

ASTRA_API_EX astra_status_t astra_depthstream_get_hfov(astra_depthstream_t depthStream,
                                                       float* hFov)
{
    return astra_stream_get_parameter_fixed(depthStream,
                                            ASTRA_PARAMETER_IMAGE_HFOV,
                                            sizeof(float),
                                            reinterpret_cast<astra_parameter_data_t*>(hFov));
}

ASTRA_API_EX astra_status_t astra_depthstream_get_vfov(astra_depthstream_t depthStream,
                                                       float* vFov)
{
    return astra_stream_get_parameter_fixed(depthStream,
                                            ASTRA_PARAMETER_IMAGE_VFOV,
                                            sizeof(float),
                                            reinterpret_cast<astra_parameter_data_t*>(vFov));
}

ASTRA_API_EX astra_status_t astra_depthstream_set_registration(astra_depthstream_t depthStream,
                                                               bool enabled)
{
    return astra_stream_set_parameter(depthStream, ASTRA_PARAMETER_DEPTH_REGISTRATION, sizeof(bool), &enabled);
}

ASTRA_API_EX astra_status_t astra_depthstream_get_registration(astra_depthstream_t depthStream,
                                                               bool* enabled)
{
    return astra_stream_get_parameter_fixed(depthStream,
                                            ASTRA_PARAMETER_DEPTH_REGISTRATION,
                                            sizeof(bool),
                                            reinterpret_cast<astra_parameter_data_t*>(enabled));
}

ASTRA_API_EX astra_status_t astra_frame_get_depthframe(astra_reader_frame_t readerFrame,
                                                       astra_depthframe_t* depthFrame)
{
    return astra_reader_get_imageframe(readerFrame,
                                       ASTRA_STREAM_DEPTH,
                                       DEFAULT_SUBTYPE,
                                       depthFrame);
}

ASTRA_API_EX astra_status_t astra_frame_get_depthframe_with_subtype(astra_reader_frame_t readerFrame,
                                                                    astra_stream_subtype_t subtype,
                                                                    astra_depthframe_t* depthFrame)
{
    return astra_reader_get_imageframe(readerFrame,
                                       ASTRA_STREAM_DEPTH,
                                       subtype,
                                       depthFrame);
}

ASTRA_API_EX astra_status_t astra_depthframe_get_frameindex(astra_depthframe_t depthFrame,
                                                            astra_frame_index_t* index)
{
    return astra_generic_frame_get_frameindex(depthFrame, index);
}

ASTRA_API_EX astra_status_t astra_depthframe_get_data_byte_length(astra_depthframe_t depthFrame,
                                                                  size_t* byteLength)
{
    return astra_imageframe_get_data_byte_length(depthFrame, byteLength);
}

ASTRA_API_EX astra_status_t astra_depthframe_get_data_ptr(astra_depthframe_t depthFrame,
                                                          int16_t** data,
                                                          size_t* byteLength)
{
    void* voidData = nullptr;
    astra_imageframe_get_data_ptr(depthFrame, &voidData, byteLength);
    *data = static_cast<int16_t*>(voidData);

    return ASTRA_STATUS_SUCCESS;
}

ASTRA_API_EX astra_status_t astra_depthframe_copy_data(astra_depthframe_t depthFrame,
                                                       int16_t* data)
{
    return astra_imageframe_copy_data(depthFrame, data);
}

ASTRA_API_EX astra_status_t astra_depthframe_get_metadata(astra_depthframe_t depthFrame,
                                                          astra_image_metadata_t* metadata)
{
    return astra_imageframe_get_metadata(depthFrame, metadata);
}

ASTRA_END_DECLS
