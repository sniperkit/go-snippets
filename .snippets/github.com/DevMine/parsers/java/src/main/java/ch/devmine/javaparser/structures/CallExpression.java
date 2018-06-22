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
public class CallExpression extends Expr {

    private static final String CALL = "CALL";

    @Getter @SerializedName(value = "expression_name")
    private final String exprName = CALL;

    @Getter @Setter @SerializedName(value = "function")
    private FuncRef function;

    @Getter @Setter
    private List<Expr> arguments;

    @Getter @Setter
    private int line;

    public CallExpression() {}

    public CallExpression(FuncRef function, List<Expr> arguments, int line) {
        this.function = function;
        this.arguments = arguments;
        this.line = line;
    }

    @Override
    public String toString() {
        String e = this.exprName;
        String f = this.function != null ? this.function.toString() : "null";
        String a = this.arguments != null ? this.arguments.toString() : "null";
        String l = this.line != 0 ? Integer.toString(this.line) : "null";

        return "expression name : ".concat(e)
                .concat(", function : ").concat(f)
                .concat(", arguments : ").concat(a)
                .concat(", line : ").concat(l);
    }
}
