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
public class ClassDecl {

    @Getter @Setter
    private List<String> doc;

    @Getter @Setter
    private String name;

    @Getter @Setter
    private String visibility;

    @Getter @Setter @SerializedName(value = "extended_classes")
    private List<ClassRef> extendedClasses;

    @Getter @Setter @SerializedName(value = "implemented_interfaces")
    private List<InterfaceRef> implementedInterfaces;

    @Getter @Setter
    private List<Attribute> attributes;

    @Getter @Setter
    private List<ConstructorDecl> constructors;

    @Getter @Setter
    private List<Method> methods;

    @Getter @Setter @SerializedName(value = "nested_classes")
    private List<ClassDecl> nestedClasses;

    public ClassDecl() {}

    public ClassDecl(List<String> doc, String name, String visibility, List<ClassRef> extendedClasses,
                     List<InterfaceRef> implementedInterfaces, List<Attribute> attributes,
                     List<ConstructorDecl> constructors, List<Method> methods, List<ClassDecl> nestedClasses)
    {
        this.doc = doc;
        this.name = name;
        this.visibility = visibility;
        this.extendedClasses = extendedClasses;
        this.implementedInterfaces = implementedInterfaces;
        this.attributes = attributes;
        this.constructors = constructors;
        this.methods = methods;
        this.nestedClasses = nestedClasses;
    }

    @Override
    public String toString() {
        String d = this.doc != null ? this.doc.toString() : "null";
        String n = this.name != null ? this.name : "null";
        String v = this.visibility != null ? this.visibility : "null";
        String e = this.extendedClasses != null ? this.extendedClasses.toString() : "null";
        String i = this.implementedInterfaces != null ? this.implementedInterfaces.toString() : "null";
        String a = this.attributes != null ? this.attributes.toString() : "null";
        String c = this.constructors != null ? this.constructors.toString() : "null";
        String m = this.methods != null ? this.methods.toString() : "null";
        String ne = this.nestedClasses != null ? this.nestedClasses.toString() : "null";

        return "doc : ".concat(d)
                .concat("name : ").concat(n)
                .concat(", visibility : ").concat(v)
                .concat(", extended classe : ").concat(e)
                .concat(", implemented interfaces : ").concat(i)
                .concat(", attributes : ").concat(a)
                .concat(", constructors : ").concat(c)
                .concat(", methods : ").concat(m)
                .concat(", nested classes : ").concat(ne);
    }
}
