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

#ifndef BITMAP_HPP
#define BITMAP_HPP

#include <algorithm>
#include <memory>
#include <cstdint>
#include <limits>
#include "hnd_pixel.hpp"
#include "hnd_point.hpp"
#include "hnd_size.hpp"

namespace astra { namespace hand {

    //! Encapsulates a 2D matrix of "pixel" data
    template<typename T>
    class Bitmap
    {
    public:
        using ValueType = T;
        using SizeType = std::size_t;

        Bitmap()
            : def_(std::make_shared<BitmapDef>(0, 0, sizeof(T), 0))
        {}

        Bitmap(void* data, unsigned width, unsigned height)
            : Bitmap(data, width, height, sizeof(T))
        {}

        Bitmap(void* data, Size2i size)
            : Bitmap(data, size.width(), size.height(), sizeof(T))
        {}

        Bitmap(unsigned width, unsigned height)
            : Bitmap(width, height, sizeof(T))
        {}

        Bitmap(Size2i size)
            : Bitmap(size.width(), size.height(), sizeof(T))
        {}

        Bitmap(Bitmap&& other)
            : def_(std::move(other.def_))
        {}

        Bitmap& operator=(Bitmap&& other)
        {
            def_ = std::move(other.def_);
            return *this;
        }

        Bitmap(Bitmap& other) = default;
        Bitmap& operator=(Bitmap& other) = default;

        inline unsigned width() const { return def_->width; }
        inline unsigned height() const { return def_->height; }
        inline Size2i size() const { return Size2i(def_->width, def_->height); }
        inline unsigned bpp() const { return def_->bpp; }
        inline unsigned length() const { return def_->width * def_->height; }
        inline unsigned byte_length() const { return def_->byteLength; }
        inline void* raw_ptr() { return def_->data; }

        inline const std::uint8_t* bytes() const { return def_->data; }

        inline const T* data() const { return reinterpret_cast<const T*>(def_->data); }
        inline const T* data(int row) const { return reinterpret_cast<const T*>(def_->data) + row * width(); }

        inline T* data() { return reinterpret_cast<T*>(def_->data); }
        inline T* data(int row) { return reinterpret_cast<T*>(def_->data) + row * width(); }

        inline Bitmap<T> submap(int left, int top, int width, int height)
        {
            Bitmap<T> sub(width, height);
            const auto* srcPtr = this->data();
            auto* destPtr = sub.data();

            for(int y = top; y < y + height && y < this->height(); ++y, srcPtr+= this->width())
            {
                for(int x = left; x < x + width && x < this->width(); ++x, ++destPtr)
                {
                    *destPtr = *(srcPtr + x);
                }
            }

            return sub;
        }

        inline Bitmap<T> clone() const
        {
            Bitmap<T> copy(size());
            copy_to(*this, copy);
            return copy;
        }

        void fill(const T& value)
        {
            auto* ptr = data();
            const auto* end = ptr + length();

            for(; ptr != end; ++ptr)
                *ptr = value;
        }

        template<typename U>
        void fill(const T& value, const Bitmap<U>& mask)
        {
            assert(this->size() == mask.size());

            auto* ptr = data();
            const auto* end = ptr + length();
            const auto* maskPtr = mask.data();

            for(; ptr != end; ++ptr, ++maskPtr)
            {
                if (*maskPtr == 0)
                    continue;

                *ptr = value;
            }
        }

        void recreate(unsigned width, unsigned height)
        {
            if (width == this->width() && height == this->height())
                return;

            assert(width < 10000);
            assert(height < 10000);

            *def_ = BitmapDef(width, height, this->bpp(), width*height*this->bpp());
            def_->reallocate();
        }

        void recreate(Size2i newSize)
        {
            recreate(newSize.width(), newSize.height());
        }

        const T& operator[](SizeType index) const
        {
            return reinterpret_cast<const T&>(bytes()[index * bpp()]);
        }

        T& operator[](SizeType index)
        {
            return reinterpret_cast<T&>(bytes()[index * bpp()]);
        }

        const T& at(const unsigned x, const unsigned y) const
        {
            const auto index = (width() * y + x) * bpp();
            return reinterpret_cast<const T&>(bytes()[index]);
        }

        const T& at(const Point2i p) const
        {
            const auto index = (width() * p.y + p.x) * bpp();
            return reinterpret_cast<const T&>(bytes()[index]);
        }

        T& at(const unsigned x, const unsigned y)
        {
            const auto index = (width() * y + x) * bpp();
            return reinterpret_cast<T&>(bytes()[index]);
        }

        T& at(const Point2i p)
        {
            const auto index = (width() * p.y + p.x) * bpp();
            return reinterpret_cast<T&>(bytes()[index]);
        }

