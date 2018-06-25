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
#include <sstream>
#include <iomanip>
#include <deque>
#include <unordered_map>
#include <chrono>
#include <SFML/Graphics.hpp>
#include <astra_core/astra_core.hpp>
#include <astra/astra.hpp>
#include "LitDepthVisualizer.hpp"
#include <key_handler.h>

class sfLine : public sf::Drawable
{
public:
    sfLine(const sf::Vector2f& point1, const sf::Vector2f& point2, sf::Color color, float thickness)
        : color_(color)
    {
        const sf::Vector2f direction = point2 - point1;
        const sf::Vector2f unitDirection = direction / std::sqrt(direction.x*direction.x + direction.y*direction.y);
        const sf::Vector2f normal(-unitDirection.y, unitDirection.x);

        const sf::Vector2f offset = (thickness / 2.f) * normal;

        vertices_[0].position = point1 + offset;
        vertices_[1].position = point2 + offset;
        vertices_[2].position = point2 - offset;
        vertices_[3].position = point1 - offset;

        for (int i = 0; i<4; ++i)
            vertices_[i].color = color;
    }

    void draw(sf::RenderTarget &target, sf::RenderStates states) const
    {
        target.draw(vertices_, 4, sf::Quads);
    }

private:
    sf::Vertex vertices_[4];
    sf::Color color_;
};

class HandFrameListener : public astra::FrameListener
{
public:
    using PointList = std::deque<astra::Vector2i>;
    using PointMap = std::unordered_map<int, PointList>;

    HandFrameListener()
    {
        font_.loadFromFile("Inconsolata.otf");
        prev_ = ClockType::now();
    }

    HandFrameListener(const HandFrameListener&) = delete;
    HandFrameListener& operator=(const HandFrameListener&) = delete;

