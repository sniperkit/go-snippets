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

#ifndef SIZE_HPP
#define SIZE_HPP

#include <type_traits>
#include <cassert>
#include <tuple>

namespace astra { namespace hand {

    template<typename T>
    class Size2
    {
        static_assert(std::is_arithmetic<T>::value, "size requires arithmetic elements");

    public:
        Size2()
            : width_(T(0)),
              height_(T(0))
        {}

        Size2(T width, T height)
            : width_(width),
              height_(height)
        {}

        explicit Size2(T side)
            : width_(side),
              height_(side)
        {}

        inline T width() const { return width_; }
        inline T height() const { return height_; }

        inline operator std::tuple<T, T>() const { return std::make_tuple<T, T>(width_, height_); }

        void set_width(T width) { width_ = width; }
        void set_height(T height) { height_ = height; }

        bool is_zero() const { return width_ == T(0) || height_ == T(0); }

    private:
        T width_{0};
        T height_{0};
    };

    template<typename T>
    bool operator==(const Size2<T>& lhs, const Size2<T>& rhs)
    {
        return lhs.width() == rhs.width() && lhs.height() == rhs.height();
    }

    template<typename T>
    bool operator!=(const Size2<T>& lhs, const Size2<T>& rhs)
    {
        return lhs.width() != rhs.width() || lhs.height() != rhs.height();
    }

    using Size2i = Size2<int>;
    using Size2f = Size2<float>;
}}

#endif /* SIZE_HPP */