    private:
        Bitmap(void* data, unsigned width, unsigned height, unsigned bpp)
            : def_(std::make_shared<BitmapDef>(width, height, bpp, width*height*bpp))
        {
            assert(width < 10000);
            assert(height < 10000);

            def_->data = static_cast<std::uint8_t*>(data);
        }

        Bitmap(unsigned width, unsigned height, unsigned bpp)
            : def_(std::make_shared<BitmapDef>(width, height, bpp, width*height*bpp))
        {
            assert(width < 10000);
            assert(height < 10000);

            def_->reallocate();
        }

        inline std::uint8_t* bytes() { return def_->data; }

        struct BitmapDef
        {
            unsigned width{0};
            unsigned height{0};
            unsigned bpp{0};
            unsigned byteLength{0};
            std::uint8_t* data{nullptr};
            std::unique_ptr<std::uint8_t[]> ownedData{nullptr};

            BitmapDef(unsigned width, unsigned height, unsigned bpp, unsigned byteLength)
                : width(width), height(height), bpp(bpp), byteLength(byteLength)
            {}

            // explicitly defined for VS2013 support
            BitmapDef(BitmapDef&& other)
                : width(other.width),
                  height(other.height),
                  bpp(other.bpp),
                  byteLength(other.byteLength),
                  data(other.data),
                  ownedData(std::move(other.ownedData))
            {}

            BitmapDef& operator=(BitmapDef&& rhs)
            {
                width = rhs.width;
                height = rhs.height;
                bpp = rhs.bpp;
                byteLength = rhs.byteLength;
                data = rhs.data;
                ownedData = std::move(rhs.ownedData);
                return *this;
            }

            void reallocate()
            {
                if (byteLength > 0)
                {
                    ownedData = std::unique_ptr<std::uint8_t[]>(new std::uint8_t[byteLength]);
                    data = ownedData.get();
                }
                else
                {
                    ownedData = nullptr;
                    data = nullptr;
                }
            }
        };

        std::shared_ptr<BitmapDef> def_{nullptr};
    };

    template<typename T>
    bool operator==(const Bitmap<T>& lhs, const Bitmap<T>& rhs)
    {
        const auto* lhs_data = lhs.bytes();
        const auto* rhs_data = rhs.bytes();
        return lhs_data != nullptr && rhs_data != nullptr && lhs_data == rhs_data;
    }

    template<typename T>
    bool operator!=(const Bitmap<T>& lhs, const Bitmap<T>& rhs)
    {
        return !(lhs == rhs);
    }

    template<typename T>
    struct MinMax
    {
        T min;
        T max;
    };

    template<typename T>
    struct MinMaxLoc
    {
        T min;
        T max;
        Point2i minLoc;
        Point2i maxLoc;
    };

    template<typename T>
    MinMax<T> find_min_max(const Bitmap<T>& bitmap)
    {
        auto max = std::numeric_limits<T>::min();
        auto min = std::numeric_limits<T>::max();

        const auto* ptr = bitmap.data();
        assert(ptr != nullptr);

        const auto* end = ptr + bitmap.length();

        for (; ptr != end; ++ptr)
        {
            const auto& v = *ptr;
            if (v > max) { max = v; }
            if (v < min) { min = v; }
        }

        return MinMax<T>{ min, max };
    }

    template<typename T, typename U>
    MinMax<T> find_min_max(const Bitmap<T>& bitmap, const Bitmap<U>& mask)
    {
        auto max = std::numeric_limits<T>::min();
        auto min = std::numeric_limits<T>::max();

        const auto* ptr = bitmap.data();
        assert(ptr != nullptr);

        const auto* end = ptr + bitmap.length();
        const auto* maskPtr = mask.data();

        for (; ptr != end; ++ptr, ++maskPtr)
        {
            if (*maskPtr == 0)
                continue;

            const auto& v = *ptr;
            if (v > max) { max = v; }
            if (v < min) { min = v; }
        }

        return MinMax<T>{ min, max };
    }

    template<typename T>
    MinMaxLoc<T> find_min_max_loc(const Bitmap<T>& bitmap)
    {
        auto max = std::numeric_limits<T>::min();
        auto min = std::numeric_limits<T>::max();

        Point2i maxLoc(-1, -1);
        Point2i minLoc(-1, -1);

        const auto* ptr = bitmap.data();
        assert(ptr != nullptr);

        for (int y = 0; y < bitmap.height(); ++y)
        {
            for (int x = 0; x < bitmap.width(); ++x, ++ptr)
            {
                const auto& v = *ptr;
                if (v > max) { max = v; maxLoc.x = x; maxLoc.y = y; }
                if (v < min) { min = v; minLoc.x = x; minLoc.y = y; }
            }
        }

        return MinMaxLoc<T>{ min, max, minLoc, maxLoc };
    }

