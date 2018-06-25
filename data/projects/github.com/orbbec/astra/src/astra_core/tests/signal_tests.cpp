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
#define CATCH_CONFIG_MAIN
#include "catch.hpp"
#include "../astra_signal.hpp"

TEST_CASE("Can add to a callback list", "[signal]") {
    astra::CallbackList<void, int> cbList;
    size_t id = cbList.add([] (int v) { /* do nothing */ });
    id = cbList.add([] (int v) { /* do nothing */ });
    cbList.add([] (int v) { /* do nothing */ });

    cbList.remove(id);
    REQUIRE(cbList.debug_count() == 2);
}

TEST_CASE("Can add a slot to a signal", "[signal]") {
    astra::signal<std::string> signal;
    size_t id = signal += [] (std::string v) { };

    REQUIRE(id != size_t(NULL));
}

TEST_CASE("Can raise signal", "[signal]") {
    astra::signal<int> signal;
    int test = 0;
    signal += [&test] (int v) { test++; };
    signal.raise(1);

    REQUIRE(test == 1);
}

TEST_CASE("Can raise two signals", "[signal]") {
    astra::signal<int> signal;
    int test = 0;
    signal += [&test] (int v) { test++; };
    signal += [&test] (int v) { test+=2; };
    signal.raise(1);
    REQUIRE(test == 3);
}

TEST_CASE("Can raise void signal", "[signal]") {
    astra::signal<void> signal;
    int test = 0;
    signal += [&test] () { test++; };

    signal.raise();
    REQUIRE(test == 1);
}

TEST_CASE("Can raise multi param signal", "[signal]") {
    astra::signal<bool, int, std::string> signal;
    int test = 1;
    signal += [&test] (bool one, int two, std::string three) { test+=two; };
    size_t id = signal += [&test] (bool one, int two, std::string three) { test+=5; };
    signal += [&test] (bool one, int two, std::string three) { test+=30; };
    signal -= id;

    signal.raise(false, 5, "Test");
    REQUIRE(test == 36);
}
