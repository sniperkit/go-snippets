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
public class LoopStatement extends Stmt {

    private static final String LOOP = "LOOP";

    @Getter @SerializedName(value = "statement_name")
    private final String statementName = LOOP;

    @Getter @Setter
    private List<Stmt> initialization;

    @Getter @Setter
    private Expr condition;

    @Getter @Setter @SerializedName(value = "post_iteration_statement")
    private List<Stmt> postItStmt;

    @Getter @Setter
    private List<Stmt> body;

    @Getter @Setter @SerializedName(value = "is_post_evaluated")
    private boolean postEvaluated;

    @Getter @Setter
    private int line;

    public LoopStatement() {}

    public LoopStatement(List<Stmt> initialization, Expr condition, List<Stmt> postItStmt,
                         List<Stmt> body, boolean postEvaluated, int line)
    {
        this.initialization = initialization;
        this.condition = condition;
        this.postItStmt = postItStmt;
        this.body = body;
        this.postEvaluated = postEvaluated;
        this.line = line;
    }

    @Override
    public String toString() {
        String t = (this.statementName != null) ? this.statementName : "null";
        String s = (this.body != null) ? this.body.toString() : "null";
        String l = (this.line != 0) ? Integer.toString(this.line) : "null";

        return "type : ".concat(t)
                .concat(", statement_list : ").concat(s)
                .concat(", line : ").concat(l);
    }
}
