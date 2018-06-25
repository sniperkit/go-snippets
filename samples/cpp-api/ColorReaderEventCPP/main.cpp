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
#include <astra/astra.hpp>
#include <cstdio>
#include <iostream>
#include <key_handler.h>

class SampleFrameListener : public astra::FrameListener
{
private:
    using buffer_ptr = std::unique_ptr<astra::RgbPixel []>;
    buffer_ptr buffer_;
    unsigned int lastWidth_;
    unsigned int lastHeight_;

public:
    virtual void on_frame_ready(astra::StreamReader& reader,
                                astra::Frame& frame) override
    {
        const astra::ColorFrame colorFrame = frame.get<astra::ColorFrame>();

        if (colorFrame.is_valid())
        {
            print_color(colorFrame);
        }
    }

    void print_color(const astra::ColorFrame& colorFrame)
    {
        if (colorFrame.is_valid())
        {
            int width = colorFrame.width();
            int height = colorFrame.height();
            int frameIndex = colorFrame.frame_index();

            if (width != lastWidth_ || height != lastHeight_){
                buffer_ = buffer_ptr(new astra::RgbPixel[colorFrame.length()]);
                lastWidth_ = width;
                lastHeight_ = height;
            }
            colorFrame.copy_to(buffer_.get());

            size_t index = ((width * (height / 2.0f)) + (width / 2.0f));
            astra::RgbPixel middle = buffer_[index];

            std::cout << "color frameIndex: " << frameIndex
                      << " r: " << static_cast<int>(middle.r)
                      << " g: " << static_cast<int>(middle.g)
                      << " b: " << static_cast<int>(middle.b)
                      << std::endl;
        }
    }
};

int main(int argc, char** argv)
{
    astra::initialize();

    set_key_handler();

    astra::StreamSet streamSet;
    astra::StreamReader reader = streamSet.create_reader();

    SampleFrameListener listener;

    reader.stream<astra::ColorStream>().start();

    std::cout << "colorStream -- hFov: "
              << reader.stream<astra::ColorStream>().hFov()
              << " vFov: "
              << reader.stream<astra::ColorStream>().vFov()
              << std::endl;

    reader.add_listener(listener);

    do
    {
        astra_temp_update();
    } while (shouldContinue);

    reader.remove_listener(listener);

    astra::terminate();
}