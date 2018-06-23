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
public class Field {

    @Getter @Setter
    private List<String> doc;

    @Getter @Setter
    private String name;

    @Getter @Setter
    private String type;

    public Field() {}

    public Field(List<String> doc, String name, String type) {
        this.doc = doc;
        this.name = name;
        this.type = type;
    }

    @Override
    public String toString() {
        String d = this.doc != null ? this.doc.toString() : "null";
        String n = this.name != null ? this.name : "null";
        String t = this.type != null ? this.type : "null";

        return "doc : ".concat(d)
                .concat(", type : ").concat(t)
                .concat(", name : ").concat(n);
    }
}
