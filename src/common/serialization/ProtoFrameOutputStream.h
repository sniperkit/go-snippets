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
#ifndef PROTOFRAMEOUTPUTSTREAM_H
#define PROTOFRAMEOUTPUTSTREAM_H

#include <common/serialization/FrameOutputStream.h>
#include "FrameDescription.pb.h"
#include "StreamHeader.pb.h"
#include "Frame.pb.h"

#include <google/protobuf/io/zero_copy_stream_impl.h>

#include <memory>

using namespace google::protobuf::io;

namespace google{namespace protobuf{namespace io{
    class ZeroCopyOutputStream;
}}}

namespace astra { namespace serialization {

    class ProtoFrameOutputStream : public FrameOutputStream
    {
    public:
        ProtoFrameOutputStream(ZeroCopyOutputStream* outputStream);
        virtual ~ProtoFrameOutputStream() override;

        void stage_frame(Frame& frame) override;
        void stage_frame_description(FrameDescription& frameDesc) override;
        void stage_stream_header(StreamHeader& streamHeader) override;
        bool write_frame() override;
        bool write_frame_description() override;
        bool write_stream_header() override;

    private:
        void populate_frame_message(size_t byteLength, int frameIndex, void* rawFrameWrapper, proto::Frame& message);
        void populate_frame_description_message(double framePeriod, int bufferLength, proto::FrameDescription& frameDescriptionMessage);
        void populate_stream_header_message(int frameType, proto::StreamHeader& streamHeaderMessage);

        std::unique_ptr<ZeroCopyOutputStream> m_outputStream;

        proto::Frame m_frameMessage;
        proto::FrameDescription m_frameDescriptionMessage;
        proto::StreamHeader m_streamHeaderMessage;
    };

}}

#endif /* PROTOFRAMEOUTPUTSTREAM_H */
