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
public class IncDecExpr extends Expr {

    private static final String INC_DEC = "INC_DEC";

    @Getter @SerializedName(value = "expression_name")
    private final String exprName = INC_DEC;

    @Getter @Setter
    private Expr operand;

    @Getter @Setter
    private String operator;

    @Getter @Setter @SerializedName(value = "is_pre")
    private boolean prefix;

    public IncDecExpr() {}

    public IncDecExpr(Expr operand, String operator, boolean prefix) {
        this.operand = operand;
        this.operator = operator;
        this.prefix = prefix;
    }

    @Override
    public String toString() {
        String e = this.exprName;
        String o = this.operand != null ? this.operand.toString() : "null";
        String op = this.operator != null ? this.operator : "null";
        String p = this.prefix ? "true" : "false";

        return "expression name : ".concat(e)
                .concat(", operand : ").concat(o)
                .concat(", operator : ").concat(op);
    }
}
