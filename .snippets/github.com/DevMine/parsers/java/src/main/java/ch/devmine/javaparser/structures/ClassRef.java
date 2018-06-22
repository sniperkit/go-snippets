/*
 * Copyright 2014-2015 The DevMine Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ch.devmine.javaparser.structures;

import com.google.gson.annotations.SerializedName;
import lombok.Getter;
import lombok.Setter;

/**
 *
 * @author lweingart
 */
public class ClassRef {

    @Getter @Setter
    private String namespace;

    @Getter @Setter @SerializedName(value = "class_name")
    private String className;

    public ClassRef() {}

    public ClassRef(String namespace, String className) {
        this.namespace = namespace;
        this.className = className;
    }

    @Override
    public String toString() {
        String n = (this.namespace != null) ? this.namespace : "null";
        String c = (this.className != null) ? this.className : "null";

        return "namespace : ".concat(n)
                .concat(", class name : ").concat(c);
    }
}