    template<typename T, typename U>
    MinMaxLoc<T> find_min_max_loc(const Bitmap<T>& bitmap, const Bitmap<U>& mask)
    {
        T max = std::numeric_limits<T>::min();
        T min = std::numeric_limits<T>::max();

        Point2i maxLoc(-1, -1);
        Point2i minLoc(-1, -1);

        const auto* ptr = bitmap.data();
        assert(ptr != nullptr);

        const auto* maskPtr = mask.data();

        for (int y = 0; y < bitmap.height(); ++y)
        {
            for (int x = 0; x < bitmap.width(); ++x, ++ptr, ++maskPtr)
            {
                if (*maskPtr == 0)
                    continue;

                const auto& v = *ptr;
                if (v > max) { max = v; maxLoc.x = x; maxLoc.y = y; }
                if (v < min) { min = v; minLoc.x = x; minLoc.y = y; }
            }
        }

        return MinMaxLoc<T>{ min, max, minLoc, maxLoc };
    }

    template<typename T, typename U>
    void bitwise_or(const Bitmap<T>& left, const Bitmap<T>& right, Bitmap<U>& dest)
    {
        assert(!left.size().is_zero());
        assert(!right.size().is_zero());
        assert(left.size() == right.size());

        const auto* leftPtr = left.data();
        const auto* rightPtr = right.data();

        dest.recreate(left.size());
        auto* destPtr = dest.data();
        const auto* end = destPtr + dest.length();

        for(; destPtr != end; ++leftPtr, ++rightPtr, ++destPtr)
        {
            *destPtr = *leftPtr | *rightPtr;
        }
    }

    template<typename T, typename U>
    void bitwise_or(const Bitmap<T>& left, const Bitmap<T>& right, Bitmap<T>& dest, const Bitmap<U>& mask)
    {
        assert(!left.size().is_zero());
        assert(!right.size().is_zero());
        assert(left.size() == right.size());
        assert(left.size() == mask.size());

        const auto* leftPtr = left.data();
        const auto* rightPtr = right.data();
        const auto* maskPtr = mask.data();

        dest.recreate(left.size());
        auto* destPtr = dest.data();
        const auto* end = destPtr + dest.length();

        for(; destPtr != end; ++leftPtr, ++rightPtr, ++destPtr, ++maskPtr)
        {
            if (*maskPtr == 0)
                continue;

            *destPtr = *leftPtr | *rightPtr;
        }
    }

    template<typename T, typename U>
    void bitwise_and(const Bitmap<T>& left, const Bitmap<T>& right, Bitmap<U>& dest)
    {
        assert(!left.size().is_zero());
        assert(!right.size().is_zero());
        assert(left.size() == right.size());

        const auto* leftPtr = left.data();
        const auto* rightPtr = right.data();

        dest.recreate(left.size());
        auto* destPtr = dest.data();
        const auto* end = destPtr + dest.length();

        for(; destPtr != end; ++leftPtr, ++rightPtr, ++destPtr)
        {
            *destPtr = *leftPtr & *rightPtr;
        }
    }

    template<typename T, typename U>
    void bitwise_and(const Bitmap<T>& left, const Bitmap<T>& right, Bitmap<T>& dest, const Bitmap<U>& mask)
    {
        assert(!left.size().is_zero());
        assert(!right.size().is_zero());
        assert(left.size() == right.size());
        assert(left.size() == mask.size());

        const auto* leftPtr = left.data();
        const auto* rightPtr = right.data();
        const auto* maskPtr = mask.data();

        dest.recreate(left.size());
        auto* destPtr = dest.data();
        const auto* end = destPtr + dest.length();

        for(; destPtr != end; ++leftPtr, ++rightPtr, ++destPtr, ++maskPtr)
        {
            if (*maskPtr == 0)
                continue;

            *destPtr = *leftPtr & *rightPtr;
        }
    }

    enum class NormType
    {
        Infinite,
        L1,
        L2
    };

    template<NormType Type>
    class NormPolicy;

    template<>
    class NormPolicy<NormType::Infinite>
    {
        static double sum(double sum, double value);
        static double normalize(double sum);
    };

    template<>
    class NormPolicy<NormType::L1>
    {
        static double sum(double sum, double value);
        static double normalize(double sum);
    };

    template<>
    class NormPolicy<NormType::L2>
    {
        static double sum(double sum, double value);
        static double normalize(double sum);
    };

    template<typename T, NormType Type = NormType::L2>
    double norm(const Bitmap<T>& bitmap)
    {
        const int length = bitmap.width() * bitmap.height();

        double sum = 0.0;
        for(int i = 0; i < length; i++)
        {
            sum = NormPolicy<Type>::sum(sum, bitmap[i]);
        }

        return NormPolicy<Type>::normalize(sum);
    }

    template<typename T>
    class BitmapViewIterator;

    template<typename T>
    bool operator==(const BitmapViewIterator<T>& lhs, const BitmapViewIterator<T>& rhs)
    {
        return lhs.pt_ == rhs.pt_;
    }