    void init_texture(int width, int height)
    {
        if (displayBuffer_ == nullptr || width != depthWidth_ || height != depthHeight_)
        {
            depthWidth_ = width;
            depthHeight_ = height;
            const int byteLength = depthWidth_ * depthHeight_ * 4;

            displayBuffer_ = BufferPtr(new uint8_t[byteLength]);
            std::fill(&displayBuffer_[0], &displayBuffer_[0] + byteLength, 0);

            texture_.create(depthWidth_, depthHeight_);
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

    void process_depth(astra::Frame& frame)
    {
        const astra::PointFrame pointFrame = frame.get<astra::PointFrame>();

        const int width = pointFrame.width();
        const int height = pointFrame.height();

        init_texture(width, height);

        visualizer_.update(pointFrame);
        const astra::RgbPixel* vizBuffer = visualizer_.get_output();

        for (int i = 0; i < width * height; i++)
        {
            const int rgbaOffset = i * 4;
            displayBuffer_[rgbaOffset] = vizBuffer[i].r;
            displayBuffer_[rgbaOffset + 1] = vizBuffer[i].b;
            displayBuffer_[rgbaOffset + 2] = vizBuffer[i].g;
            displayBuffer_[rgbaOffset + 3] = 255;
        }

        texture_.update(displayBuffer_.get());
    }

    void update_hand_trace(int trackingId, const astra::Vector2i& position)
    {
        auto it = pointMap_.find(trackingId);
        if (it == pointMap_.end())
        {
            PointList list;
            for (int i = 0; i < maxTraceLength_; i++)
            {
                list.push_back(position);
            }
            pointMap_.insert({trackingId, list});
        }
        else
        {
            PointList& list = it->second;
            while (list.size() < maxTraceLength_)
            {
                list.push_back(position);
            }
        }
    }

    void shorten_hand_traces()
    {
        auto it = pointMap_.begin();

        while (it != pointMap_.end())
        {
            PointList& list = it->second;
            if (list.size() > 1)
            {
                list.pop_front();
                ++it;
            }
            else
            {
                it = pointMap_.erase(it);
            }
        }
    }

    void process_hand_frame(astra::Frame& frame)
    {
        const astra::HandFrame handFrame = frame.get<astra::HandFrame>();

        handPoints_ = handFrame.handpoints();

        shorten_hand_traces();
        for (const auto& handPoint : handPoints_)
        {
            if (handPoint.status() == HAND_STATUS_TRACKING)
            {
                update_hand_trace(handPoint.tracking_id(), handPoint.depth_position());
            }
        }
    }

    virtual void on_frame_ready(astra::StreamReader& reader,
                                astra::Frame& frame) override
    {
        process_depth(frame);
        process_hand_frame(frame);

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

    void draw_shadow_text(sf::RenderWindow& window, sf::Text& text, sf::Color color, int x, int y)
    {
        text.setColor(sf::Color::Black);
        text.setPosition(x + 5, y + 5);
        window.draw(text);

        text.setColor(color);
        text.setPosition(x, y);

        window.draw(text);
    }

    void draw_hand_label(sf::RenderWindow& window,
                         const float radius,
                         const float x,
                         const float y,
                         const astra::HandPoint& handPoint)
    {
        const auto trackingId = handPoint.tracking_id();

        std::stringstream str;
        str << trackingId;
        if (handPoint.status() == HAND_STATUS_LOST) { str << "Lost"; }

        sf::Text label(str.str(), font_);
        const int characterSize = 60;
        label.setCharacterSize(characterSize);

        const auto bounds = label.getLocalBounds();
        label.setOrigin(bounds.left + bounds.width / 2.f, characterSize);
        draw_shadow_text(window, label, sf::Color::White, x, y - radius - 10);
    }

    void draw_hand_position(sf::RenderWindow& window,
                            const float radius,
                            const float x,
                            const float y,
                            const astra::HandPoint& handPoint)
    {
        const auto worldPosition = handPoint.world_position();

        std::stringstream str;
        str << std::fixed
            << std::setprecision(0)
            << worldPosition.x << ", "
            << worldPosition.y << ", "
            << worldPosition.z;

        sf::Text label(str.str(), font_);
        const int characterSize = 60;
        label.setCharacterSize(characterSize);

        const auto bounds = label.getLocalBounds();
        label.setOrigin(bounds.left + bounds.width / 2.f, 0);
        draw_shadow_text(window, label, sf::Color::White, x, y + radius + 10);
    }

    void draw_hand_trace(sf::RenderWindow& window,
                         const PointList& pointList,
                         const sf::Color& color,
                         const float depthScale)
    {
        if (pointList.size() < 2) { return; }

        const float thickness = 4;
        auto it = pointList.begin();
        astra::Vector2i previousPoint = *it;

        while (it != pointList.end())
        {
            const astra::Vector2i currentPoint = *it;
            ++it;

            const sf::Vector2f p1((previousPoint.x + .5f) * depthScale,
                                  (previousPoint.y + .5f) * depthScale);
            const sf::Vector2f p2((currentPoint.x + .5f) * depthScale,
                                  (currentPoint.y + .5f) * depthScale);
            previousPoint = currentPoint;

            window.draw(sfLine(p1, p2, color, thickness));
        }
    }

    void draw_hand_points(sf::RenderWindow& window, const float depthScale)
    {
        const float radius = 16;
        const sf::Color candidateColor(255, 255, 0);
        const sf::Color lostColor(255, 0, 0);
        const sf::Color trackingColor(0, 139, 69);

        for (const auto& handPoint : handPoints_)
        {
            sf::Color color;

            switch(handPoint.status())
            {
            case HAND_STATUS_LOST:
                color = lostColor;
            case HAND_STATUS_CANDIDATE:
                color = candidateColor;
            default:
                color = trackingColor;
            }

            const astra::Vector2i& p = handPoint.depth_position();

            const float circleX = (p.x + .5f) * depthScale;
            const float circleY = (p.y + .5f) * depthScale;

            draw_circle(window, radius, circleX, circleY, color);
            draw_hand_label(window, radius, circleX, circleY, handPoint);

            if (handPoint.status() == HAND_STATUS_TRACKING)
            {
                draw_hand_position(window, radius, circleX, circleY, handPoint);
            }
        }

        const sf::Color lineColor(0, 0, 255);
        for (const auto& kvp : pointMap_)
        {
            const PointList& list = kvp.second;
            draw_hand_trace(window, list, lineColor, depthScale);
        }
    }

    void draw_to(sf::RenderWindow& window)
    {
        if (displayBuffer_)
        {
            const float depthScale = window.getView().getSize().x / depthWidth_;
            sprite_.setScale(depthScale, depthScale);
            window.draw(sprite_);
            draw_hand_points(window, depthScale);
        }
    }

private:
    samples::common::LitDepthVisualizer visualizer_;

    using DurationType = std::chrono::milliseconds;
    using ClockType = std::chrono::high_resolution_clock;

    ClockType::time_point prev_;
    float elapsedMillis_{.0f};

    sf::Texture texture_;
    sf::Sprite sprite_;
    sf::Font font_;

    using BufferPtr = std::unique_ptr<uint8_t[]>;
    BufferPtr displayBuffer_{nullptr};

    std::vector<astra::HandPoint> handPoints_;

    PointMap pointMap_;

    int depthWidth_{0};
    int depthHeight_{0};
    int maxTraceLength_{15};
};

int main(int argc, char** argv)
{
    astra::initialize();

    set_key_handler();

    sf::RenderWindow window(sf::VideoMode(1280, 960), "Hand Viewer");

    astra::StreamSet streamSet;
    astra::StreamReader reader = streamSet.create_reader();

    reader.stream<astra::PointStream>().start();
    reader.stream<astra::HandStream>().start();

    HandFrameListener listener;

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

        listener.draw_to(window);
        window.display();

        if (!shouldContinue)
        {
            window.close();
        }
    }

    astra::terminate();

    return 0;
}
