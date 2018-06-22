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
public class CatchClause {

    @Getter @Setter
    private List<Field> parameters;

    @Getter @Setter
    private List<Stmt> body;

    public CatchClause() {}

    public CatchClause(List<Field> parameters, List<Stmt> body) {
        this.parameters = parameters;
        this.body = body;
    }

    @Override
    public String toString() {
        String p = (this.parameters != null) ? this.parameters.toString() : "null";
        String b = (this.body != null) ? this.body.toString() : "null";

        return "parameters : ".concat(p)
                .concat(", body : ").concat(b);
    }
}