    template<typename T>
    bool operator!=(const BitmapViewIterator<T>& lhs, const BitmapViewIterator<T>& rhs)
    {
        return !(lhs == rhs);
    }

    enum class MaskStyle
    {
        Normal,
        Inverted
    };

    using MaskType = std::uint8_t;

    template<typename T>
    class BitmapView
    {
    public:
        BitmapView(Bitmap<T>& src)
            : src_(src),
              owned_(false),
              mask_(nullptr),
              maskedValue_(T())
        {}

        BitmapView(Bitmap<T>& src,
                    const BitmapView<MaskType>& mask,
                    const T& maskedValue,
                    MaskStyle maskStyle = MaskStyle::Normal)
            : src_(src),
              owned_(false),
              mask_(&mask),
              maskedValue_(maskedValue),
              maskStyle_(maskStyle)
        {}

        BitmapView(Bitmap<T>& src,
                    BitmapView<MaskType>&& mask,
                    const T& maskedValue,
                    MaskStyle maskStyle = MaskStyle::Normal)
            : src_(src),
              owned_(true),
              mask_(new BitmapView<std::uint8_t>(std::move(mask))),
              maskedValue_(maskedValue),
              maskStyle_(maskStyle)
        {}

        BitmapView(const BitmapView& other)
            : src_(other.src_),
              owned_(other.owned_),
              maskedValue_(other.maskedValue_),
              maskStyle_(other.maskStyle_)
        {
            if (other.owned_)
            {
                assert(other.mask_ != nullptr);
                mask_ = new BitmapView<std::uint8_t>(*other.mask_);
            }
            else
            {
                mask_ = other.mask_;
            }
        }

        BitmapView& operator=(const BitmapView& other)
        {
            src_ = other.src_;
            owned_ = other.owned_;
            maskedValue_ = other.maskedValue_;
            maskStyle_ = other.maskStyle_;

            if (other.owned_)
            {
                assert(other.mask_ != nullptr);
                mask_ = new BitmapView<std::uint8_t>(*other.mask_);
            }
            else
            {
                mask_ = other.mask_;
            }

            return *this;
        }

        BitmapView(BitmapView&& other)
            : src_(other.src_),
              maskedValue_(other.maskedValue_),
              maskStyle_(other.maskStyle_)
        {
            std::swap(owned_, other.owned_);
            std::swap(mask_, other.mask_);
        }

        BitmapView& operator=(BitmapView&& other)
        {
            src_ = other.src_;
            maskedValue_ = other.maskedValue_;
            maskStyle_ = other.maskStyle_;

            std::swap(owned_, other.owned_);
            std::swap(mask_, other.mask_);

            return *this;
        }

        ~BitmapView()
        {
            if (owned_)
            {
                assert(mask_ != nullptr);

                delete mask_;
                mask_ = nullptr;
                owned_ = false;
            }
        }

        const T& operator[](unsigned index) const
        {
            if (!mask_)
                return src_.data()[index];

            bool maskVal = (*mask_)[index] != 0;

            if (maskStyle_ == MaskStyle::Inverted)
                maskVal = !maskVal;

            if (maskVal)
                return src_.data()[index];

            return maskedValue_;
        }

        BitmapViewIterator<T> begin() const;
        BitmapViewIterator<T> end() const;

        const T& masked_value() const { return maskedValue_; }
        const BitmapView<std::uint8_t>* mask() const { return mask_; }

        unsigned width() const { return src_.width(); }
        unsigned height() const { return src_.height(); }
        Size2i size() const { return src_.size(); }
        unsigned bpp() const { return src_.bpp(); }
        unsigned length() const { return src_.length(); }
        unsigned byte_length() const { return src_.byte_length(); }

        const Bitmap<T>& underlying_bitmap() const { return src_; }

    private:
        Bitmap<T>& src_;
        bool owned_{false};
        const BitmapView<std::uint8_t>* mask_{nullptr};
        const T maskedValue_{T()};
        MaskStyle maskStyle_{MaskStyle::Normal};
    };

    template<typename T>
    class BitmapViewIterator
    {
    public:
        BitmapViewIterator(const BitmapView<T>& view, Point2i pt)
            : view_(view), pt_(pt)
        {}

        BitmapViewIterator(BitmapViewIterator&& other)
            : view_(other.view_), pt_(other.pt_)
        {}

        BitmapViewIterator& operator=(BitmapViewIterator&& rhs)
        {
            view_ = rhs.view_;
            pt_ = rhs.pt_;

            return *this;
        }

        operator Point2i() const { return pt_; }

        unsigned x() const { return pt_.x; }
        unsigned y() const { return pt_.y; }

