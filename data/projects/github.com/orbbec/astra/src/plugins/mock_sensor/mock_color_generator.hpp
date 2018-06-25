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
#ifndef MOCK_COLOR_GENERATOR_H
#define MOCK_COLOR_GENERATOR_H

#include <SFML/Graphics.hpp>
#include "mock_frame_generator.hpp"

namespace orbbec { namespace mocks {

    class color_generator : public frame_generator
    {
    public:
        color_generator(std::uint32_t width, std::uint32_t height);
        virtual ~color_generator();

    private:
        virtual void on_generate() override;
    };
}}

#endif /* MOCK_COLOR_GENERATOR_H */
