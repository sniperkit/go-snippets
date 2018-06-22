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
public class ConstructorDecl {

    @Getter @Setter
    private List<String> doc;

    @Getter @Setter
    private String name;

    @Getter @Setter
    private List<Field> parameters;

    @Getter @Setter
    private List<Stmt> body;

    @Getter @Setter
    private String visibility;

    @Getter @Setter
    private int loc;

    public ConstructorDecl() {}

    public ConstructorDecl(List<String> doc, String name, List<Field> parameters,
                           List<Stmt> body, String visibility, int loc)
    {
        this.doc = doc;
        this.name = name;
        this.parameters = parameters;
        this.body = body;
        this.visibility = visibility;
        this.loc = loc;
    }

    @Override
    public String toString() {
       String d = (this.doc != null) ? this.doc.toString() : "null";
       String n = (this.name != null) ? this.name : "null";
       String p = (this.parameters != null) ? this.parameters.toString() : "null";
       String b = (this.body != null) ? this.body.toString() : "null";
       String v = (this.visibility != null) ? this.visibility : "null";
       String l = (this.loc != 0) ? Integer.toString(this.loc) : "null";

               return "doc : ".concat(d)
                       .concat(", name : ").concat(n)
                       .concat(", parameters : ").concat(p)
                       .concat(", body : ").concat(b)
                       .concat(", visibility: ").concat(v)
                       .concat(", loc : ").concat(l);
    }
}
