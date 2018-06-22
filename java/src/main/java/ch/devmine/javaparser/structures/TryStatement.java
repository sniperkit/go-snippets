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
public class TryStatement extends Stmt {

    private static final String TRY = "TRY";

    @Getter @SerializedName(value = "statement_name")
    private final String stmtName = TRY;

    @Getter @Setter
    private List<Stmt> body;

    @Getter @Setter @SerializedName(value = "catch_clauses")
    private List<CatchClause> catchClauses;

    @Getter @Setter @SerializedName(value = "finally")
    private List<Stmt> finallyStmts;

    public TryStatement() {}

    public TryStatement(List<Stmt> body, List<CatchClause> catchClauses, List<Stmt> finallyStmts) {
        this.body = body;
        this.catchClauses = catchClauses;
        this.finallyStmts = finallyStmts;
    }

    @Override
    public String toString() {
        String s = this.stmtName;
        String b = this.body != null ? this.body.toString() : "null";
        String c = this.catchClauses != null ? this.catchClauses.toString() : "null";
        String f = this.finallyStmts != null ? this.finallyStmts.toString() : "null";

        return "statement name : ".concat(s)
                .concat(", body : ").concat(b)
                .concat(", catch clause : ").concat(c)
                .concat(", finally : ").concat(f);
    }
}
