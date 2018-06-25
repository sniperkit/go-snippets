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
#ifndef MOCK_GRAPHICS_H
#define MOCK_GRAPHICS_H

#include <cstdint>
#include <cstddef>

namespace orbbec {

    inline void copy_rgba_to_rgb(std::uint8_t* dest, const std::uint8_t* src, std::size_t pixelCount)
    {
        for(int i = 0; i < pixelCount; i++, dest+=3, src+=4)
        {
            dest[0] = src[0];
            dest[1] = src[1];
            dest[2] = src[2];
        }
    }

    inline void copy_rgb_to_rgba(
        std::uint8_t* dest,
        const std::uint8_t* src,
        std::size_t pixelCount,
        std::uint8_t alpha = 255)
    {
        for(int i = 0; i < pixelCount; i++, dest+=4, src+=3)
        {
            dest[0] = src[0];
            dest[1] = src[1];
            dest[2] = src[2];
            dest[3] = alpha;
        }
    }
}

#endif /* MOCK_GRAPHICS_H */