        BitmapViewIterator<T>& operator++()
        {
            unsigned width = view_.width();
            unsigned height = view_.height();

            Point2i b(pt_.x, pt_.y);
            b.x++;

            //if past the last column
            if (b.x >= width)
            {
                //if still have at least one more row
                if (b.y < height - 1)
                {
                    //move to beginning of the next row
                    b.y++;
                    b.x = 0;
                }
                else
                {
                    //otherwise, end of iteration
                    pt_ = end_point(view_);
                    return *this;
                }
            }

            unsigned x = b.x;

            for(unsigned y = b.y; y < height; y++, x=0)
            {
                for(; x < width; x++)
                {
                    b.x = x;
                    b.y = y;

                    unsigned index = y * width + x;
                    if (view_[index] != view_.masked_value())
                    {
                        pt_ = b;
                        return *this;
                    }
                }
            }

            pt_ = end_point(view_);
            return *this;
        }

        BitmapViewIterator<T> operator++(int)
        {
            auto pre = *this;

            ++(*this);

            return pre;
        }

        const T& operator*() const
        {
            unsigned index = view_.width() * pt_.y + pt_.x;
            return (view_.underlying_bitmap()).data()[index];
        }

        friend bool operator==<>(const BitmapViewIterator<T>& lhs, const BitmapViewIterator<T>& rhs);
        friend bool operator!=<>(const BitmapViewIterator<T>& lhs, const BitmapViewIterator<T>& rhs);

        static Point2i end_point(const BitmapView<T>& ref)
        {
            return Point2i(ref.width(), ref.height() - 1);
        }
    private:

        const BitmapView<T>& view_;
        Point2i pt_;
    };

    template<typename T>
    BitmapViewIterator<T> BitmapView<T>::begin() const
    {
        unsigned w = width();
        unsigned h = height();

        for(unsigned y = 0; y < h; y++)
        {
            for(unsigned x = 0; x < w; x++)
            {
                unsigned index = y * w + x;
                if ((*this)[index] != masked_value())
                {
                    const BitmapView<T>& ref = *this;
                    return BitmapViewIterator<T>(ref, Point2i(x,y));
                }
            }
        }

        return end();
    }

    template<typename T>
    BitmapViewIterator<T> BitmapView<T>::end() const
    {
        const BitmapView<T>& ref = *this;
        return BitmapViewIterator<T>(ref, BitmapViewIterator<T>::end_point(ref));
    }

    enum class ThresholdOp {
        Binary,
        BinaryInverse,
        Truncate,
        Zero,
        ZeroInverse
    };

    template<typename T, typename U>
    void threshold(const Bitmap<T>& src,
                   Bitmap<T>& dest,
                   U threshold,
                   U maxValue,
                   ThresholdOp op)
    {
        assert(!src.size().is_zero());
        dest.recreate(src.size());

        const auto* srcPtr = src.data();
        auto* destPtr = dest.data();

        const auto* srcPtrEnd = srcPtr + src.length();
        for (; srcPtr != srcPtrEnd; srcPtr++, destPtr++)
        {
            bool overThreshold = *srcPtr > threshold;
            switch (op)
            {
            case ThresholdOp::Binary:
                *destPtr = overThreshold ? maxValue : T(0);
                break;
            case ThresholdOp::BinaryInverse:
                *destPtr = overThreshold ? T(0) : maxValue;
                break;
            case ThresholdOp::Truncate:
                *destPtr = overThreshold ? threshold : *srcPtr;
                break;
            case ThresholdOp::Zero:
                *destPtr = overThreshold ? *srcPtr : T(0);
                break;
            case ThresholdOp::ZeroInverse:
                *destPtr = overThreshold ? T(0) : *srcPtr;
                break;
            }
        }
    }

    template<typename T, typename U, typename V>
    void band_pass(const Bitmap<T>& src,
                   Bitmap<U>& dest,
                   V lowerBound,
                   V upperBound,
                   U inBandValue,
                   ThresholdOp op)
    {
        assert(!src.size().is_zero());

        dest.recreate(src.size());

        const auto* srcPtr = src.data();
        auto* destPtr = dest.data();

        const auto* srcPtrEnd = srcPtr + src.length();
        for (; srcPtr != srcPtrEnd; srcPtr++, destPtr++)
        {
            bool under = *srcPtr < lowerBound;
            bool over = *srcPtr > upperBound;
            bool outOfBand = under || over;

            switch (op)
            {
            case ThresholdOp::Binary:
                *destPtr = outOfBand ? T(0) : inBandValue;
                break;
            case ThresholdOp::BinaryInverse:
                *destPtr = outOfBand ? inBandValue : T(0);
                break;
            case ThresholdOp::Truncate:
                *destPtr = outOfBand ?
                    (under ? lowerBound : upperBound) : *srcPtr;
                break;
            case ThresholdOp::Zero:
                *destPtr = outOfBand ? T(0) : *srcPtr;
                break;
            case ThresholdOp::ZeroInverse:
                *destPtr = outOfBand ? *srcPtr : T(0);
                break;
            }
        }
    }

