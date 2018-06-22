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
public class IfStatement extends Stmt {

    private static final String IF = "IF";

    @Getter @SerializedName(value = "statement_name")
    private final String stmtName = IF;

    @Getter @Setter
    private Stmt initialization;

    @Getter @Setter
    private Expr condition;

    @Getter @Setter
    private List<Stmt> body;

    @Getter @Setter @SerializedName(value = "else")
    private List<Stmt> elseStmt;

    @Getter @Setter
    private int line;

    public IfStatement() {}

    public IfStatement(Stmt initialization, Expr condition, List<Stmt> stmtList, List<Stmt> elseStmt, int line) {
        this.initialization = initialization;
        this.condition = condition;
        this.body = stmtList;
        this.elseStmt = elseStmt;
        this.line = line;
    }

    @Override
    public String toString() {
        String t = this.stmtName;
        String i = this.initialization != null ? this.initialization.toString() : "null";
        String c = this.condition != null ? this.condition.toString() : "null";
        String s = this.body != null ? this.body.toString() : "null";
        String e = this.elseStmt != null ? this.elseStmt.toString() : "null";
        String l = this.line != 0 ? Integer.toString(this.line) : "null";

        return "name : ".concat(t)
                .concat(", initialization : ").concat(i)
                .concat(", condition : ").concat(c)
                .concat(", statement_list : ").concat(s)
                .concat(", else : ").concat(e)
                .concat(", line : ").concat(l);
    }
}
