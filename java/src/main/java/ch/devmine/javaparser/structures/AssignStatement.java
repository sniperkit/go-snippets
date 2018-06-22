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
public class AssignStatement extends Stmt {

    private static final String ASSIGN = "ASSIGN";

    @Getter @SerializedName(value = "statement_name")
    private final String stmtName = ASSIGN;

    @Getter @Setter @SerializedName(value = "left_hand_side")
    private List<Expr> leftHandExpr;

    @Getter @Setter @SerializedName(value = "right_hand_side")
    private List<Expr> rightHandExpr;

    @Getter @Setter
    private int line;

    public AssignStatement() {}

    public AssignStatement(List<Expr> leftHandStmt, List<Expr> rightHandStmt,
                           int line)
    {
        this.leftHandExpr = leftHandStmt;
        this.rightHandExpr = rightHandStmt;
        this.line = line;
    }

    @Override
    public String toString() {
        String t = this.stmtName != null ? this.stmtName : "null";
        String n = this.leftHandExpr != null ? this.leftHandExpr.toString() : "null";
        String v = this.rightHandExpr != null ? this.rightHandExpr.toString() : "null";
        String l = this.line != 0 ? Integer.toString(this.line) : "null";

        return "statement name : ".concat(t)
                .concat(", name : ").concat(n)
                .concat(", value : ").concat(v)
                .concat(", line num : ").concat(l);
    }
}
