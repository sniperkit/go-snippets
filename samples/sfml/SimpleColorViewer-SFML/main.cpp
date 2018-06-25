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
#include <SFML/Graphics.hpp>
#include <astra_core/astra_core.hpp>
#include <astra/astra.hpp>

#include <chrono>
#include <iostream>
#include <iomanip>
#include <cstring>
#include <key_handler.h>

class ColorFrameListener : public astra::FrameListener
{
public:
    ColorFrameListener()
    {
        prev_ = ClockType::now();
    }

    void init_texture(int width, int height)
    {
        if (displayBuffer_ == nullptr || width != displayWidth_ || height != displayHeight_)
        {
            displayWidth_ = width;
            displayHeight_ = height;

            // texture is RGBA
            int byteLength = displayWidth_ * displayHeight_ * 4;

            displayBuffer_ = BufferPtr(new uint8_t[byteLength]);
            std::memset(displayBuffer_.get(), 0, byteLength);

            texture_.create(displayWidth_, displayHeight_);
            sprite_.setTexture(texture_);
            sprite_.setPosition(0, 0);
        }
    }

    void check_fps()
    {
        const float frameWeight = .2f;

        const ClockType::time_point now = ClockType::now();
        const float elapsedMillis = std::chrono::duration_cast<DurationType>(now - prev_).count();

        elapsedMillis_ = elapsedMillis * frameWeight + elapsedMillis_ * (1.f - frameWeight);
        prev_ = now;

        const float fps = 1000.f / elapsedMillis;

        const auto precision = std::cout.precision();

        std::cout << std::fixed
                  << std::setprecision(1)
                  << fps << " fps ("
                  << std::setprecision(1)
                  << elapsedMillis_ << " ms)"
                  << std::setprecision(precision)
                  << std::endl;
    }

    virtual void on_frame_ready(astra::StreamReader& reader, astra::Frame& frame) override
    {
        const astra::ColorFrame colorFrame = frame.get<astra::ColorFrame>();

        int width = colorFrame.width();
        int height = colorFrame.height();

        init_texture(width, height);

        const astra::RgbPixel* colorData = colorFrame.data();

        for (int i = 0; i < width * height; i++)
        {
            int rgbaOffset = i * 4;
            displayBuffer_[rgbaOffset] = colorData[i].r;
            displayBuffer_[rgbaOffset + 1] = colorData[i].g;
            displayBuffer_[rgbaOffset + 2] = colorData[i].b;
            displayBuffer_[rgbaOffset + 3] = 255;
        }

        texture_.update(displayBuffer_.get());
        check_fps();
    }

    void drawTo(sf::RenderWindow& window)
    {
        if (displayBuffer_ != nullptr)
        {
            float imageScale = window.getView().getSize().x / displayWidth_;
            sprite_.setScale(imageScale, imageScale);
            window.draw(sprite_);
        }
    }

private:
    using DurationType = std::chrono::milliseconds;
    using ClockType = std::chrono::high_resolution_clock;

    ClockType::time_point prev_;
    float elapsedMillis_{.0f};

    sf::Texture texture_;
    sf::Sprite sprite_;

    using BufferPtr = std::unique_ptr<uint8_t[]>;
    BufferPtr displayBuffer_{nullptr};

    int displayWidth_{0};
    int displayHeight_{0};
};

int main(int argc, char** argv)
{
    astra::initialize();

    set_key_handler();

    sf::RenderWindow window(sf::VideoMode(1280, 960), "Color Viewer");

    astra::StreamSet streamSet;
    astra::StreamReader reader = streamSet.create_reader();

    reader.stream<astra::ColorStream>().start();

    ColorFrameListener listener;
    reader.add_listener(listener);

    while (window.isOpen())
    {
        astra_temp_update();

        sf::Event event;
        while (window.pollEvent(event))
        {
            switch (event.type)
            {
            case sf::Event::Closed:
                window.close();
                break;
            case sf::Event::KeyPressed:
                {
                    if (event.key.code == sf::Keyboard::Escape ||
                        (event.key.code == sf::Keyboard::C && event.key.control))
                    {
                        window.close();
                    }
                }
            default:
                break;
            }
        }

        // clear the window with black color
        window.clear(sf::Color::Black);

        listener.drawTo(window);
        window.display();

        if (!shouldContinue)
        {
            window.close();
        }
    }

    astra::terminate();
    return 0;
}
