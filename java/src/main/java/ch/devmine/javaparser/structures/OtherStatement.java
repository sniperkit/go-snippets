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
public class OtherStatement extends Stmt {

    private static final String OTHER = "OTHER";

    @Getter @SerializedName(value = "statement_name")
    private final String stmtName = OTHER;

    @Getter @Setter
    private List<Stmt> body;

    @Getter @Setter
    private int line;

    public OtherStatement() {}

    public OtherStatement(List<Stmt> stmtList, int line) {
        this.body = stmtList;
        this.line = line;
    }

    @Override
    public String toString() {
        String t = this.stmtName;
        String s = this.body != null ? this.body.toString() : "null";
        String l = this.line != 0 ? Integer.toString(this.line) : "null";

        return "name : ".concat(t)
                .concat(", body : ").concat(s)
                .concat(", line : ").concat(l);
    }
}
