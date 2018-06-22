/*
 * Copyright 2014-2015 The DevMine Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ch.devmine.javaparser.structures;

import com.google.gson.annotations.SerializedName;
import java.util.List;
import lombok.Getter;
import lombok.Setter;

/**
 *
 * @author lweingart
 */
public class ArrayLit extends Expr {

    private static final String ARRAY_LIT = "ARRAY_LIT";

    @Getter @SerializedName(value = "expression_name")
    private final String expressionName = ARRAY_LIT;

    @Getter @Setter
    private ArrayType type;

    @Getter @Setter
    private List<Expr> elements;

    public ArrayLit() {}

    public ArrayLit(ArrayType type, List<Expr> elements) {
        this.type = type;
        this.elements = elements;
    }

    @Override
    public String toString() {
        String t = this.type != null ? this.type.toString() : "null";
        String e = this.elements != null ? this.elements.toString() : "null";

        return "type : ".concat(t)
                .concat(", elements : ").concat(e);
    }
}
