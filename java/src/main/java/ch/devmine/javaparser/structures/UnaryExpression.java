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
public class UnaryExpression extends Expr {

    private static final String UNARY = "UNARY";

    @Getter @SerializedName(value = "expression_name")
    private final String expressionName = UNARY;

    @Getter @Setter
    private String operator;

    @Getter @Setter
    private Expr operand;

    public UnaryExpression() {}

    public UnaryExpression(String operator, Expr operand) {
        this.operator = operator;
        this.operand = operand;
    }

    @Override
    public String toString() {
        String e = this.expressionName;
        String op1 = this.operator != null ? this.operator : "null";
        String op2 = this.operand != null ? this.operand.toString() : "null";

        return "Expression name : ".concat(e)
                .concat(", operator : ").concat(op1)
                .concat(", operand : " ).concat(op2);
    }
}
