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
public class BasicLit extends Expr {

    private static final String BASIC_LIT = "BASIC_LIT";

    @Getter @SerializedName(value = "expression_name")
    private final String exprName = BASIC_LIT;

    @Getter @Setter
    private String kind;

    @Getter @Setter
    private String value;

    public BasicLit() {}

    public BasicLit(String kind, String value) {
        this.kind = kind;
        this.value = value;
    }

    @Override
    public String toString() {
        String e = this.exprName;
        String k = this.kind != null ? this.kind : "null";
        String v = this.value != null ? this.value : "null";

        return "expression name : ".concat(e)
                .concat(", kind : ").concat(k)
                .concat(". value : ").concat(v);
    }
}
