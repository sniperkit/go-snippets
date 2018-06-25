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
#ifndef ONI_MAPPERS_H
#define ONI_MAPPERS_H

#include <tuple>
#include <OpenNI.h>
#include <astra/streams/Image.hpp>
#include <cstdint>

namespace orbbec { namespace ni {

    std::tuple<astra_pixel_format_t, uint8_t> convert_format(const openni::PixelFormat& oniFormat);
    std::tuple<openni::PixelFormat, uint8_t> convert_format(const astra_pixel_format_t& format);

    astra::ImageStreamMode convert_mode(const openni::VideoMode& oniMode);
    openni::VideoMode convert_mode(const astra::ImageStreamMode& mode);
}}

#endif /* ONI_MAPPERS_H */
