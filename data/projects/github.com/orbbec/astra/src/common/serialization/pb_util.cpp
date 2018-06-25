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
#include "pb_util.h"

#include "google/protobuf/message_lite.h"
#include "google/protobuf/io/coded_stream.h"
#include "google/protobuf/io/zero_copy_stream_impl.h"
#include <stdint.h>

using namespace google::protobuf;
using namespace google::protobuf::io;

// http://stackoverflow.com/questions/22881876/protocol-buffers-how-to-serialize-and-deserialize-multiple-messages-into-a-file

namespace astra { namespace serialization { namespace proto {

    bool write_delimited_to( const MessageLite& message, ZeroCopyOutputStream* rawOutput)
    {
        // We create a new coded stream for each message.  Don't worry, this is fast.
        CodedOutputStream output(rawOutput);

        // Write the size.
        const int size = message.ByteSize();
        output.WriteLittleEndian32(size);

        uint8_t* buffer = output.GetDirectBufferForNBytesAndAdvance(size);
        if (buffer != NULL) {
            // Optimization:  The message fits in one buffer, so use the faster
            // direct-to-array serialization path.
            message.SerializeWithCachedSizesToArray(buffer);
        }
        else {
            // Slightly-slower path when the message is multiple buffers.
            message.SerializeWithCachedSizes(&output);
            if (output.HadError()) return false;
        }

        return true;
    }

    bool read_delimited_to(ZeroCopyInputStream* rawInput, MessageLite* message)
    {
        // We create a new coded stream for each message.  Don't worry, this is fast,
        // and it makes sure the 64MB total size limit is imposed per-message rather
        // than on the whole stream.  (See the CodedInputStream interface for more
        // info on this limit.)
        CodedInputStream input(rawInput);

        // Read the size.
        uint32_t size;
        if (!input.ReadLittleEndian32(&size)) return false;

        // Tell the stream not to read beyond that size.
        CodedInputStream::Limit limit =
            input.PushLimit(size);

        // Parse the message.
        if (!message->MergeFromCodedStream(&input)) return false;
        if (!input.ConsumedEntireMessage()) return false;

        // Release the limit.
        input.PopLimit(limit);

        return true;
    }

}}}
