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
#include <astra/streams/Depth.hpp>
#include <astra_core/capi/plugins/astra_plugin.h>
#include <astra_core/capi/astra_core.h>
#include <astra_core/capi/astra_types.h>

#include <memory>
#include <chrono>
#include <iostream>
#include <iomanip>

#include <SFML/Graphics.hpp>

#include "LitDepthVisualizer.hpp"
#include <common/serialization/FrameStreamWriter.h>
#include <common/serialization/FrameOutputStream.h>

using namespace astra;

class Recorder
{
public:
    explicit Recorder(const char* filename) {
        outputFile_ = fopen(filename, "wb");
        frameOutputStream_ = serialization::open_frame_output_stream(outputFile_);
        frameStreamWriter_ = FrameStreamWriterPtr(new serialization::FrameStreamWriter(*frameOutputStream_));
        frameStreamWriter_->begin_write();
    }

    ~Recorder() {
        frameStreamWriter_->end_write();
        serialization::close_frame_output_stream(frameOutputStream_);
        frameStreamWriter_ = nullptr;
        fclose(outputFile_);
        printf("closed file\n");
    }

    void add_frame(const DepthFrame& depthFrame) {
        if (!depthFrame.is_valid()) {
            return;
        }
        bool result = frameStreamWriter_->write(depthFrame);
        printf("Saving frame: %d %s\n", frameCount_, result ? "" : "failure");
        ++frameCount_;
    }

private:
    FILE* outputFile_;
    serialization::FrameOutputStream* frameOutputStream_;
    int frameCount_ { 0 };

    using FrameStreamWriterPtr = std::unique_ptr<serialization::FrameStreamWriter>;
    FrameStreamWriterPtr frameStreamWriter_;
};

class Viewer : public FrameListener
{
public:
    Viewer(StreamSet& streamset) :
        reader_(streamset.create_reader())
    {
        lastTimepoint_ = clock_type::now();

        reader_.stream<DepthStream>().start();
        reader_.stream<PointStream>().start();

        reader_.add_listener(*this);
    }

    void draw_to(sf::RenderWindow& window)
    {
        if (displayBuffer_ != nullptr)
        {
            float depthScale = window.getView().getSize().x / displayWidth_;

            sprite_.setScale(depthScale, depthScale);

            window.draw(sprite_);
        }
    }

    void start_recording()
    {
        recorder_ = RecorderPtr(new Recorder("test.df"));
    }

    void stop_recording()
    {
        recorder_ = nullptr;
    }

private:
    void init_texture(int width, int height)
    {
        if (displayBuffer_ == nullptr || width != displayWidth_ || height != displayHeight_)
        {
            displayWidth_ = width;
            displayHeight_ = height;

            // texture is RGBA
            int byteLength = displayWidth_ * displayHeight_ * 4;

            displayBuffer_ = BufferPtr(new uint8_t[byteLength]);
            memset(displayBuffer_.get(), 0, byteLength);

            texture_.create(displayWidth_, displayHeight_);
            sprite_.setTexture(texture_);
            sprite_.setPosition(0, 0);
        }
    }

    void check_fps()
    {
        const double frameWeight = 0.2;

        auto newTimepoint = clock_type::now();
        auto frameDuration = std::chrono::duration_cast<duration_type>(newTimepoint - lastTimepoint_);

        frameDuration_ = frameDuration * frameWeight + frameDuration_ * (1 - frameWeight);
        lastTimepoint_ = newTimepoint;

        double fps = 1.0 / frameDuration_.count();

        auto precision = std::cout.precision();
        std::cout << std::fixed
                  << std::setprecision(1)
                  << fps << " fps ("
                  << std::setprecision(2)
                  << frameDuration.count() * 1000 << " ms)"
                  << std::setprecision(precision)
                  << std::endl;
    }

    void visualize_frame(const PointFrame& pointFrame)
    {
        if (!pointFrame.is_valid()) {
            return;
        }
        int width = pointFrame.width();
        int height = pointFrame.height();

        init_texture(width, height);

        visualizer_.update(pointFrame);

        astra_rgb_pixel_t* vizBuffer = visualizer_.get_output();
        for (int i = 0; i < width * height; i++)
        {
            int rgbaOffset = i * 4;
            displayBuffer_[rgbaOffset] = vizBuffer[i].r;
            displayBuffer_[rgbaOffset + 1] = vizBuffer[i].b;
            displayBuffer_[rgbaOffset + 2] = vizBuffer[i].g;
            displayBuffer_[rgbaOffset + 3] = 255;
        }
        texture_.update(displayBuffer_.get());
    }

    virtual void on_frame_ready(StreamReader& reader,
                                Frame& frame) override
    {
        const PointFrame pointFrame = frame.get<PointFrame>();
        const DepthFrame depthFrame = frame.get<DepthFrame>();

        if (recorder_ != nullptr) {
            recorder_->add_frame(depthFrame);
        }

        visualize_frame(pointFrame);

        //check_fps();
    }

    samples::common::LitDepthVisualizer visualizer_;

    using duration_type = std::chrono::duration<double>;
    duration_type frameDuration_{0.0};

    using clock_type = std::chrono::system_clock;
    std::chrono::time_point<clock_type> lastTimepoint_;
    sf::Texture texture_;
    sf::Sprite sprite_;

    using BufferPtr = std::unique_ptr<uint8_t[]>;
    BufferPtr displayBuffer_{ nullptr };
    int displayWidth_{0};
    int displayHeight_{0};

    StreamReader reader_;
    using RecorderPtr = std::unique_ptr<Recorder>;
    RecorderPtr recorder_;
};

void handle_key(sf::Keyboard::Key key, sf::RenderWindow& window, Viewer& viewer)
{
    if (key == sf::Keyboard::Escape)
    {
        window.close();
    }
    else if (key == sf::Keyboard::S)
    {
        viewer.stop_recording();
    }
    if (key == sf::Keyboard::R)
    {
        viewer.start_recording();
    }
}

int main(int argc, char** argv)
{
    astra::initialize();

    sf::RenderWindow window(sf::VideoMode(1280, 960), "Stream Recorder");

    StreamSet streamset;

/*
    streamset streamPlayer("stream_player");
    StreamReader streamPlayerReader = streamPlayer.create_reader();

    auto streamPlayerPs = streamPlayerReader.stream<pointstream>();
*/
    Viewer viewer(streamset);

    while (window.isOpen())
    {
        astra_temp_update();

        sf::Event event;

        window.clear(sf::Color::Black);

        while (window.pollEvent(event))
        {
            switch (event.type)
            {
            case sf::Event::Closed:
                window.close();
                break;
            case sf::Event::KeyPressed:
            {
                sf::Keyboard::Key key = event.key.code;
                handle_key(key, window, viewer);
                break;
            }
            default:
                break;
            }
        }

        viewer.draw_to(window);

        window.display();
    }

    astra::terminate();

    return 0;
}
