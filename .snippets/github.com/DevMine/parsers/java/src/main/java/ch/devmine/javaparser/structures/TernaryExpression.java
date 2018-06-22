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
public class TernaryExpression extends Expr {

    private static final String TERNARY = "TERNARY";

    @Getter @SerializedName(value = "expression_name")
    private final String expressionName = TERNARY;

    @Getter @Setter
    private Expr condition;

    @Getter @Setter
    private Expr then;

    @Getter @Setter @SerializedName(value = "else")
    private Expr els;

    public TernaryExpression() {}

    public TernaryExpression(Expr condition, Expr then, Expr els) {
        this.condition = condition;
        this.then = then;
        this.els = els;
    }

    @Override
    public String toString() {
        String e = this.expressionName;
        String c = this.condition != null ? this.condition.toString() : "null";
        String t = this.then != null ? this.then.toString() : "null";
        String el = this.els != null ? this.els.toString() : "null";

        return "expression name : ".concat(e)
                .concat(", condition : ").concat(c)
                .concat(", then : ").concat(t)
                .concat(", else : ").concat(el);
    }
}
