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
public class FuncType {

    @Getter @Setter
    private List<Field> parameters;

    @Getter @Setter
    private List<Field> results;

    public FuncType() {}

    public FuncType(List<Field> parameters, List<Field> results) {
        this.parameters = parameters;
        this.results = results;
    }

    @Override
    public String toString() {
        String p = this.parameters != null ? this.parameters.toString() : "null";
        String r = this.results != null ? this.results.toString() : "null";

        return "params : ".concat(p)
                .concat(r).concat(r);
    }

}
