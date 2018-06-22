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
public class ExprStatement extends Stmt {

    private static final String EXPR = "EXPR";

    @Getter @SerializedName(value = "statement_name")
    private final String exprStmtName = EXPR;

    @Getter @Setter
    private Expr expression;

    public ExprStatement() {}

    public ExprStatement(Expr expression) {
        this.expression = expression;
    }

    @Override
    public String toString() {
        String e = this.exprStmtName;
        String ex = this.expression != null ? this.expression.toString() : "null";

        return "expression name : ".concat(e)
                .concat(", expression : ").concat(ex);
    }
}
