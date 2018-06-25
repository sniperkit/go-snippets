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
#include <astra/capi/astra.h>
#include <stdio.h>
#include <key_handler.h>

void print_hand_frame(astra_handframe_t handFrame)
{
    astra_handpoint_t* handPoints;
    size_t handCount;

    astra_handframe_get_hand_count(handFrame, &handCount);

    handPoints = (astra_handpoint_t*)malloc(handCount * sizeof(astra_handpoint_t));
    astra_handframe_copy_hands(handFrame, handPoints);

    astra_frame_index_t frameIndex;
    astra_handframe_get_frameindex(handFrame, &frameIndex);

    int activeHandCount = 0;

    astra_handpoint_t* point = handPoints;

    int i;
    for (i = 0; i < handCount; i++)
    {
        if (point->status == HAND_STATUS_TRACKING)
        {
            ++activeHandCount;
        }
        ++point;
    }

    free(handPoints);

    printf("index %d active hand count %d\n", frameIndex, activeHandCount);
}

void runHandStream(astra_reader_t reader)
{
    do
    {
        astra_temp_update();

        astra_reader_frame_t frame;
        astra_status_t rc = astra_reader_open_frame(reader, 0, &frame);

        if (rc == ASTRA_STATUS_SUCCESS)
        {
            astra_handframe_t handFrame;
            astra_frame_get_handframe(frame, &handFrame);

            print_hand_frame(handFrame);

            astra_colorframe_t handDebugimageframe;
            astra_frame_get_debug_handframe(frame, &handDebugimageframe);

            astra_image_metadata_t metadata;
            astra_colorframe_get_metadata(handDebugimageframe, &metadata);

            astra_reader_close_frame(&frame);
        }

    } while (shouldContinue);
}

int main(int argc, char* argv[])
{
    set_key_handler();

    astra_initialize();

    astra_streamsetconnection_t sensor;

    astra_streamset_open("device/default", &sensor);

    astra_reader_t reader;
    astra_reader_create(sensor, &reader);

    astra_handstream_t handStream;
    astra_reader_get_handstream(reader, &handStream);
    astra_stream_start(handStream);

    astra_colorstream_t handDebugimagestream;
    astra_reader_get_debug_handstream(reader, &handDebugimagestream);
    astra_stream_start(handDebugimagestream);

    runHandStream(reader);

    astra_reader_destroy(&reader);
    astra_streamset_close(&sensor);

    astra_terminate();
}
