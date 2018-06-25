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
#include "hnd_scaling_coordinate_mapper.hpp"
#include <Shiny.h>

namespace astra { namespace hand {

    void convert_depth_to_world_f(const conversion_cache_t& depthToWorldData,
                                  float depthX, float depthY, float depthZ,
                                  float& worldX, float& worldY, float& worldZ)
    {
        float normalizedX = depthX / depthToWorldData.resolutionX - .5f;
        float normalizedY = .5f - depthY / depthToWorldData.resolutionY;

        worldX = normalizedX * depthZ * depthToWorldData.xzFactor;
        worldY = normalizedY * depthZ * depthToWorldData.yzFactor;
        worldZ = depthZ;
    }

    Vector3f cv_convert_depth_to_world(const conversion_cache_t& depthToWorldData,
                                       float depthX, float depthY, float depthZ)
    {
        float worldX, worldY, worldZ;

        convert_depth_to_world_f(depthToWorldData,
                                 depthX, depthY, depthZ,
                                 worldX, worldY, worldZ);

        return Vector3f(worldX, worldY, worldZ);
    }

    Vector3f cv_convert_depth_to_world(const conversion_cache_t& depthToWorldData,
                                       int depthX, int depthY, float depthZ)
    {
        return cv_convert_depth_to_world(depthToWorldData,
                                         static_cast<float>(depthX),
                                         static_cast<float>(depthY),
                                         depthZ);
    }

    Vector3f cv_convert_depth_to_world(const conversion_cache_t& depthToWorldData,
                                       const Vector3f& depth)
    {
        return cv_convert_depth_to_world(depthToWorldData, depth.x, depth.y, depth.z);
    }

    Vector3f cv_convert_world_to_depth(const conversion_cache_t& depthToWorldData,
                                       float worldX, float worldY, float worldZ)
    {
        float depthX = depthToWorldData.coeffX * worldX / worldZ + depthToWorldData.halfResX;
        float depthY = depthToWorldData.halfResY - depthToWorldData.coeffY * worldY / worldZ;
        float depthZ = worldZ;

        return Vector3f(depthX, depthY, depthZ);
    }

    Vector3f cv_convert_world_to_depth(const conversion_cache_t& depthToWorldData,
                                       const Vector3f& world)
    {
        return cv_convert_world_to_depth(depthToWorldData, world.x, world.y, world.z);
    }


    Point2i scaling_coordinate_mapper::offset_pixel_location_by_mm(const Point2i& position,
                                                                   float offsetX,
                                                                   float offsetY,
                                                                   float depthZ) const
    {
        if (depthZ == 0)
        {
            return position;
        }

        const conversion_cache_t& depthToWorldData = depthToWorldData_;

        const float scaledDepth = depthZ * scale_;

        //This bakes the world to depth conversion and full size to scaled conversion into
        //a single factor that we apply to the requested offset.
        //This process is equivalent to depth_to_world, adding offsetX/Y to world position,
        //then doing world_to_depth.

        const float xFactor = depthToWorldData.resolutionX / (scaledDepth * depthToWorldData.xzFactor);
        const float finalDepthX = position.x + offsetX * xFactor;

        const float yFactor = depthToWorldData.resolutionY / (scaledDepth * depthToWorldData.yzFactor);
        const float finalDepthY = position.y - offsetY * yFactor;

        return Point2i(static_cast<int>(finalDepthX), static_cast<int>(finalDepthY));
    }
}}
