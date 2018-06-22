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
public class BinaryExpression extends Expr {

    private static final String BINARY = "BINARY";

    @Getter @SerializedName(value = "expression_name")
    private final String expressionName = BINARY;

    @Getter @Setter @SerializedName(value = "left_expression")
    private Expr leftExpr;

    @Getter @Setter
    private String operator;

    @Getter @Setter @SerializedName(value ="right_expression")
    private Expr rightExpr;

    public BinaryExpression() {}

    public BinaryExpression(Expr leftExpr, String operator, Expr rightExpr) {
        this.leftExpr = leftExpr;
        this.operator = operator;
        this.rightExpr = rightExpr;
    }

    @Override
    public String toString() {
        String e = this.expressionName;
        String l = this.leftExpr != null ? this.leftExpr.toString() : "null";
        String o = this.operator != null ? this.operator : "null";
        String r = this.rightExpr != null ? this.rightExpr.toString() : "null";

        return "expression name : ".concat(e)
                .concat(", left expr : ").concat(l)
                .concat(", operator : ").concat(o)
                .concat(", tight expr : ").concat(r);
    }
}
