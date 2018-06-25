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
#include <astra/astra.hpp>
#include <iostream>
#include <cstring>

class skeleton_visualizer : public astra::FrameListener
{
public:
    void init_texture(int width, int height)
    {
        if (displayBuffer_ == nullptr || width != depthWidth_ || height != depthHeight_)
        {
            depthWidth_ = width;
            depthHeight_ = height;
            int byteLength = depthWidth_ * depthHeight_ * 4;

            displayBuffer_ = BufferPtr(new uint8_t[byteLength]);
            std::memset(displayBuffer_.get(), 0, byteLength);

            texture_.create(depthWidth_, depthHeight_);
            sprite_.setTexture(texture_);
            sprite_.setPosition(0, 0);
        }
    }

    void check_fps()
    {
        double fpsFactor = 0.02;

        std::clock_t newTimepoint= std::clock();
        long double frameDuration = (newTimepoint - lastTimepoint_) / static_cast<long double>(CLOCKS_PER_SEC);

        frameDuration_ = frameDuration * fpsFactor + frameDuration_ * (1 - fpsFactor);
        lastTimepoint_ = newTimepoint;
        double fps = 1.0 / frameDuration_;

        printf("FPS: %3.1f (%3.4Lf ms)\n", fps, frameDuration_ * 1000);
    }

    void processDepth(astra::Frame& frame)
    {
        const astra::DepthFrame depthFrame = frame.get<astra::DepthFrame>();

        int width = depthFrame.width();
        int height = depthFrame.height();

        init_texture(width, height);

        const int16_t* depthPtr = depthFrame.data();
        for(int y = 0; y < height; y++)
        {
            for(int x = 0; x < width; x++)
            {
                int index = (x + y * width);
                int index4 = index * 4;

                int16_t depth = depthPtr[index];
                uint8_t value = depth % 255;

                displayBuffer_[index4] = value;
                displayBuffer_[index4 + 1] = value;
                displayBuffer_[index4 + 2] = value;
                displayBuffer_[index4 + 3] = 255;
            }
        }

        texture_.update(displayBuffer_.get());
    }

    void processSkeletons(astra::Frame& frame)
    {
        astra::skeletonframe skeletonFrame = frame.get<astra::skeletonframe>();

        skeletons_ = skeletonFrame.skeletons();
        jointPositions_.clear();

        for (auto skeleton : skeletons_)
        {
            for(auto joint : skeleton.joints())
            {
                auto depthPosition =
                    mapper_->convert_world_to_depth(joint.position());

                jointPositions_.push_back(depthPosition);
            }
        }
    }

    virtual void on_frame_ready(astra::StreamReader& reader,
                                astra::Frame& frame) override
    {
        if (mapper_ == nullptr)
        {
            auto& mapper = reader.stream<astra::DepthStream>().coordinateMapper();
            mapper_ = astra::make_unique<astra::CoordinateMapper>(mapper);
        }

        processDepth(frame);
        processSkeletons(frame);

        check_fps();
    }

    void draw_circle(sf::RenderWindow& window, float radius, float x, float y, sf::Color color)
    {
        sf::CircleShape shape(radius);

        shape.setFillColor(color);

        shape.setOrigin(radius, radius);
        shape.setPosition(x, y);
        window.draw(shape);
    }

    void draw_skeletons(sf::RenderWindow& window, float depthScale)
    {
        float radius = 16;
        sf::Color trackingColor(10, 10, 200);

        for (auto position : jointPositions_)
        {
            draw_circle(window,
                       radius,
                       position.x * depthScale,
                       position.y * depthScale,
                       trackingColor);
        }
    }

    void draw_to(sf::RenderWindow& window)
    {
        if (displayBuffer_ != nullptr)
        {
            float depthScale = window.getView().getSize().x / depthWidth_;

            sprite_.setScale(depthScale, depthScale);

            window.draw(sprite_);

            draw_skeletons(window, depthScale);
        }
    }

private:
    long double frameDuration_{ 0 };
    std::clock_t lastTimepoint_ { 0 };
    sf::Texture texture_;
    sf::Sprite sprite_;

    using BufferPtr = std::unique_ptr < uint8_t[] >;
    BufferPtr displayBuffer_{ nullptr };

    std::unique_ptr<astra::CoordinateMapper> mapper_;
    std::vector<astra::skeleton> skeletons_;
    std::vector<astra::Vector3f> jointPositions_;

    int depthWidth_{0};
    int depthHeight_{0};
};

int main(int argc, char** argv)
{
    astra::initialize();

    sf::RenderWindow window(sf::VideoMode(1280, 960), "Skeleton Viewer");

    astra::StreamSet sensor;
    astra::StreamReader reader = sensor.create_reader();

    skeleton_visualizer listener;

    reader.stream<astra::DepthStream>().start();
    reader.stream<astra::skeletonstream>().start();
    reader.stream<astra::skeletonstream>().set_zMin(5);
    reader.stream<astra::skeletonstream>().set_zMax(1000);
    reader.add_listener(listener);

    while (window.isOpen())
    {
        astra_temp_update();

        sf::Event event;
        while (window.pollEvent(event))
        {
            if (event.type == sf::Event::Closed)
                window.close();
            if ((event.type == sf::Event::KeyPressed) && (event.key.code == sf::Keyboard::Escape))
                window.close();
        }

        // clear the window with black color
        window.clear(sf::Color::Black);

        listener.draw_to(window);
        window.display();
    }

    astra::terminate();

    return 0;
}
