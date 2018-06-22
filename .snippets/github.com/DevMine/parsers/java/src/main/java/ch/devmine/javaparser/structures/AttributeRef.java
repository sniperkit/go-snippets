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
public class AttributeRef extends Expr {

    private static final String ATTR_REF = "ATTR_REF";

    @Getter @SerializedName(value = "expression_name")
    private final String exprName = ATTR_REF;

    @Getter @Setter
    private Ident name;

    public AttributeRef() {}

    public AttributeRef(Ident name) {
        this.name = name;
    }

    @Override
    public String toString() {
        String e = this.exprName != null ? this.exprName : "null";
        String i = this.name != null ? this.name.toString() : "null";

        return "expression name : ".concat(e)
                .concat(", name : ").concat(i);
    }
}
