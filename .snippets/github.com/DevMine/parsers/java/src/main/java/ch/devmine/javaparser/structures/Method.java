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
public class Method {

    @Getter @Setter
    private List<String> doc;

    @Getter @Setter
    private String name;

    @Getter @Setter
    private FuncType type;

    @Getter @Setter
    private List<Stmt> body;

    @Getter @Setter
    private String visibility;

    @Getter @Setter
    private int loc;

    @Getter @Setter
    private boolean override;

    public Method() {}

    public Method(List<String> doc, String name, FuncType type, List<Stmt> body,
                  String visibility, int loc, boolean override)
    {
        this.doc = doc;
        this.name = name;
        this.type = type;
        this.body = body;
        this.visibility = visibility;
        this.loc = loc;
        this.override = override;
    }

    @Override
    public String toString() {
        String d = this.doc != null ? this.doc.toString() : "null";
        String n = this.name != null ? this.name : "null";
        String t = this.type != null ? this.type.toString() : "null";
        String b = this.body != null ? this.body.toString() : "null";
        String v = this.visibility != null ? this.visibility : "null";
        String o = this.override ? "true" : "false";
        String l = this.loc != 0 ? Integer.toString(this.loc) : "null";

        return "doc : ".concat(d)
                .concat(", name : ").concat(n)
                .concat(", type : ").concat(t)
                .concat(", body : ").concat(b)
                .concat(", visibility : ").concat(v)
                .concat(", override : ").concat(o)
                .concat(", line : ").concat(l);
    }
}
