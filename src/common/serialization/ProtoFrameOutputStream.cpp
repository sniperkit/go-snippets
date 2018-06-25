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
#include "ProtoFrameOutputStream.h"

#include "pb_util.h"

namespace astra { namespace serialization {

    ProtoFrameOutputStream::ProtoFrameOutputStream(ZeroCopyOutputStream* outputStream):
        FrameOutputStream(),
            m_outputStream(outputStream)
    {

    }

    ProtoFrameOutputStream::~ProtoFrameOutputStream()
    {

    }

    void ProtoFrameOutputStream::stage_frame(Frame& frame)
    {
        proto::Frame frameMessage;
        populate_frame_message(frame.byteLength, frame.frameIndex, frame.rawFrameWrapper, frameMessage);

        m_frameMessage = frameMessage;

    }

    void ProtoFrameOutputStream::stage_frame_description(FrameDescription& frameDesc)
    {
        proto::FrameDescription frameDescMessage;
        populate_frame_description_message(frameDesc.framePeriod, frameDesc.bufferLength, frameDescMessage);

        m_frameDescriptionMessage = frameDescMessage;
    }

    void ProtoFrameOutputStream::stage_stream_header(StreamHeader& streamHeader)
    {
        proto::StreamHeader streamHeaderMessage;
        populate_stream_header_message(streamHeader.frameType, streamHeaderMessage);

        m_streamHeaderMessage = streamHeaderMessage;
    }

    bool ProtoFrameOutputStream::write_frame()
    {
        return proto::write_delimited_to(m_frameMessage, m_outputStream.get());
    }

    bool ProtoFrameOutputStream::write_frame_description()
    {
        return proto::write_delimited_to(m_frameDescriptionMessage, m_outputStream.get());
    }

    bool ProtoFrameOutputStream::write_stream_header()
    {
        return proto::write_delimited_to(m_streamHeaderMessage, m_outputStream.get());
    }

    void ProtoFrameOutputStream::populate_frame_message(size_t byteLength, int frameIndex, void* rawFrameWrapper, proto::Frame& message)
    {
        message.clear_rawframewrapper();

        message.set_bytelength(byteLength);
        message.set_frameindex(frameIndex);
        message.set_rawframewrapper(rawFrameWrapper, byteLength);
    }

    void ProtoFrameOutputStream::populate_frame_description_message(double framePeriod, int bufferLength, proto::FrameDescription& frameDescriptionMessage)
    {
        frameDescriptionMessage.set_frameperiod(framePeriod);
        frameDescriptionMessage.set_bufferlength(bufferLength);
    }

    void ProtoFrameOutputStream::populate_stream_header_message(int frameType, proto::StreamHeader& streamHeaderMessage)
    {
        streamHeaderMessage.set_frametype(frameType);
    }
}}