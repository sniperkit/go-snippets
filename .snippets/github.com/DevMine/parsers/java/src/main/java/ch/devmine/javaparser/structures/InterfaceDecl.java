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
public class InterfaceDecl {

    @Getter @Setter
    private List<String> doc;

    @Getter @Setter
    private String name;

    @Getter @Setter @SerializedName(value = "implemented_interfaces")
    private List<InterfaceRef> implementedInterfaces;

    @Getter @Setter
    private List<ProtoDecl> prototypes;

    @Getter @Setter
    private String visibility;

    public InterfaceDecl() {}

    public InterfaceDecl(List<String> doc, String name, List<ProtoDecl> prototypes, String visibility) {
        this.doc = doc;
        this.name = name;
        this.prototypes = prototypes;
        this.visibility = visibility;
    }

    @Override
    public String toString() {
        String d = this.doc != null ? this.doc.toString() : "null";
        String n = this.name != null ? this.name : "null";
        String p = this.prototypes != null ? this.prototypes.toString() : "null";
        String v = this.visibility != null ? this.visibility : "null";

        return "doc : ".concat(d)
                .concat("name : ").concat(n)
                .concat(", prototypes : ").concat(p)
                .concat(", visibility : ").concat(v);
    }
}
