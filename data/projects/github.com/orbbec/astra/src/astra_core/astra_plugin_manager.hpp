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
#ifndef ASTRA_PLUGIN_MANAGER_H
#define ASTRA_PLUGIN_MANAGER_H

#include <string>
#include <vector>
#include <memory>
#include "astra_plugin_service.hpp"
#include "astra_shared_library.hpp"
#include <astra_core/capi/plugins/astra_pluginservice_proxy.h>
#include "astra_logger.hpp"

namespace astra {

    using initialize_fn = void(*)(astra_pluginservice_proxy_t*);
    using terminate_fn = void(*)();
    using update_fn = void(*)();

    struct PluginFuncs
    {
        initialize_fn initialize{nullptr};
        terminate_fn terminate{nullptr};
        update_fn update{nullptr};
        process::lib_handle libHandle{nullptr};

        bool is_valid()
        {
            return initialize != nullptr &&
                terminate != nullptr &&
                update != nullptr &&
                libHandle != nullptr;
        }
    };

    class plugin_manager
    {
    public:
        plugin_manager(streamset_catalog& setCatalog);
        ~plugin_manager();

        void load_plugins(std::string searchPath);
        void load_plugin(std::string pluginPath);

        void update();
        void unload_all_plugins();
        size_t plugin_count() const { return pluginList_.size(); }

        void notify_host_event(astra_event_id id, const void* data, size_t dataSize);
    private:
        std::vector<std::string> find_libraries(const std::string& pluginsPath);
        void try_load_plugin(const std::string& path);

        using PluginList = std::vector<PluginFuncs>;
        PluginList pluginList_;

        using plugin_service_ptr = std::unique_ptr<plugin_service>;
        plugin_service_ptr pluginService_;

        astra_pluginservice_proxy_t* pluginServiceProxy_{nullptr};
    };
}

#endif /* ASTRA_PLUGIN_MANAGER_H */
