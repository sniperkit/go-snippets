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
public class ValueSpec extends Expr {

    private static final String VALUE_SPEC = "VALUE_SPEC";

    @Getter @SerializedName(value = "expression_name")
    private final String exprName = VALUE_SPEC;

    @Getter @Setter
    private Ident name;

    @Getter @Setter
    private Ident type;

    public ValueSpec() {}

    public ValueSpec(Ident name, Ident type) {
        this.name = name;
        this.type = type;
    }

    @Override
    public String toString() {
        String e = this.exprName;
        String n = this.name != null ? this.name.toString() : "null";
        String t = this.type != null ? this.type.toString() : "null";

        return "expression name : ".concat(e)
                .concat(", name : ").concat(n)
                .concat(", type : ").concat(t);
    }
}
