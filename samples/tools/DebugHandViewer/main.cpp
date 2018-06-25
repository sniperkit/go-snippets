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
#include <sstream>
#include <iomanip>
#include <cstring>

class debug_frame_listener : public astra::FrameListener
{
public:
    debug_frame_listener()
    {
        font_.loadFromFile("Inconsolata.otf");
    }

    debug_frame_listener(const debug_frame_listener&) = delete;
    debug_frame_listener& operator=(const debug_frame_listener&) = delete;

    void init_texture(int width, int height)
    {
        if (displayBuffer_ == nullptr || width != displayWidth_ || height != displayHeight_)
        {
            displayWidth_ = width;
            displayHeight_ = height;
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
        double fpsFactor = 0.02;

        std::clock_t newTimepoint = std::clock();
        long double frameDuration = (newTimepoint - lastTimepoint_) / static_cast<long double>(CLOCKS_PER_SEC);

        frameDuration_ = frameDuration * fpsFactor + frameDuration_ * (1 - fpsFactor);
        lastTimepoint_ = newTimepoint;
        double fps = 1.0 / frameDuration_;

        if (outputFPS_)
        {
            printf("FPS: %3.1f (%3.4Lf ms)\n", fps, frameDuration_ * 1000);
        }
    }

    void processDepth(astra::Frame& frame)
    {
        const astra::DepthFrame depthFrame = frame.get<astra::DepthFrame>();

        int width = depthFrame.width();
        int height = depthFrame.height();
        depthWidth_ = width;
        depthHeight_ = height;
        /*
        init_texture(width, height);

        const int16_t* depthPtr = depthFrame.data();
        for(int y = 0; y < height; y++)
        {
        for(int x = 0; x < width; x++)
        {
        int index = (x + y * width);
        int16_t depth = depthPtr[index];
        uint8_t value = depth % 255;
        displayBuffer_[index * 4] = value;
        displayBuffer_[index * 4 + 1] = value;
        displayBuffer_[index * 4 + 2] = value;
        displayBuffer_[index * 4 + 3] = 255;
        }
        }

        texture_.update(displayBuffer_.get());
        */
    }

    void processHandFrame(astra::Frame& frame)
    {
        const astra::HandFrame handFrame = frame.get<astra::HandFrame>();

        handPoints_ = handFrame.handpoints();
    }

    void processdebug_handframe(astra::Frame& frame)
    {
        const astra::DebugHandFrame handFrame = frame.get<astra::DebugHandFrame>();

        int width = handFrame.width();
        int height = handFrame.height();

        init_texture(width, height);

        const astra::RgbPixel* imagePtr = handFrame.data();
        for (int y = 0; y < height; y++)
        {
            for (int x = 0; x < width; x++)
            {
                int index = (x + y * width);
                int index4 = index * 4;
                astra::RgbPixel rgb = imagePtr[index];

                displayBuffer_[index4] = rgb.r;
                displayBuffer_[index4 + 1] = rgb.g;
                displayBuffer_[index4 + 2] = rgb.b;
                displayBuffer_[index4 + 3] = 255;
            }
        }

        texture_.update(displayBuffer_.get());
    }

    virtual void on_frame_ready(astra::StreamReader& reader,
                                astra::Frame& frame) override
    {
        processDepth(frame);
        processHandFrame(frame);
        processdebug_handframe(frame);

        viewType_ = reader.stream<astra::DebugHandStream>().get_view_type();

        check_fps();
    }

    void drawCircle(sf::RenderWindow& window, float radius, float x, float y, sf::Color color)
    {
        sf::CircleShape shape(radius);

        shape.setFillColor(color);

        shape.setOrigin(radius, radius);
        shape.setPosition(x, y);
        window.draw(shape);
    }

    void drawShadowText(sf::RenderWindow& window, sf::Text& text, sf::Color color, int x, int y)
    {
        text.setColor(sf::Color::Black);
        text.setPosition(x + 5, y + 5);
        window.draw(text);

        text.setColor(color);
        text.setPosition(x, y);
        window.draw(text);
    }

    void drawHandLabel(sf::RenderWindow& window, float radius, float x, float y, astra::HandPoint& handPoint)
    {
        int32_t trackingId = handPoint.tracking_id();
        std::stringstream str;
        str << trackingId;
        if (handPoint.status() == HAND_STATUS_LOST)
        {
            str << " Lost";
        }
        sf::Text label(str.str(), font_);
        int characterSize = 60;
        label.setCharacterSize(characterSize);

        auto bounds = label.getLocalBounds();
        label.setOrigin(bounds.left + bounds.width / 2.0, characterSize);
        drawShadowText(window, label, sf::Color::White, x, y - radius - 10);
    }

    void drawHandPosition(sf::RenderWindow& window, float radius, float x, float y, astra::HandPoint& handPoint)
    {
        auto worldPosition = handPoint.world_position();
        std::stringstream str;
        str << std::fixed << std::setprecision(0);
        str << worldPosition.x << "," << worldPosition.y << "," << worldPosition.z;
        sf::Text label(str.str(), font_);
        int characterSize = 60;
        label.setCharacterSize(characterSize);

        auto bounds = label.getLocalBounds();
        label.setOrigin(bounds.left + bounds.width / 2.0, 0);
        drawShadowText(window, label, sf::Color::White, x, y + radius + 10);
    }

    void drawhandpoints(sf::RenderWindow& window, float depthScale)
    {
        sf::Color candidateColor(255, 255, 0);
        sf::Color lostColor(255, 0, 0);
        sf::Color trackingColor(128, 138, 0);

        for (auto handPoint : handPoints_)
        {
            sf::Color color = trackingColor;
            astra_handstatus_t status = handPoint.status();
            if (status == HAND_STATUS_LOST)
            {
                color = lostColor;
            }
            else if (status == HAND_STATUS_CANDIDATE)
            {
                color = candidateColor;
            }

            const astra::Vector2i& p = handPoint.depth_position();

            float circleX = (p.x + 0.5) * depthScale;
            float circleY = (p.y + 0.5) * depthScale;
            if (showCircles_)
            {
                drawCircle(window, circleRadius_, circleX, circleY, color);
            }

            drawHandLabel(window, circleRadius_, circleX, circleY, handPoint);
            if (status == HAND_STATUS_TRACKING)
            {
                drawHandPosition(window, circleRadius_, circleX, circleY, handPoint);
            }
        }
    }

    std::string getViewName(astra::DebugHandViewType viewType)
    {
        switch (viewType)
        {
        case DEBUG_HAND_VIEW_DEPTH:
            return "Depth";
            break;
        case DEBUG_HAND_VIEW_DEPTH_MOD:
            return "Depth mod";
            break;
        case DEBUG_HAND_VIEW_DEPTH_AVG:
            return "Depth avg";
            break;
        case DEBUG_HAND_VIEW_VELOCITY:
            return "Velocity";
            break;
        case DEBUG_HAND_VIEW_FILTEREDVELOCITY:
            return "Filtered velocity";
            break;
        case DEBUG_HAND_VIEW_UPDATE_SEGMENTATION:
            return "Update segmentation";
            break;
        case DEBUG_HAND_VIEW_CREATE_SEGMENTATION:
            return "Create segmentation";
            break;
        case DEBUG_HAND_VIEW_UPDATE_SEARCHED:
            return "Update searched";
            break;
        case DEBUG_HAND_VIEW_CREATE_SEARCHED:
            return "Create searched";
            break;
        case DEBUG_HAND_VIEW_CREATE_SCORE:
            return "Create score";
            break;
        case DEBUG_HAND_VIEW_UPDATE_SCORE:
            return "Update score";
            break;
        case DEBUG_HAND_VIEW_HANDWINDOW:
            return "Hand window";
            break;
        case DEBUG_HAND_VIEW_TEST_PASS_MAP:
            return "Test pass map";
            break;
        default:
            return "Unknown view";
            break;
        }
    }

    void drawDebugViewName(sf::RenderWindow& window)
    {
        std::string viewName = getViewName(viewType_);
        sf::Text text(viewName, font_);
        int characterSize = 60;
        text.setCharacterSize(characterSize);
        text.setStyle(sf::Text::Bold);

        int x = 10;
        int y = window.getView().getSize().y - 20 - characterSize;

        drawShadowText(window, text, sf::Color::White, x, y);
    }

    void drawTo(sf::RenderWindow& window)
    {
        if (displayBuffer_ != nullptr)
        {
            float debugScale = window.getView().getSize().x / displayWidth_;
            float depthScale = window.getView().getSize().x / depthWidth_;

            sprite_.setScale(debugScale, debugScale);

            window.draw(sprite_);

            drawhandpoints(window, depthScale);

            drawDebugViewName(window);
        }
    }

    void toggle_output_fps()
    {
        outputFPS_ = !outputFPS_;
    }

    void toggle_circles()
    {
        showCircles_ = !showCircles_;
    }

private:
    long double frameDuration_{ 0 };
    std::clock_t lastTimepoint_{ 0 };
    sf::Texture texture_;
    sf::Sprite sprite_;

    sf::Font font_;

    using BufferPtr = std::unique_ptr < uint8_t[] > ;
    BufferPtr displayBuffer_;

    std::vector<astra::HandPoint> handPoints_;
    astra::DebugHandViewType viewType_;
    int depthWidth_{ 1 };
    int depthHeight_{ 1 };
    int displayWidth_{ 1 };
    int displayHeight_{ 1 };
    bool outputFPS_{ false };
    float circleRadius_ { 16 };
    bool showCircles_ { true };
};

void request_view_mode(astra::StreamReader& reader, astra::DebugHandViewType view)
{
    reader.stream<astra::DebugHandStream>().set_view_type(view);
}

void process_mouse_move(sf::RenderWindow& window, astra::StreamReader& reader)
{
    sf::Vector2i position = sf::Mouse::getPosition(window);
    astra::Vector2f normPosition;
    auto windowSize = window.getSize();
    normPosition.x = position.x / (float)windowSize.x;
    normPosition.y = position.y / (float)windowSize.y;

    reader.stream<astra::DebugHandStream>().set_mouse_position(normPosition);
}

static bool g_mouseProbe = false;
static bool g_pauseInput = false;
static bool g_lockSpawnPoint = false;

void toggle_mouse_probe(astra::StreamReader& reader)
{
    g_mouseProbe = !g_mouseProbe;
    reader.stream<astra::DebugHandStream>().set_use_mouse_probe(g_mouseProbe);
}

void toggle_pause_input(astra::StreamReader& reader)
{
    g_pauseInput = !g_pauseInput;
    reader.stream<astra::DebugHandStream>().set_pause_input(g_pauseInput);
}

void toggle_spawn_lock(astra::StreamReader& reader)
{
    g_lockSpawnPoint = !g_lockSpawnPoint;
    reader.stream<astra::DebugHandStream>().set_lock_spawn_point(g_lockSpawnPoint);
}

void process_key_input(astra::StreamReader& reader, debug_frame_listener& listener, sf::Event::KeyEvent key)
{
    if (key.code == sf::Keyboard::F)
    {
        listener.toggle_output_fps();
    }
    if (key.code == sf::Keyboard::Tilde ||
        key.code == sf::Keyboard::Dash)
    {
        request_view_mode(reader, DEBUG_HAND_VIEW_DEPTH);
    }
    else if (key.code == sf::Keyboard::Num1)
    {
        request_view_mode(reader, DEBUG_HAND_VIEW_DEPTH_MOD);
    }
    else if (key.code == sf::Keyboard::Num2)
    {
        request_view_mode(reader, DEBUG_HAND_VIEW_VELOCITY);
    }
    else if (key.code == sf::Keyboard::Num3)
    {
        request_view_mode(reader, DEBUG_HAND_VIEW_FILTEREDVELOCITY);
    }
    else if (key.code == sf::Keyboard::Num4)
    {
        request_view_mode(reader, DEBUG_HAND_VIEW_UPDATE_SCORE);
    }
    else if (key.code == sf::Keyboard::Num5)
    {
        request_view_mode(reader, DEBUG_HAND_VIEW_UPDATE_SEGMENTATION);
    }
    else if (key.code == sf::Keyboard::Num6)
    {
        request_view_mode(reader, DEBUG_HAND_VIEW_UPDATE_SEARCHED);
    }
    else if (key.code == sf::Keyboard::Num7)
    {
        request_view_mode(reader, DEBUG_HAND_VIEW_CREATE_SCORE);
    }
    else if (key.code == sf::Keyboard::Num8)
    {
        request_view_mode(reader, DEBUG_HAND_VIEW_CREATE_SEGMENTATION);
    }
    else if (key.code == sf::Keyboard::Num9)
    {
        request_view_mode(reader, DEBUG_HAND_VIEW_CREATE_SEARCHED);
    }
    else if (key.code == sf::Keyboard::Num0)
    {
        request_view_mode(reader, DEBUG_HAND_VIEW_HANDWINDOW);
    }
    else if (key.code == sf::Keyboard::W)
    {
        request_view_mode(reader, DEBUG_HAND_VIEW_DEPTH_AVG);
    }
    else if (key.code == sf::Keyboard::E)
    {
        request_view_mode(reader, DEBUG_HAND_VIEW_TEST_PASS_MAP);
    }
    else if (key.code == sf::Keyboard::M)
    {
        toggle_mouse_probe(reader);
    }
    else if (key.code == sf::Keyboard::P)
    {
        toggle_pause_input(reader);
    }
    else if (key.code == sf::Keyboard::C)
    {
        listener.toggle_circles();
    }
}

int main(int argc, char** argv)
{
    astra::initialize();

    sf::RenderWindow window(sf::VideoMode(1280, 960), "Hand Debug Viewer");

    astra::StreamSet streamSet;
    astra::StreamReader reader = streamSet.create_reader();

    debug_frame_listener listener;

    reader.stream<astra::DepthStream>().start();
    auto handStream = reader.stream<astra::HandStream>();
    handStream.start();
    handStream.set_include_candidate_points(true);

    reader.stream<astra::DebugHandStream>().start();
    reader.add_listener(listener);

    while (window.isOpen())
    {
        astra_temp_update();

        sf::Event event;
        while (window.pollEvent(event))
        {
            if (event.type == sf::Event::Closed)
            {
                window.close();
            }
            if (event.type == sf::Event::KeyPressed)
            {
                if (event.key.code == sf::Keyboard::Escape ||
                    (event.key.code == sf::Keyboard::C && event.key.control))
                {
                    window.close();
                }
                else
                {
                    process_key_input(reader, listener, event.key);
                }
            }
            else if (event.type == sf::Event::MouseMoved)
            {
                process_mouse_move(window, reader);
            }
            else if (event.type == sf::Event::MouseButtonPressed)
            {
                if (event.mouseButton.button == sf::Mouse::Right)
                {
                    toggle_pause_input(reader);
                }
                else
                {
                    process_mouse_move(window, reader);
                    toggle_spawn_lock(reader);
                }
            }
        }

        // clear the window with black color
        window.clear(sf::Color::Black);

        listener.drawTo(window);
        window.display();
    }

    astra::terminate();

    return 0;
}
