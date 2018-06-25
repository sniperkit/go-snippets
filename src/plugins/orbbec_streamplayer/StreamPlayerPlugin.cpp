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
#include "StreamPlayerPlugin.h"
#include <astra_core/astra_core.hpp>
#include "PlaybackStreamSet.h"

EXPORT_PLUGIN(astra::plugins::streamplayer::StreamPlayerPlugin)

namespace astra { namespace plugins { namespace streamplayer {

    StreamPlayerPlugin::StreamPlayerPlugin(PluginServiceProxy* pluginService)
        : plugin_base(pluginService, "orbbec_streamplayer")
    {
        create_streamset();
    }

    StreamPlayerPlugin::~StreamPlayerPlugin()
    {
        m_sets.clear();
    }

    void StreamPlayerPlugin::create_streamset()
    {
        PlaybackStreamSetBase* set = new PlaybackStreamSet(pluginService(), "stream_player");
        set->open();

        m_sets.push_back(SetPtr(set));
    }

    void StreamPlayerPlugin::temp_update()
    {
        read_streams();
    }

    astra_status_t StreamPlayerPlugin::read_streams()
    {
        for (auto& set : m_sets)
        {
            set->read();
        }

        return ASTRA_STATUS_SUCCESS;
    }

}}}
