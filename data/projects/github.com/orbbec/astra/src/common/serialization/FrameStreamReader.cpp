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
#include <common/serialization/FrameStreamReader.h>

#include <memory>

#include "ProtoFrameInputStream.h"

namespace astra { namespace serialization {

    FrameStreamReader::FrameStreamReader(FrameInputStream* frameStream) :
        m_inputStream(frameStream)
    {
        if (m_inputStream->get_position() != 0)
        {
            m_inputStream->seek(-(m_inputStream->get_position()));
        }

        m_inputStream->read_stream_header(m_streamHeader);
        m_inputStream->read_frame_description(m_frameDescription);

        m_pulser.set_period(1 / m_frameDescription->framePeriod);

        m_pulser.start();
    }

    FrameStreamReader::~FrameStreamReader()
    {

    }

    void FrameStreamReader::close()
    {
        m_inputStream->close();
    }

    bool FrameStreamReader::read()
    {
        bool isSuccessful = true;

        if (m_isEndOfFile)
        {
            m_inputStream->seek_to_first_frame();

            isSuccessful = m_inputStream->read_frame_description(m_frameDescription);
            if (isSuccessful)
            {
                m_pulser.set_period(1 / m_frameDescription->framePeriod);
            }

            m_isEndOfFile = false;
        }

        m_pulser.stop();

        if (!m_pulser.is_pulse())
        {
            return !isSuccessful;
        }

        isSuccessful = m_inputStream->read_frame(m_frame);

        if (!is_end_of_file())
        {
            isSuccessful = m_inputStream->read_frame_description(m_frameDescription);

            if (isSuccessful)
            {
                m_pulser.set_period(1 / m_frameDescription->framePeriod);
            }
        }

        m_pulser.start();

        return isSuccessful;
    }

    bool FrameStreamReader::seek(int numberOfFrames)
    {
        int sizeMarkerOffset = 8;
        int totalFrameSize = m_inputStream->get_frame_description_size();

        int offset = numberOfFrames * (totalFrameSize + sizeMarkerOffset);
        return m_inputStream->seek(offset);
    }

    int FrameStreamReader::get_stream_type()
    {
        return m_streamHeader->frameType;
    }

    int FrameStreamReader::get_buffer_length()
    {
        return m_frameDescription->bufferLength;
    }

    bool FrameStreamReader::is_end_of_file()
    {
        m_isEndOfFile = m_inputStream->is_end_of_file();
        return m_isEndOfFile;
    }

    Frame& FrameStreamReader::peek()
    {
        return *m_frame;
    }

    FrameInputStream* FrameStreamReader::get_frame_input_stream()
    {
        return m_inputStream;
    }

    FrameInputStream* open_frame_input_stream(const char* path)
    {
        return new ProtoFrameInputStream(path);
    }

}}