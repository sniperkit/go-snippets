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
#ifndef PROTOFRAMEINPUTSTREAM_H
#define PROTOFRAMEINPUTSTREAM_H

#include <common/serialization/FrameInputStream.h>

#include "Frame.pb.h"
#include "FrameDescription.pb.h"
#include "StreamHeader.pb.h"

#include <google/protobuf/io/zero_copy_stream_impl.h>

#include <cstdint>
#include <memory>
#include <sys/stat.h>

using namespace google::protobuf::io;

namespace astra { namespace serialization {

    class ProtoFrameInputStream : public FrameInputStream
    {
    public:
        ProtoFrameInputStream(const char* path);

        virtual ~ProtoFrameInputStream();

        void close() override final;
        bool read_stream_header(StreamHeader*& streamHeader) override;
        bool read_frame(Frame*& frame) override;
        bool read_frame_description(FrameDescription*& frameDescription) override;
        bool seek(int offset) override;
        bool seek_to_first_frame() override;
        int64_t get_position() override;
        bool is_end_of_file() override;
        int get_frame_description_size() override;
        int get_stream_header_size() override;

    private:
        long get_file_size(int fd);
        int get_file_descriptor(FILE* file);

        FILE* m_file;
        int m_fileDescriptor{-1};

        int m_positionOffset{0};

        std::unique_ptr<ZeroCopyInputStream> m_inputStream;
        proto::Frame m_frameMessage;
        proto::FrameDescription m_frameDescriptionMessage;
        proto::StreamHeader m_streamHeaderMessage;

        Frame m_frame;
        FrameDescription m_frameDescription;
        StreamHeader m_streamHeader;
        long m_fileSize{-1};
    };

}}

#endif /* PROTOFRAMEINPUTSTREAM_H */
