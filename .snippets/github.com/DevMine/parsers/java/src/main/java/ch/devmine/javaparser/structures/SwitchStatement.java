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
public class SwitchStatement extends Stmt {

    private static final String SWITCH = "SWITCH";

    @Getter @SerializedName(value = "statement_name")
    private final String statementName = SWITCH;

    @Getter @Setter
    private Expr initialization;

    @Getter @Setter
    private Expr condition;

    @Getter @Setter @SerializedName(value = "case_clauses")
    private List<CaseClause> caseClauses;

    @Getter @Setter @SerializedName(value = "default")
    private List<Stmt> defaultStmt;

    public SwitchStatement() {}

    public SwitchStatement(Expr initialization, Expr condition, List<CaseClause> caseClauses, List<Stmt> defaultStmt) {
        this.initialization = initialization;
        this.condition = condition;
        this.caseClauses = caseClauses;
        this.defaultStmt = defaultStmt;
    }

    @Override
    public String toString() {
        String s = this.statementName;
        String c = this.condition != null ? this.condition.toString() : "null";
        String cs = this.caseClauses != null ? this.caseClauses.toString() : "null";
        String d = this.defaultStmt != null ? this.defaultStmt.toString() : "null";

        return "Statement name : ".concat(s)
                .concat(", condition : ").concat(c)
                .concat(", caseClauses : ").concat(cs)
                .concat(", defaultStmt : ").concat(d);
    }

}
