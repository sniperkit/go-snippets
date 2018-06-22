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
public class Variable {

    @Getter @Setter
    private List<String> doc;

    @Getter @Setter
    private String name;

    @Getter @Setter
    private String type;

    @Getter @Setter
    private String value;

    public Variable() {}

    public Variable(List<String> doc, String name, String type, String value) {
        this.doc = doc;
        this.name = name;
        this.type = type;
        this.value = value;
    }



    @Override
    public String toString() {
        String d = (this.doc != null) ? this.doc.toString() : "null";
        String n = (this.name != null) ? this.name : "null";
        String t = (this.type != null) ? this.type : "null";
        String v = (this.value != null) ? this.value : "null";

        return "doc : ".concat(d)
                .concat(", name : ").concat(n)
                .concat(", type : ").concat(t)
                .concat(", value : ").concat(v);
    }
}
