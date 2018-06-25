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
#ifndef HND_MORPHOLOGY_HPP
#define HND_MORPHOLOGY_HPP

#include "hnd_bitmap.hpp"

namespace astra { namespace hand {

    enum class MorphShape
    {
        Rect,
        Cross
    };

    inline BitmapMask get_structuring_element(MorphShape shape, Size2i size)
    {
        BitmapMask element(size);
        int width = size.width();
        int height = size.height();
        int centerX = width >> 1;
        int centerY = height >> 1;

        MaskType* data = element.data();
        for (int y = 0; y < height; ++y, data += width)
        {
            MaskType* row = data;
            //default: MorphShape::Rect
            int maskStart = 0;
            int maskStop = width;
            if (shape == MorphShape::Cross && y != centerY)
            {
                maskStart = centerX;
                maskStop = maskStart+1;
            }

            for (int x = 0; x < maskStart; ++x)
            {
                row[x] = 0;
            }
            for (int x = maskStart; x < maskStop; ++x)
            {
                row[x] = 1;
            }
            for (int x = maskStop; x < width; ++x)
            {
                row[x] = 0;
            }
        }
        return element;
    }

    //TODO look at mallocs
    template<typename T>
    void erode(const Bitmap<T>& input, Bitmap<T>& output, const BitmapMask& element)
    {
        Bitmap<T> temp;
        const Bitmap<T>* source = &input;
        if (input == output)
        {
            //input and output are the same. Create a temporary bitmap for the source data
            temp = input.clone();
            source = &temp;
        } else {
            //enforce allocation and same size as the input
            output.recreate(input.size());
        }

        //Use source instead of input from here down

        const int width = source->width();
        const int height = source->height();
        const int elementWidth = element.width();
        const int elementHeight = element.height();
        const int halfElementHeight = elementHeight >> 1;
        const int halfElementWidth = elementWidth >> 1;

        const T borderValue = std::numeric_limits<T>::max();

        const T* const sourceData = source->data();
        T* outputData = output.data();
        const MaskType* elementData = element.data();

        const T** rows = new const T*[elementHeight];

        for (int y = 0; y < height; ++y)
        {
            for (int i = 0; i < elementHeight; ++i)
            {
                int rowY = y + i - halfElementHeight;
                rows[i] = &sourceData[rowY * width];
            }

            for (int x = 0; x < width; ++x)
            {
                //apply kernel:
                T minValue = borderValue;
                for (int elemY = 0; elemY < elementHeight; ++elemY)
                {
                    int currentY = y + elemY - halfElementHeight;
                    const T* currentRow = rows[elemY];
                    for (int elemX = 0; elemX < elementWidth; ++elemX)
                    {
                        int elementIndex = elemX + elemY * elementWidth;
                        //Only check if the kernel element is non-zero
                        if (elementData[elementIndex] != 0)
                        {
                            int currentX = x + elemX - halfElementWidth;
                            if (currentX >= 0 && currentX < width &&
                                currentY >= 0 && currentY < height)
                            {
                                //in range of input
                                T currentValue = currentRow[currentX];
                                //Erode is the MIN operation
                                if (currentValue < minValue)
                                {
                                    minValue = currentValue;
                                }
                            }
                            //If out of range, would use the borderValue.
                            //For erode, the default value is borderValue, so no check required.
                        }
                    }
                }
                int outIndex = x + y * width;
                outputData[outIndex] = minValue;
            }
        }

        delete[] rows;
    }

        //TODO look at mallocs
    template<typename T>
    void dilate(const Bitmap<T>& input, Bitmap<T>& output, const BitmapMask& element)
    {
        Bitmap<T> temp;
        const Bitmap<T>* source = &input;
        if (input == output)
        {
            //input and output are the same. Create a temporary bitmap for the source data
            temp = input.clone();
            source = &temp;
        } else {
            //enforce allocation and same size as the input
            output.recreate(input.size());
        }

        //Use source instead of input from here down

        const int width = source->width();
        const int height = source->height();
        const int elementWidth = element.width();
        const int elementHeight = element.height();
        const int halfElementHeight = elementHeight >> 1;
        const int halfElementWidth = elementWidth >> 1;

        const T borderValue = std::numeric_limits<T>::min();

        const T* const sourceData = source->data();
        T* outputData = output.data();
        const MaskType* elementData = element.data();

        const T** rows = new const T*[elementHeight];

        for (int y = 0; y < height; ++y)
        {
            for (int i = 0; i < elementHeight; ++i)
            {
                int rowY = y + i - halfElementHeight;
                rows[i] = &sourceData[rowY * width];
            }

            for (int x = 0; x < width; ++x)
            {
                //apply kernel:
                T maxValue = borderValue;
                for (int elemY = 0; elemY < elementHeight; ++elemY)
                {
                    int currentY = y + elemY - halfElementHeight;
                    const T* currentRow = rows[elemY];
                    for (int elemX = 0; elemX < elementWidth; ++elemX)
                    {
                        int elementIndex = elemX + elemY * elementWidth;
                        //Only check if the kernel element is non-zero
                        if (elementData[elementIndex] != 0)
                        {
                            int currentX = x + elemX - halfElementWidth;
                            if (currentX >= 0 && currentX < width &&
                                currentY >= 0 && currentY < height)
                            {
                                //in range of input
                                T currentValue = currentRow[currentX];
                                //Dilate is the MAX operation
                                if (currentValue > maxValue)
                                {
                                    maxValue = currentValue;
                                }
                            }
                            //If out of range, would use the borderValue.
                            //For dilate, the default value is borderValue, so no check required.
                        }
                    }
                }
                int outIndex = x + y * width;
                outputData[outIndex] = maxValue;
            }
        }

        delete[] rows;
    }
}}


#endif /* HND_MORPHOLOGY_HPP */
