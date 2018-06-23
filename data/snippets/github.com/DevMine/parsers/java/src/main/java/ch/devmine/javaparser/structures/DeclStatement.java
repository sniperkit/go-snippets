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
public class DeclStatement extends Stmt {

    private static final String DECL = "DECL";

    @Getter @SerializedName(value = "statement_name")
    private final String statementName = DECL;

    @Getter @Setter @SerializedName(value = "left_hand_side")
    private List<Expr> leftHandExpr;

    @Getter @Setter @SerializedName(value = "right_hand_side")
    private List<Expr> rightHandExpr;

    @Getter @Setter
    private int line;

    @Getter @Setter
    private String kind;

    public DeclStatement(){}

    public DeclStatement(List<Expr> leftHandExpr, List<Expr> rightHandExpr, int line, String kind)
    {
        this.leftHandExpr = leftHandExpr;
        this.rightHandExpr = rightHandExpr;
        this.line = line;
        this.kind = kind;
    }

    @Override
    public String toString() {
        String t = this.statementName;
        String n = this.leftHandExpr != null ? this.leftHandExpr.toString() : "null";
        String v = this.rightHandExpr != null ? this.rightHandExpr.toString() : "null";
        String l = this.line != 0 ? Integer.toString(this.line) : "null";
        String k = this.kind != null ? this.kind : "null";

        return "type : " + t
                .concat(", statement name : " + n)
                .concat(", value : " + v)
                .concat(", line num : " + l)
                .concat(", kind : " + k);
    }
}
