/*
 * Copyright 2014-2015 The DevMine Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ch.devmine.javaparser.structures;

import java.util.List;
import lombok.Getter;
import lombok.Setter;

/**
 *
 * @author lweingart
 */
public class CaseClause {

    @Getter @Setter
    private List<Expr> conditions;

    @Getter @Setter
    private List<Stmt> body;

    public CaseClause() {}

    public CaseClause(List<Expr> conditions, List<Stmt> body) {
        this.conditions = conditions;
        this.body = body;
    }

    @Override
    public String toString() {
        String c = this.conditions != null ? this.conditions.toString() : "null";
        String b = this.body != null ? this.body.toString() : "null";

        return "conditions : ".concat(c)
                .concat(", body : ").concat(b);
    }
}
