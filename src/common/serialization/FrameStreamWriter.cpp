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
#include <common/serialization/FrameStreamWriter.h>
#include "ProtoFrameOutputStream.h"

namespace astra { namespace serialization {

    FrameOutputStream* open_frame_output_stream(FILE* file)
    {
#ifdef _MSC_VER
        int fileDescriptor = _fileno(file);
#else
        int fileDescriptor = fileno(file);
#endif

        ZeroCopyOutputStream* outputStream = new FileOutputStream(fileDescriptor);

        return new ProtoFrameOutputStream(outputStream);
    }

    void close_frame_output_stream(FrameOutputStream*& stream)
    {
        delete stream;
        stream = nullptr;
    }

    FrameStreamWriter::FrameStreamWriter(FrameOutputStream& frameOutputStream):
        m_outputStream(frameOutputStream),
        m_swatchName("FrameStreamWriter")
    {
        m_swatch.set_mode(REAL_TIME);
    }

    FrameStreamWriter::~FrameStreamWriter()
    {

    }

    bool FrameStreamWriter::begin_write()
    {
        bool isSuccessful = true;

        StreamHeader streamHeader;
        streamHeader.frameType = 1;

        m_outputStream.stage_stream_header(streamHeader);

        isSuccessful = m_outputStream.write_stream_header();

        if (isSuccessful)
        {
            m_shouldWrite = true;
            m_swatch.start(m_swatchName);
        }

        return isSuccessful;
    }

    bool FrameStreamWriter::end_write()
    {
        m_shouldWrite = false;

        return true;
    }

    bool FrameStreamWriter::write(const DepthFrame& depthFrame)
    {
        bool isSuccessful = true;

        if (!m_shouldWrite)
        {
            return !isSuccessful;
        }

        m_swatch.stop(m_swatchName);

        astra_imageframe_t imageFrame = depthFrame.handle();
        astra_frame_t* astraFrame = imageFrame->frame;

        stage_frame(*astraFrame);
        stage_frame_description(*astraFrame, 1 / m_swatch.get_time_so_far(m_swatchName));

        isSuccessful = m_outputStream.write_frame_description();
        isSuccessful = isSuccessful && m_outputStream.write_frame();

        m_swatch.start(m_swatchName);

        return isSuccessful;
    }

    void FrameStreamWriter::stage_frame(astra_frame_t& astraFrame)
    {
        Frame frame;
        populate_frame(astraFrame, frame);
        m_outputStream.stage_frame(frame);
    }

    void FrameStreamWriter::stage_frame_description(astra_frame_t& astraFrame, double fps)
    {
        FrameDescription frameDesc;
        populate_frame_description(astraFrame, frameDesc, fps);
        m_outputStream.stage_frame_description(frameDesc);
    }

    void FrameStreamWriter::populate_frame(astra_frame_t& astraFrame, Frame& frame)
    {

        frame.frameIndex = astraFrame.frameIndex;
        frame.byteLength = astraFrame.byteLength;
        frame.rawFrameWrapper = astraFrame.data;
    }

    void FrameStreamWriter::populate_frame_description(astra_frame_t& astraFrame, FrameDescription& frameDescription, double fps)
    {
        frameDescription.bufferLength = astraFrame.byteLength;
        frameDescription.framePeriod = fps;
    }
}}
