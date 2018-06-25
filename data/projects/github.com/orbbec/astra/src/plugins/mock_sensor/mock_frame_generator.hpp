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
#ifndef MOCK_FRAME_GENERATOR_H
#define MOCK_FRAME_GENERATOR_H

#include <SFML/Graphics.hpp>
#include <cstdint>

namespace orbbec { namespace mocks {

    class frame_generator
    {
    public:
        frame_generator(std::uint32_t width, std::uint32_t height);
        virtual ~frame_generator();

        frame_generator(const frame_generator&) = delete;
        const frame_generator& operator=(const frame_generator&) = delete;

        const std::uint8_t* pixels() const;
        void generate();

        std::uint32_t width() const { return width_; }
        std::uint32_t height() const { return height_; }

        void set_size(std::uint32_t width, std::uint32_t height);

        const std::string& overlayText() { return overlayText_; }
        void set_overlayText(const std::string& overlayText)
        {
            overlayText_ = overlayText;

            if (overlayText_.size() > 0)
            {
                text_.setString(overlayText_);
            }
        }

        void set_overlayColor(std::uint8_t r, std::uint8_t g, std::uint8_t b, std::uint8_t a)
        {
            sf::Color color;
            color.r = r;
            color.g = g;
            color.b = b;
            color.a = a;

            text_.setColor(color);
        }

    protected:
        sf::RenderTexture& target() { return target_; }

    private:
        virtual void on_generate() = 0;

        std::uint32_t width_{0};
        std::uint32_t height_{0};
        std::string overlayText_;

        sf::RenderTexture target_;
        sf::Image image_;
        sf::Font font_;
        sf::Text text_;
    };
}}
#endif /* MOCK_FRAME_GENERATOR_H */
