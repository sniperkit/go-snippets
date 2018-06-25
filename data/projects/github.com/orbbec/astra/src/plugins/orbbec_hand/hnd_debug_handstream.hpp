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
#ifndef HND_DEBUG_HANDSTREAM_H
#define HND_DEBUG_HANDSTREAM_H

#include <astra_core/plugins/SingleBinStream.hpp>
#include <astra/capi/astra_ctypes.h>
#include <astra/capi/streams/stream_types.h>
#include <astra/Vector.hpp>

namespace astra { namespace hand {

    using debug_handview_type = astra_debug_hand_view_type_t;

    class debug_handstream : public plugins::single_bin_stream<astra_imageframe_wrapper_t>

    {
    public:
        debug_handstream(PluginServiceProxy& pluginService,
                         astra_streamset_t streamSet,
                         uint32_t width,
                         uint32_t height,
                         uint32_t bytesPerPixel)
            : single_bin_stream(pluginService,
                                streamSet,
                                StreamDescription(ASTRA_STREAM_DEBUG_HAND,
                                                  DEFAULT_SUBTYPE),
                                width * height * bytesPerPixel)
        { }

        debug_handview_type view_type() const { return viewType_; }
        void set_view_type(debug_handview_type view) { viewType_ = view; }

        bool use_mouse_probe() const { return useMouseProbe_; }
        const Vector2f& mouse_norm_position() const { return mouseNormPosition_; }
        bool pause_input() const { return pauseInput_; }
        bool spawn_point_locked() const { return lockSpawnPoint_; }
        const Vector2f& spawn_norm_position() const { return spawnNormPosition_; }

    protected:
        virtual void on_set_parameter(astra_streamconnection_t connection,
                                      astra_parameter_id id,
                                      size_t inByteLength,
                                      astra_parameter_data_t inData) override;

        virtual void on_get_parameter(astra_streamconnection_t connection,
                                      astra_parameter_id id,
                                      astra_parameter_bin_t& parameterBin) override;

        virtual void on_invoke(astra_streamconnection_t connection,
                               astra_command_id commandId,
                               size_t inByteLength,
                               astra_parameter_data_t inData,
                               astra_parameter_bin_t& parameterBin) override;
    private:
        void get_view_parameter(astra_parameter_bin_t& parameterBin);
        void set_view_parameter(size_t inByteLength,
                                astra_parameter_data_t& inData);
        void set_use_mouse_probe(size_t inByteLength, astra_parameter_data_t& inData);
        void set_mouse_norm_position(size_t inByteLength, astra_parameter_data_t& inData);
        void set_pause_input(size_t inByteLength, astra_parameter_data_t& inData);
        void set_lock_spawn_point(size_t inByteLength, astra_parameter_data_t& inData);

        debug_handview_type viewType_{ DEBUG_HAND_VIEW_DEPTH };
        bool useMouseProbe_{false};
        Vector2f mouseNormPosition_;
        Vector2f spawnNormPosition_;
        bool pauseInput_{false};
        bool lockSpawnPoint_{false};
    };
}}

#endif /* HND_DEBUG_HANDSTREAM_H */
