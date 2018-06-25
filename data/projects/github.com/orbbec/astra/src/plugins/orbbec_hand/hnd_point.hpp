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

#ifndef POINT_HPP
#define POINT_HPP

#include <type_traits>
#include <cmath>

namespace astra { namespace hand {

    template<typename T>
    class Point2
    {
        static_assert(std::is_arithmetic<T>::value, "Point2 requires arithmetic elements");

    public:
        T x{0};
        T y{0};

        Point2(const T& x, const T& y)
            : x(x),
              y(y)
        { }

        Point2()
            : x(static_cast<T>(0)),
              y(static_cast<T>(0))
        { }

        const float length() const { return std::sqrt(x*x + y*y); }

        void operator+=(const Point2<T>& rhs)
        {
            this->x = this->x + rhs.x;
            this->y = this->y + rhs.y;
        }

        template<typename S>
        void operator/=(const S& scalar)
        {
            this->x = this->x / scalar;
            this->y = this->y / scalar;
        }
    };


    template<typename T>
    bool operator==(const Point2<T>& left, const Point2<T>& right)
    {
        return left.x == right.x && left.y == right.y;
    }

    template<typename T>
    bool operator!=(const Point2<T>& left, const Point2<T>& right)
    {
        return !(left == right);
    }

    template<typename T>
    Point2<T> operator-(const Point2<T>& left, const Point2<T>& right)
    {
        return Point2<T>(left.x - right.x, left.y - right.y);
    }

    template<typename T>
    Point2<T> operator+(const Point2<T>& left, const Point2<T>& right)
    {
        return Point2<T>(left.x + right.x, left.y + right.y);
    }

    template<typename T>
    Point2<T> operator-(const Point2<T>& p)
    {
        return Point2<T>(-p.x, -p.y);
    }

    template<typename T, typename S>
    Point2<T> operator*(const Point2<T>& left, const S& right)
    {
        return Point2<T>(left.x * right, left.y * right);
    }

    template<typename T, typename S>
    Point2<T> operator/(const Point2<T>& left, const S& right)
    {
        return Point2<T>(left.x / right, left.y / right);
    }

    using Point2i = Point2<int>;
    using Point2f = Point2<float>;

    template<typename T>
    class Vector3
    {
        static_assert(std::is_arithmetic<T>::value, "Vector3 requires arithmetic elements");

    public:
        T x{0};
        T y{0};
        T z{0};


        Vector3()
            : x(static_cast<T>(0)),
              y(static_cast<T>(0)),
              z(static_cast<T>(0))
        { }

        Vector3(const T& x, const T& y, const T& z)
            : x(x),
              y(y),
              z(z)
        { }

        Vector3(const Point2<T>& p2, const T& z)
            : x(p2.x),
              y(p2.y),
              z(z)
        { }

        inline float length() const
        {
            return std::sqrt(x * x + y * y + z * z);
        }

        inline float length_squared() const
        {
            return x * x + y * y + z * z;
        }

        inline float dot(const Vector3& v) const
        {
            return x * v.x + y * v.y + z * v.z;
        }

        inline Vector3 cross(const Vector3& v) const
        {
            return Vector3(y*v.z - z*v.y, z*v.x - x*v.z, x*v.y - y*v.x);
        }

        static inline Vector3 zero()
        {
            static Vector3 zero;
            return zero;
        }

        static inline Vector3 normalize(Vector3 v)
        {
            double length = std::sqrt(v.x*v.x + v.y*v.y + v.z*v.z);
            if (length < 1e-9)
            {
                return Vector3(0.0f, 0.0f, 0.0f);
            }
            else
            {
                return Vector3(
                    static_cast<float>(v.x / length),
                    static_cast<float>(v.y / length),
                    static_cast<float>(v.z / length));
            }
        }

        inline bool is_zero() const
        {
            return *this == zero();
        }

        inline Vector3& operator+=(const Vector3& rhs)
        {
            this->x = this->x + rhs.x;
            this->y = this->y + rhs.y;
            this->z = this->z + rhs.z;
            return *this;
        }

        inline Vector3& operator-=(const Vector3& rhs)
        {
            this->x = this->x - rhs.x;
            this->y = this->y - rhs.y;
            this->z = this->z - rhs.z;
            return *this;
        }

        inline Vector3& operator*=(const float& rhs)
        {
            this->x = this->x * rhs;
            this->y = this->y * rhs;
            this->z = this->z * rhs;
            return *this;
        }

        inline Vector3& operator*=(const Vector3& rhs)
        {
            this->x = this->x * rhs.x;
            this->y = this->y * rhs.y;
            this->z = this->z * rhs.z;
            return *this;
        }

        inline Vector3& operator/=(const float& rhs)
        {
            this->x = this->x / rhs;
            this->y = this->y / rhs;
            this->z = this->z / rhs;
            return *this;
        }

        inline Vector3 operator-()
        {
            return Vector3(-this->x, -this->y, -this->z);
        }
    };

    template<typename T>
    inline bool operator==(const Vector3<T>& lhs, const Vector3<T>& rhs)
    {
        return lhs.x == rhs.x && lhs.y == rhs.y && lhs.z == rhs.z;
    }

    template<typename T>
    inline bool operator!=(const Vector3<T>& lhs, const Vector3<T>& rhs)
    {
        return !(lhs == rhs);
    }

    template<typename T>
    inline Vector3<T> operator+(const Vector3<T>& lhs, const Vector3<T>& rhs)
    {
        return Vector3<T>(lhs.x + rhs.x, lhs.y + rhs.y, lhs.z + rhs.z);
    }

    template<typename T>
    inline Vector3<T> operator-(const Vector3<T>& lhs, const Vector3<T>& rhs)
    {
        return Vector3<T>(lhs.x - rhs.x, lhs.y - rhs.y, lhs.z - rhs.z);
    }

    template<typename T, typename S>
    inline Vector3<T> operator*(const Vector3<T>& lhs, const S& rhs)
    {
        return Vector3<T>(lhs.x * rhs, lhs.y * rhs, lhs.z * rhs);
    }

    template<typename T>
    inline Vector3<T> operator*(const Vector3<T>& lhs, const Vector3<T>& rhs)
    {
        return Vector3<T>(lhs.x * rhs.x, lhs.y * rhs.y, lhs.z * rhs.z);
    }

    template<typename T, typename S>
    inline Vector3<T> operator*(const S& lhs, const Vector3<T>& rhs)
    {
        return rhs * lhs;
    }

    template<typename T, typename S>
    inline Vector3<T> operator/(const Vector3<T>& lhs, const S& rhs)
    {
        return Vector3<T>(lhs.x / rhs, lhs.y / rhs, lhs.z / rhs);
    }

    using Vector3i = Vector3<int>;
    using Vector3f = Vector3<float>;
}}

#endif /* POINT_HPP */
