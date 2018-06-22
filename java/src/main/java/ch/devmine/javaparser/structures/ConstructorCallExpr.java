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
public class ConstructorCallExpr extends Expr {

    private static final String CONSTRUCTOR_CALL = "CONSTRUCTOR_CALL";

    @Getter @SerializedName(value = "expression_name")
    private final String expressionName = CONSTRUCTOR_CALL;

    @Getter @Setter
    private FuncRef function;

    @Getter @Setter
    private List<Expr> arguments;

    @Getter @Setter
    private int line;

    public ConstructorCallExpr() {}

    public ConstructorCallExpr(FuncRef function, List<Expr> arguments, int line) {
        this.function = function;
        this.arguments = arguments;
        this.line = line;
    }

    @Override
    public String toString() {
        String e = this.expressionName;
        String f = this.function != null ? this.function.toString() : "null";
        String a = null;
        if (this.arguments != null) {
            a = "[";
            for (Expr argument : arguments) {
                String s = argument != null ? argument.toString() : "null";
                a.concat(s);
            }
            a.concat("]");
        }
        a = a != null ? a : "null";
        String l = this.line != 0 ? Integer.toString(this.line) : "null";

        return "expression name : ".concat(e)
                .concat(", function : ").concat(f)
                .concat(", arguments : ").concat(a)
                .concat(", line : ").concat(l);
    }

}