    template<typename T, typename U, typename V>
    void in_range(const Bitmap<T>& src, Bitmap<U>& dest, V lowerBound, V upperBound)
    {
        assert(!src.size().is_zero());
        band_pass(src, dest, lowerBound, upperBound, static_cast<U>(1), ThresholdOp::Binary);
    }

    template<typename T>
    void range_normalize(const Bitmap<T>& src, Bitmap<T>& dest, T lowerBound, T upperBound)
    {
        assert(!src.size().is_zero());
        dest.recreate(src.size());
        const auto minMax = find_min_max(src);

        const auto length = dest.width() * dest.height();
        const auto srcRange = minMax.max - minMax.min;
        const auto dstRange = upperBound - lowerBound;

        for(int i = 0; i < length; ++i)
        {
            const auto norm = (src[i] - minMax.min) / srcRange;
            dest[i] = norm * dstRange + lowerBound;
        }
    }

    template<typename T, typename U, typename V>
    void range_normalize(const Bitmap<T>& src, Bitmap<T>& dest, V lowerBound, V upperBound, const Bitmap<U>& mask)
    {
        assert(!src.size().is_zero());
        assert(src.size() == mask.size());

        dest.recreate(src.size());
        const auto minMax = find_min_max(src, mask);

        const auto length = dest.width() * dest.height();
        const auto srcRange = minMax.max - minMax.min;
        const auto dstRange = upperBound - lowerBound;

        for(int i = 0; i < length; ++i)
        {
            if (mask[i] == 0)
                continue;

            const auto norm = (src[i] - minMax.min) / srcRange;
            dest[i] = norm * dstRange + lowerBound;
        }
    }

    template<typename T>
    void resize(const Bitmap<T>& src, Bitmap<T>& dest)
    {
        resize(src, dest, dest.size());
    }

    template<typename T>
    void resize(const Bitmap<T>& src, Bitmap<T>& dest, Size2i size)
    {
        assert(!src.size().is_zero());

        dest.recreate(size);

        const float scaleX = src.width() / static_cast<float>(dest.width());
        const float scaleY = src.height() / static_cast<float>(dest.height());

        const unsigned srcWidth = src.width();
        const unsigned destWidth = dest.width();
        const unsigned destHeight = dest.height();

        const T* srcPtr = src.data();
        T* destPtr = dest.data();

        for(unsigned y = 0; y < destHeight; y++)
        {
            const unsigned srcY = std::max<unsigned>(y * scaleY, 0);

            const T* srcRowPtr = srcPtr + (srcY * srcWidth);
            T* destRowPtr = destPtr + (y * destWidth);

            for(unsigned x = 0; x < destWidth; x++)
            {
                const unsigned srcX = std::max<unsigned>(x * scaleX, 0);
                destRowPtr[x] = srcRowPtr[srcX];
            }
        }
    }

    template<typename T, typename U>
    void convert_to(const Bitmap<T>& src, Bitmap<U>& dest)
    {
        assert(!src.size().is_zero());
        dest.recreate(src.size());

        const unsigned width = src.width();
        const unsigned height = src.height();

        const T* srcPtr = src.data();
        U* destPtr = dest.data();

        for(unsigned y = 0; y < height; y++)
        {
            const unsigned rowIndex = y * width;

            const T* srcRowPtr = &srcPtr[rowIndex];
            U* destRowPtr = &destPtr[rowIndex];

            for(unsigned x = 0; x < width; x++)
            {
                destRowPtr[x] = static_cast<U>(srcRowPtr[x]);
            }
        }
    }

    template<typename T, typename U, typename TFunc>
    void convert_to(const Bitmap<T>& src, Bitmap<U>& dest, TFunc converter)
    {
        assert(!src.size().is_zero());
        dest.recreate(src.size());

        const unsigned width = src.width();
        const unsigned height = src.height();

        const T* srcPtr = src.data();
        U* destPtr = dest.data();

        for(unsigned y = 0; y < height; y++)
        {
            const unsigned rowIndex = y * width;

            const T* srcRowPtr = &srcPtr[rowIndex];
            U* destRowPtr = &destPtr[rowIndex];

            for(unsigned x = 0; x < width; x++)
            {
                destRowPtr[x] = converter(srcRowPtr[x]);
            }
        }
    }

    template<typename T>
    void copy_to(const Bitmap<T>& src, Bitmap<T>& dest)
    {
        assert(!src.size().is_zero());
        dest.recreate(src.size());

        const unsigned width = src.width();
        const unsigned height = src.height();

        const T* srcPtr = src.data();
        T* destPtr = dest.data();

        std::copy(srcPtr, srcPtr+(width*height), destPtr);
    }

