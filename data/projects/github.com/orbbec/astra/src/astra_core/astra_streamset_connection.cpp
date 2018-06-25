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
#include "astra_streamset_connection.hpp"
#include <memory>
#include "astra_cxx_compatibility.hpp"

namespace astra {

    stream_reader* streamset_connection::create_reader()
    {
        ReaderPtr reader = astra::make_unique<stream_reader>(*this);
        stream_reader* rawPtr = reader.get();

        readers_.push_back(std::move(reader));
        return rawPtr;
    }

    bool streamset_connection::destroy_reader(astra::stream_reader* reader)
    {
        return true;
    }
}
