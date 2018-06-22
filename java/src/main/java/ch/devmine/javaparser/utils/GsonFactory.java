/*
 * Copyright 2014-2015 The DevMine Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ch.devmine.javaparser.utils;

import ch.devmine.javaparser.annotations.GsonTransient;
import com.google.gson.ExclusionStrategy;
import com.google.gson.FieldAttributes;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;


/**
 *
 * @author lweingart
 */
public class GsonFactory {

    public static Gson build() {
        GsonBuilder builder = new GsonBuilder();
        builder.addSerializationExclusionStrategy(new ExclusionStrategy() {

            @Override
            public boolean shouldSkipField(FieldAttributes f) {
                return f.getAnnotation(GsonTransient.class) != null;
            }

            @Override
            public boolean shouldSkipClass(Class<?> c) {
                return c.getAnnotation(GsonTransient.class) != null;
            }
        });
        return builder.create();
    }
}

