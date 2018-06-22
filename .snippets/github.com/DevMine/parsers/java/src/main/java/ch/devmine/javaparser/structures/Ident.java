/*
 * Copyright 2014-2015 The DevMine Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ch.devmine.javaparser.structures;

import com.google.gson.annotations.SerializedName;
import lombok.Getter;
import lombok.Setter;

/**
 *
 * @author lweingart
 */
public class Ident extends Expr {

    private static final String IDENT = "IDENT";

    @Getter @SerializedName(value = "expression_name")
    private final String exprName = IDENT;

    @Getter @Setter
    private String name;

    public Ident() {}

    public Ident(String name) {
        this.name = name;
    }

    @Override
    public String toString() {
        String e = (this.exprName != null) ? this.exprName : "null";
        String n = (this.name != null) ? this.name : "null";

        return "expression name : ".concat(e)
                .concat(", name : ").concat(n);
    }
}
