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
#include "mock_color_generator.hpp"
#include <cassert>

namespace orbbec { namespace mocks {

    color_generator::color_generator(std::uint32_t width, std::uint32_t height)
        : frame_generator(width, height)
    { }

    color_generator::~color_generator() = default;

    void color_generator::on_generate()
    {
        target().clear(sf::Color(128, 0, 0, 255));
    }
}}
