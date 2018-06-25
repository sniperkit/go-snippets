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

#ifndef PIXEL_HPP
#define PIXEL_HPP

#include <cstdint>

namespace astra { namespace hand {

    template<typename TComponent, unsigned Components>
    struct Pixel
    {
        using ComponentType = TComponent;
        static const unsigned kNumOfComponents = Components;

        TComponent values[Components];
    };

    template<typename T>
    struct Pixel<T, 3>
    {
        using ComponentType = T;
        static const unsigned kNumOfComponents = 3;

        inline Pixel()
            : r(T()), g(T()), b(T())
        {}

        inline Pixel(T r, T g, T b)
            : r(r), g(g), b(b)
        {}

        union {
            struct {
                T r;
                T g;
                T b;
            };
            struct {
                T x;
                T y;
                T z;
            };
            T values[kNumOfComponents];
        };
    };

    template<typename T>
    struct Pixel<T, 4>
    {
        using ComponentType = T;
        static const unsigned kNumOfComponents = 4;

        inline Pixel()
            : r(T()), g(T()), b(T(), a(T()))
        {}

        inline Pixel(T r, T g, T b, T a)
            : r(r), g(g), b(b), a(a)
        {}

        union {
            struct {
                T r;
                T g;
                T b;
                T a;
            };
            struct {
                T x;
                T y;
                T z;
                T w;
            };
            T values[kNumOfComponents];
        };
    };

    //! R8G8B8 pixel (24-bits)
    using RGBPixel = Pixel<std::uint8_t, 3>;
    //! R32G32B32 pixel (96-bits)
    using RGBPixelF = Pixel<float, 3>;
    //! R8G8B8A8 pixel with alpha channel (32-bits)
    using RGBAPixel = Pixel<std::uint8_t, 4>;
    //! R32G32B32A32 pixel with alpha channel (128-bits)
    using RGBAPixelF = Pixel<float, 4>;
}}

#endif /* PIXEL_HPP */
