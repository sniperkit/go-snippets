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
public class IndexExpression extends Expr {

    private static final String INDEX = "INDEX";

    @Getter @SerializedName(value = "expression_name")
    private final String exprName = INDEX;

    @Getter @Setter
    private Expr expression;

    @Getter @Setter
    private Expr index;

    public IndexExpression() {}

    public IndexExpression(Expr expression, Expr index) {
        this.expression = expression;
        this.index = index;
    }

    @Override
    public String toString() {
        String e = this.exprName;
        String ex = this.expression != null ? this.expression.toString() : "null";
        String i = this.index != null ? this.index.toString() : "null";

        return "expression name : ".concat(e)
                .concat(", expression : ").concat(ex)
                .concat(", index : ").concat(i);
    }
}
