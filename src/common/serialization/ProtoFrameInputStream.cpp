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
#include "ProtoFrameInputStream.h"

#include <astra_core/astra_cxx_make_unique.hpp>
#include "FrameDescription.pb.h"
#include "Frame.pb.h"

#include "pb_util.h"

namespace astra { namespace serialization {

    ProtoFrameInputStream::ProtoFrameInputStream(const char* path) :
        FrameInputStream()
    {
        m_file = fopen(path, "rb");

        if (m_file == nullptr)
        {
            throw ResourceNotFoundException(path);
        }

        m_fileDescriptor = get_file_descriptor(m_file);

        m_fileSize = get_file_size(m_fileDescriptor);

        m_inputStream = astra::make_unique<FileInputStream>(m_fileDescriptor);
    }

    int ProtoFrameInputStream::get_file_descriptor(FILE* file)
    {
#ifdef _MSC_VER
        int fileDescriptor = _fileno(m_file);
#else
        int fileDescriptor = fileno(m_file);
#endif

        return fileDescriptor;
    }

    long ProtoFrameInputStream::get_file_size(int fd)
    {
        struct stat stat_buf;
        int rc = fstat(fd, &stat_buf);
        return rc == 0 ? stat_buf.st_size : -1;
    }

    ProtoFrameInputStream::~ProtoFrameInputStream()
    {
        close();
    }

    void ProtoFrameInputStream::close()
    {
        fclose(m_file);
    }

    bool ProtoFrameInputStream::read_stream_header(StreamHeader*& streamHeader)
    {
        long filePos = ftell(m_file);

        if (filePos != 0)
        {
            fseek(m_file, 0, SEEK_SET);
            m_inputStream = astra::make_unique<FileInputStream>(m_fileDescriptor);
        }

        bool isSuccessful = proto::read_delimited_to(m_inputStream.get(), &m_streamHeaderMessage);

        m_streamHeader.frameType = m_streamHeaderMessage.frametype();

        if (isSuccessful)
        {
            streamHeader = &m_streamHeader;
        }
        else
        {
            streamHeader = nullptr;
        }

        return isSuccessful;
    }

    bool ProtoFrameInputStream::read_frame(Frame*& frame)
    {
        bool isSuccessful = proto::read_delimited_to(m_inputStream.get(), &m_frameMessage);

        m_frame.byteLength = m_frameMessage.bytelength();
        m_frame.frameIndex = m_frameMessage.frameindex();

        std::string* frameWrapperString = m_frameMessage.mutable_rawframewrapper();
        m_frame.rawFrameWrapper = &(*frameWrapperString)[0];

        if (isSuccessful)
        {
            frame = &m_frame;
        }
        else
        {
            frame = nullptr;
        }

        return isSuccessful;
    }

    bool ProtoFrameInputStream::read_frame_description(FrameDescription*& frameDescription)
    {
        bool isSuccessful = proto::read_delimited_to(m_inputStream.get(), &m_frameDescriptionMessage);

        m_frameDescription.framePeriod = m_frameDescriptionMessage.frameperiod();
        m_frameDescription.bufferLength = m_frameDescriptionMessage.bufferlength();

        if (isSuccessful)
        {
            frameDescription = &m_frameDescription;
        }
        else
        {
            frameDescription = nullptr;
        }

        return isSuccessful;
    }

    bool ProtoFrameInputStream::seek_to_first_frame()
    {
        bool isSuccessful = true;

        int offset = get_stream_header_size() + 4;

        fseek(m_file, offset, SEEK_SET);
        m_inputStream = astra::make_unique<FileInputStream>(m_fileDescriptor);
        m_positionOffset = ftell(m_file);

        return isSuccessful;
    }


    bool ProtoFrameInputStream::seek(int offset)
    {
        bool isSuccessful = true;

        if (offset < 0 && get_position() + offset < 0)
        {
            return !isSuccessful;
        }
        if (get_position() + offset > m_fileSize)
        {
            return !isSuccessful;
        }

        fseek(m_file, offset, SEEK_CUR);
        m_inputStream = astra::make_unique<FileInputStream>(m_fileDescriptor);
        m_positionOffset = ftell(m_file);

        return isSuccessful;
    }

    int64_t ProtoFrameInputStream::get_position()
    {
        return m_inputStream->ByteCount() + m_positionOffset;
    }

    bool ProtoFrameInputStream::is_end_of_file()
    {
        return get_position() >= m_fileSize;
    }

    int ProtoFrameInputStream::get_frame_description_size()
    {
        return m_frameDescriptionMessage.ByteSize();
    }

    int ProtoFrameInputStream::get_stream_header_size()
    {
        return m_streamHeaderMessage.ByteSize();
    }
}}
