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
#ifndef ASTRA_STREAM_MODE_H
#define ASTRA_STREAM_MODE_H

#include <cstdint>

namespace astra { namespace devices {

    class stream_mode
    {
    public:
        stream_mode(std::uint32_t width,
                    std::uint32_t height,
                    std::uint32_t framesPerSecond,
                    std::uint32_t pixelFormatCode);

        std::uint32_t width() const;
        void set_width(std::uint32_t width);

        std::uint32_t height() const;
        void set_height(std::uint32_t height);

        std::uint32_t framesPerSecond() const;
        void set_framesPerSecond(std::uint32_t framesPerSecond);

        std::uint32_t pixelFormatCode() const;
        void set_pixelFormatCode(std::uint32_t pixelFormatCode);
    private:
        std::uint32_t width_{0};
        std::uint32_t height_{0};
        std::uint32_t framesPerSecond_{0};
        std::uint32_t pixelFormatCode_{0};
    };

    bool operator==(const stream_mode& lhs, const stream_mode& rhs);
    bool operator!=(const stream_mode& lhs, const stream_mode& rhs);

}}

#endif /* ASTRA_STREAM_MODE_H */
