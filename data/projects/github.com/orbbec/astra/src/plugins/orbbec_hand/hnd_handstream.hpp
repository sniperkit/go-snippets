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
#ifndef HND_HAND_STREAM_H
#define HND_HAND_STREAM_H

#include <astra_core/plugins/SingleBinStream.hpp>
#include <astra/capi/streams/hand_types.h>
#include <astra/capi/astra_ctypes.h>
#include <astra/capi/streams/stream_types.h>
#include <Shiny.h>

namespace astra { namespace hand {

    class handstream : public plugins::single_bin_stream<astra_handframe_wrapper_t>
    {
    public:
        handstream(PluginServiceProxy& pluginService,
                   astra_streamset_t streamSet,
                   size_t maxHandCount)
            : single_bin_stream(pluginService,
                                streamSet,
                                StreamDescription(ASTRA_STREAM_HAND,
                                                  DEFAULT_SUBTYPE),
                                sizeof(astra_handpoint_t) * maxHandCount)
        { }

        bool include_candidate_points() const { return includeCandidatePoints_; }
        void set_include_candidate_points(bool includeCandidatePoints)
        {
            includeCandidatePoints_ = includeCandidatePoints;
        }
    protected:
        virtual void on_set_parameter(astra_streamconnection_t connection,
                                      astra_parameter_id id,
                                      size_t inByteLength,
                                      astra_parameter_data_t inData) override;

        virtual void on_get_parameter(astra_streamconnection_t connection,
                                      astra_parameter_id id,
                                      astra_parameter_bin_t& parameterBin) override;
    private:
        void get_include_candidates(astra_parameter_bin_t& parameterBin);
        void set_include_candidates(size_t inByteLength, astra_parameter_data_t& inData);

        virtual void on_connection_removed(astra_bin_t bin,
                                           astra_streamconnection_t connection) override
        {
            single_bin_stream::on_connection_removed(bin, connection);

#ifdef __ANDROID__
            PROFILE_UPDATE();
            PROFILE_OUTPUT("/sdcard/profile_orbbec_hand.txt");
#endif
        }

        bool includeCandidatePoints_{false};
    };

}}

#endif /* HND_HAND_STREAM_H */