    // template<typename T>
    // void copy_to(const Bitmap<T>& src, const RectI& srcWindow, Bitmap<T>& dest, const RectI& destWindow)
    // {
    //     dest.recreate(src.size());

    //     const unsigned width = src.width();
    //     const unsigned height = src.height();

    //     const T* srcPtr = src.data();
    //     T* destPtr = dest.data();

    //     std::copy(srcPtr, srcPtr+(width*height), destPtr);
    // }

    template<typename T, typename U>
    void copy_to(const Bitmap<T>& src, Bitmap<T>& dest, const Bitmap<U>& mask)
    {
        assert(!src.size().is_zero());
        assert(src.size() == mask.size());
        dest.recreate(src.size());

        const unsigned width = src.width();
        const unsigned height = src.height();

        const T* srcPtr = src.data();
        T* destPtr = dest.data();

        const MaskType* maskPtr = mask.data();

        for(unsigned y = 0; y < height; y++)
        {
            const unsigned rowIndex = y * width;

            const T* srcRowPtr = srcPtr + rowIndex;
            T* destRowPtr = destPtr + rowIndex;
            const MaskType* maskRowPtr = maskPtr + rowIndex;

            for(unsigned x = 0; x < width; x++)
            {
                const MaskType maskPixel = maskRowPtr[x];

                if (maskPixel == 0)
                    continue;

                destRowPtr[x] = srcRowPtr[x];
            }
        }
    }

    template<typename T>
    void scalar_add(const Bitmap<T>& left, const Bitmap<T>& right, Bitmap<T>& dest)
    {
        assert(!left.size().is_zero());
        assert(!right.size().is_zero());
        assert(left.size() == right.size());

        dest.recreate(left.size());
        const auto length = dest.width() * dest.height();

        for(int i = 0; i < length; ++i)
        {
            dest[i] = left[i] + right[i];
        }
    }

    template<typename T, typename U>
    void scalar_add(const Bitmap<T>& left, const Bitmap<T>& right, Bitmap<T>& dest, const Bitmap<U>& mask)
    {
        assert(!left.size().is_zero());
        assert(!right.size().is_zero());
        assert(left.size() == right.size());
        assert(left.size() == mask.size());

        dest.recreate(left.size());
        const auto length = dest.width() * dest.height();

        for(int i = 0; i < length; ++i)
        {
            if (mask[i] == 0)
                continue;

            dest[i] = left[i] + right[i];
        }
    }

    template<typename T>
    void scalar_subtract(const Bitmap<T>& left, const Bitmap<T>& right, Bitmap<T>& dest)
    {
        assert(!left.size().is_zero());
        assert(!right.size().is_zero());
        assert(left.size() == right.size());

        dest.recreate(left.size());
        const auto length = dest.width() * dest.height();

        for(int i = 0; i < length; ++i)
        {
            dest[i] = left[i] - right[i];
        }
    }

    template<typename T, typename U>
    void scalar_subtract(const Bitmap<T>& left, const Bitmap<T>& right, Bitmap<T>& dest, const Bitmap<U>& mask)
    {
        assert(!left.size().is_zero());
        assert(!right.size().is_zero());
        assert(left.size() == right.size());
        assert(left.size() == mask.size());

        dest.recreate(left.size());
        const auto length = dest.width() * dest.height();

        for(int i = 0; i < length; ++i)
        {
            if (mask[i] == 0)
                continue;

            dest[i] = left[i] - right[i];
        }
    }

    template<typename T>
    void scalar_multiply(const Bitmap<T>& left, const Bitmap<T>& right, Bitmap<T>& dest)
    {
        assert(!left.size().is_zero());
        assert(!right.size().is_zero());
        assert(left.size() == right.size());

        dest.recreate(left.size());
        const auto length = dest.width() * dest.height();

        for(int i = 0; i < length; ++i)
        {
            dest[i] = left[i] * right[i];
        }
    }

    template<typename T, typename U>
    void scalar_multiply(const Bitmap<T>& left, const Bitmap<T>& right, Bitmap<T>& dest, const Bitmap<U>& mask)
    {
        assert(!left.size().is_zero());
        assert(!right.size().is_zero());
        assert(left.size() == right.size());
        assert(left.size() == mask.size());

        dest.recreate(left.size());
        const auto length = dest.width() * dest.height();

        for(int i = 0; i < length; ++i)
        {
            if (mask[i] == 0)
                continue;

            dest[i] = left[i] * right[i];
        }
    }

    template<typename T>
    void scalar_divide(const Bitmap<T>& left, const Bitmap<T>& right, Bitmap<T>& dest)
    {
        assert(!left.size().is_zero());
        assert(!right.size().is_zero());
        assert(left.size() == right.size());

        dest.recreate(left.size());
        const auto length = dest.width() * dest.height();

        for(int i = 0; i < length; ++i)
        {
            dest[i] = left[i] / right[i];
        }
    }

