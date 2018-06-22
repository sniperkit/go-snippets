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
public class ReturnStatement extends Stmt {

    private static final String RETURN = "RETURN";

    @Getter @SerializedName(value = "statement_name")
    private final String statementName = RETURN;

    @Getter @Setter
    private List<Expr> results;

    @Getter @Setter
    private int line;

    public ReturnStatement() {}

    public ReturnStatement(List<Expr> results, int line) {
        this.results = results;
        this.line = line;
    }

    @Override
    public String toString() {
        String t = this.statementName;
        String v = this.results != null ? this.results.toString() : "null";
        String l = this.line != 0 ? Integer.toString(this.line) : "null";

        return "type : " + t
                .concat(", value : " + v)
                .concat(", line : " + l);
    }
}
