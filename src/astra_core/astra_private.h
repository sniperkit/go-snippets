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
#ifndef ASTRA_PRIVATE_H
#define ASTRA_PRIVATE_H

using astra_frame_id_t = int32_t;

using astra_frame_status_t = int8_t;
const astra_frame_status_t ASTRA_FRAME_STATUS_AVAILABLE = 0;
const astra_frame_status_t ASTRA_FRAME_STATUS_LOCKED_POLL = 1;
const astra_frame_status_t ASTRA_FRAME_STATUS_LOCKED_EVENT = 2;

struct _astra_reader_frame
{
    astra_frame_id_t id;
    astra_frame_status_t status;
    astra_reader_t reader;
};

#endif // ASTRA_PRIVATE_H
