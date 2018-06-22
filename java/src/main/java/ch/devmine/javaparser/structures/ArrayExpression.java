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
public class ArrayExpression extends Expr {

    private static final String ARRAY = "ARRAY";

    @Getter @SerializedName(value = "expression_name")
    private final String expressionName = ARRAY;

    @Getter @Setter
    private ArrayType type;

    public ArrayExpression() {}

    public ArrayExpression(ArrayType type) {
        this.type = type;
    }

    @Override
    public String toString() {
        String e = this.expressionName;
        String t = this.type!= null ? this.type.toString() : "null";

        return "name : ".concat(e)
                .concat(", type : ").concat(t);
    }
}
