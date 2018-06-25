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
#include <astra/capi/astra_ctypes.h>
#include <astra/capi/streams/stream_types.h>
#include <string.h>
#include <astra/capi/streams/image_types.h>
#include <astra/capi/streams/image_capi.h>
#include <astra/capi/streams/image_parameters.h>

ASTRA_BEGIN_DECLS

ASTRA_API_EX astra_status_t astra_imagestream_get_mirroring(astra_imagestream_t imageStream,
                                                            bool* mirroring)
{
    return astra_stream_get_parameter_fixed(imageStream,
                                            ASTRA_PARAMETER_IMAGE_MIRRORING,
                                            sizeof(bool),
                                            reinterpret_cast<astra_parameter_data_t*>(mirroring));
}

ASTRA_API_EX astra_status_t astra_imagestream_set_mirroring(astra_imagestream_t imageStream,
                                                            bool mirroring)
{
    return astra_stream_set_parameter(imageStream,
                                      ASTRA_PARAMETER_IMAGE_MIRRORING,
                                      sizeof(bool),
                                      reinterpret_cast<astra_parameter_data_t>(&mirroring));
}

ASTRA_API_EX astra_status_t astra_imagestream_get_hfov(astra_imagestream_t imageStream,
                                                       float* hFov)
{
    return astra_stream_get_parameter_fixed(imageStream,
                                            ASTRA_PARAMETER_IMAGE_HFOV,
                                            sizeof(float),
                                            reinterpret_cast<astra_parameter_data_t*>(hFov));
}

ASTRA_API_EX astra_status_t astra_imagestream_get_vfov(astra_imagestream_t imageStream,
                                                       float* vFov)
{
    return astra_stream_get_parameter_fixed(imageStream,
                                            ASTRA_PARAMETER_IMAGE_VFOV,
                                            sizeof(float),
                                            reinterpret_cast<astra_parameter_data_t*>(vFov));
}

ASTRA_API_EX astra_status_t astra_imagestream_request_modes(astra_imagestream_t imageStream,
                                                            astra_result_token_t* token,
                                                            size_t* count)
{
    return astra_generic_stream_request_array<astra_imagestream_mode_t>(imageStream,
                                                                        ASTRA_PARAMETER_IMAGE_MODES,
                                                                        token,
                                                                        count);
}

ASTRA_API_EX astra_status_t astra_imagestream_get_modes_result(astra_imagestream_t imageStream,
                                                               astra_result_token_t token,
                                                               astra_imagestream_mode_t* modes,
                                                               size_t modeCount)
{

    return astra_generic_stream_get_result_array<astra_imagestream_mode_t>(imageStream,
                                                                           token,
                                                                           modes,
                                                                           modeCount);
}

ASTRA_API_EX astra_status_t astra_imagestream_set_mode(astra_imagestream_t imageStream,
                                                       const astra_imagestream_mode_t* mode)
{
    return astra_stream_set_parameter(imageStream,
                                      ASTRA_PARAMETER_IMAGE_MODE,
                                      sizeof(astra_imagestream_mode_t),
                                      const_cast<astra_imagestream_mode_t*>(mode));
}

ASTRA_API_EX astra_status_t astra_reader_get_imageframe(astra_reader_frame_t readerFrame,
                                                        astra_stream_type_t type,
                                                        astra_stream_subtype_t subtype,
                                                        astra_imageframe_t* imageFrame)
{
    return astra_generic_frame_get<astra_imageframe_wrapper_t>(readerFrame,
                                                               type,
                                                               subtype,
                                                               imageFrame);
}

ASTRA_API_EX astra_status_t astra_imageframe_get_frameindex(astra_imageframe_t imageFrame,
                                                            astra_frame_index_t* index)
{
    return astra_generic_frame_get_frameindex(imageFrame, index);
}

ASTRA_API_EX astra_status_t astra_imageframe_get_data_byte_length(astra_imageframe_t imageFrame,
                                                                  size_t* length)
{
    astra_image_metadata_t metadata = imageFrame->metadata;

    uint8_t bpp;
    astra_pixelformat_get_bytes_per_pixel(metadata.pixelFormat, &bpp);
    size_t size = metadata.width * metadata.height * bpp;
    *length = size;

    return ASTRA_STATUS_SUCCESS;
}

ASTRA_API_EX astra_status_t astra_imageframe_get_data_ptr(astra_imageframe_t imageFrame,
                                                          void** data,
                                                          size_t* byteLength)
{
    *data = imageFrame->data;
    astra_imageframe_get_data_byte_length(imageFrame, byteLength);

    return ASTRA_STATUS_SUCCESS;
}

ASTRA_API_EX astra_status_t astra_imageframe_copy_data(astra_imageframe_t imageFrame,
                                                       void* data)
{
    size_t byteLength;
    astra_imageframe_get_data_byte_length(imageFrame, &byteLength);

    memcpy(data, imageFrame->data, byteLength);

    return ASTRA_STATUS_SUCCESS;
}

ASTRA_API_EX astra_status_t astra_imageframe_get_metadata(astra_imageframe_t imageFrame,
                                                          astra_image_metadata_t* metadata)
{
    *metadata = imageFrame->metadata;
    return ASTRA_STATUS_SUCCESS;
}

ASTRA_END_DECLS
