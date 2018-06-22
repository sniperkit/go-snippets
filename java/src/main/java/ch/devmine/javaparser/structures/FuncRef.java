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
public class FuncRef {

    @Getter @Setter
    private String namespace;

    @Getter @Setter @SerializedName(value = "function_name")
    private String funcName;

    public FuncRef() {}

    public FuncRef(String namespace, String funcName) {
        this.namespace = namespace;
        this.funcName = funcName;
    }

    @Override
    public String toString() {
        String n = (this.namespace != null) ? this.namespace : "null";
        String f = (this.funcName != null) ? this.funcName : "null";

        return "namespace : ".concat(n)
                .concat(", function_name : ").concat(f);
    }
}
