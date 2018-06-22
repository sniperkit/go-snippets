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
public class ThrowStatement extends Stmt {

    private static final String THROW = "THROW";

    @Getter @SerializedName(value = "statement_name")
    private final String statementName = THROW;

    @Getter @Setter
    private Expr expression;

    public ThrowStatement() {}

    public ThrowStatement(Expr expression) {
        this.expression = expression;
    }

    @Override
    public String toString() {
        String s = this.statementName;
        String e = (this.expression != null) ? this.expression.toString() : "null";

        return "statement name : " .concat(s)
                .concat(", expression : ").concat(e);
    }
}
