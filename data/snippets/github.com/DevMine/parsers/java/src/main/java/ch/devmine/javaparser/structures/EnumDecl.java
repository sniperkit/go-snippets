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
public class EnumDecl {

    @Getter @Setter
    private List<String> doc;

    @Getter @Setter
    private String name;

    @Getter @Setter
    private String visibility;

    @Getter @Setter @SerializedName(value = "implemented_interfaces")
    private List<InterfaceRef> implementedInterfaces;

    @Getter @Setter @SerializedName(value = "enum_constants")
    private List<Ident> enumConstants;

    @Getter @Setter
    private List<Attribute> attributes;

    @Getter @Setter
    private List<ConstructorDecl> constructors;

    @Getter @Setter
    private List<Method> methods;

    public EnumDecl() {}

    public EnumDecl(List<String> doc, String name, String visibility, List<InterfaceRef> implementedInterfaces,
                    List<Ident> enumConstants, List<Attribute> attributes, List<ConstructorDecl> constructors,
                    List<Method> methods)
    {
        this.doc = doc;
        this.name = name;
        this.visibility = visibility;
        this.implementedInterfaces = implementedInterfaces;
        this.enumConstants = enumConstants;
        this.attributes = attributes;
        this.constructors = constructors;
        this.methods = methods;
    }

    @Override
    public String toString() {
        String d = this.doc != null ? this.doc.toString() : "null";
        String n = this.name != null ? this.name : "null";
        String v = this.visibility != null ? this.visibility : "null";
        String i = this.implementedInterfaces != null ? this.implementedInterfaces.toString() : "nul";
        String e = this.enumConstants != null ? this.enumConstants.toString() : "null";
        String a = this.attributes != null ? this.attributes.toString() : "null";
        String c = this.constructors != null ? this.constructors.toString() : "null";
        String m = this.methods != null ? this.methods.toString() : "null";

        return "doc : ".concat(d)
                .concat("name : ").concat(n)
                .concat(", visibility : ").concat(v)
                .concat(", implemented interfaces : ").concat(i)
                .concat(", enum const : ").concat(e)
                .concat(", attributes : ").concat(a)
                .concat(", constructors : ").concat(c)
                .concat(", methods : ").concat(m);
    }
}
