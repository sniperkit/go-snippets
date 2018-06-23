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
public class ProtoDecl {

    @Getter @Setter
    private List<String> doc;

    @Getter @Setter
    private Ident name;

    @Getter @Setter
    private FuncType type;

    @Getter @Setter
    private String visibility;

    public ProtoDecl() {}

    public ProtoDecl(List<String> doc, Ident ident, FuncType type, String visibility) {
        this.doc = doc;
        this.name = ident;
        this.type = type;
        this.visibility = visibility;
    }

    @Override
    public String toString() {
        String d = this.doc != null ? this.doc.toString() : "null";
        String i = this.name != null ? this.name.toString() : "null";
        String t = this.type != null ? this.type.toString() : "null";
        String v = this.visibility != null ? this.visibility : "null";

        return "doc : ".concat(d)
                .concat(", ident : ").concat(i)
                .concat(", type : ").concat(t)
                .concat(", visibility : ").concat(v);

    }
}
