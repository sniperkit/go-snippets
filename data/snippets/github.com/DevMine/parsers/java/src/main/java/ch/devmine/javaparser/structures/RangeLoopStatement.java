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
public class RangeLoopStatement extends Stmt {

    private static final String RANGE_LOOP = "RANGE_LOOP";

    @Getter @SerializedName(value = "statement_name")
    private final String StatementName = RANGE_LOOP;

    @Getter @Setter
    private List<Expr> variables;

    @Getter @Setter
    private Expr iterable;

    @Getter @Setter
    private List<Stmt> body;

    @Getter @Setter
    private int line;

    public RangeLoopStatement() {}

    public RangeLoopStatement(String StatementName, List<Expr> variables, Expr iterable, List<Stmt> body, int line) {
        this.variables = variables;
        this.iterable = iterable;
        this.body = body;
        this.line = line;
    }

    @Override
    public String toString() {
        String s = this.StatementName;
        String v = this.variables != null ? this.variables.toString() : "null";
        String i = this.iterable != null ? this.iterable.toString() : "null";
        String b = this.body != null ? this.body.toString() : "null";
        String l = this.line != 0 ? Integer.toString(this.line) : "null";

        return "name : ".concat(s)
                .concat(", variables : ").concat(v)
                .concat(", iterable : ").concat(i)
                .concat(", body : ").concat(b)
                .concat(", line : ").concat(l);
    }

}