    template<typename T, typename U>
    void scalar_divide(const Bitmap<T>& left, const Bitmap<T>& right, Bitmap<T>& dest, const Bitmap<U>& mask)
    {
        assert(!left.size().is_zero());
        assert(!right.size().is_zero());
        assert(left.size() == right.size());
        assert(left.size() == mask.size());

        dest.recreate(left.size());
        const auto length = dest.width() * dest.height();

        for(int i = 0; i < length; ++i)
        {
            if (mask[i] == 0)
                continue;

            dest[i] = left[i] / right[i];
        }
    }

    template<typename T, typename TSeed, typename FoldFunc>
    TSeed foldl(const Bitmap<T>& bitmap, TSeed seed, FoldFunc f)
    {
        assert(!bitmap.size().is_zero());

        auto result = seed;
        for(int i = 0; i < bitmap.length(); ++i)
        {
            result = f(result, bitmap[i]);
        }

        return result;
    }

    template<typename T, typename TSeed, typename FoldFunc>
    TSeed foldr(const Bitmap<T>& bitmap, FoldFunc f, TSeed seed)
    {
        assert(!bitmap.size().is_zero());

        auto result = seed;
        for(int i = bitmap.length() - 1; i >= 0; --i)
        {
            result = f(bitmap[i], result);
        }

        return result;
    }

    template<typename T>
    int count_non_zero(const Bitmap<T>& bitmap)
    {
        assert(!bitmap.size().is_zero());

        return foldl(bitmap,
                     0,
                     [](int acc, T val)
                     {
                         if (val != 0) { return acc + 1; }
                         return acc;
                     });
    }

    template<typename T, typename UnaryPredicate>
    bool any(const Bitmap<T>& bitmap, UnaryPredicate predicate)
    {
        const auto* ptr = bitmap.data();
        assert(ptr != nullptr);

        const auto* end = ptr + bitmap.length();

        for (; ptr != end; ++ptr)
        {
            const auto& v = *ptr;
            if (predicate(v))
                return true;
        }

        return false;
    }

    template<typename T, typename U, typename UnaryPredicate>
    bool any(const Bitmap<T>& bitmap, UnaryPredicate predicate, const Bitmap<U>& mask)
    {
        assert(bitmap.size() == mask.size());

        const auto* ptr = bitmap.data();
        assert(ptr != nullptr);

        const auto* end = ptr + bitmap.length();
        const auto* maskPtr = mask.data();

        for (; ptr != end; ++ptr, ++maskPtr)
        {
            if (*maskPtr == 0)
                continue;

            const auto& v = *ptr;
            if (predicate(v))
                return true;
        }

        return false;
    }

    template<typename T, typename UnaryPredicate>
    bool all(const Bitmap<T>& bitmap, UnaryPredicate predicate)
    {
        const auto* ptr = bitmap.data();
        assert(ptr != nullptr);

        const auto* end = ptr + bitmap.length();

        for (; ptr != end; ++ptr)
        {
            const auto& v = *ptr;
            if (!predicate(v))
                return false;
        }

        return true;
    }

    template<typename T, typename U, typename UnaryPredicate>
    bool all(const Bitmap<T>& bitmap, UnaryPredicate predicate, const Bitmap<U>& mask)
    {
        assert(bitmap.size() == mask.size());

        const auto* ptr = bitmap.data();
        assert(ptr != nullptr);

        const auto* end = ptr + bitmap.length();
        const auto* maskPtr = mask.data();

        for (; ptr != end; ++ptr, ++maskPtr)
        {
            if (*maskPtr == 0)
                continue;

            const auto& v = *ptr;
            if (!predicate(v))
                return false;
        }

        return true;
    }

    template<typename T>
    bool all_zero(const Bitmap<T>& bitmap)
    {
        return all(bitmap, [](const T& v) -> bool { return v == 0; });
    }

    template<typename T, typename U>
    bool all_zero(const Bitmap<T>& bitmap, const Bitmap<U>& mask)
    {
        return all(bitmap, [](const T& v) -> bool { return v == 0; }, mask);
    }

    using BitmapMask = Bitmap<MaskType>;
    using BitmapU8 = Bitmap<std::uint8_t>;
    using BitmapI16 = Bitmap<std::int16_t>;
    using BitmapU16 = Bitmap<std::uint16_t>;
    using BitmapI32 = Bitmap<std::int32_t>;
    using BitmapU32 = Bitmap<std::int32_t>;
    using BitmapF = Bitmap<float>;
    using BitmapVector3f = Bitmap<Vector3f>;
    using BitmapRGB = Bitmap<RGBPixel>;
    using BitmapRGBA = Bitmap<RGBAPixel>;

}}

#endif /* BITMAP_HPP */
