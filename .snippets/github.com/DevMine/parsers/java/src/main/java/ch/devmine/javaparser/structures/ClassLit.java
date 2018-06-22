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
public class ClassLit extends Expr {

    private static final String CLASS_LIT = "CLASS_LIT";

    @Getter @SerializedName(value = "expression_name")
    private final String exprName = CLASS_LIT;

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

    public ClassLit() {}

    public ClassLit(List<ClassRef> extendedClasses, List<InterfaceRef> implementedInterfaces, List<Attribute> attributes,
                    List<ConstructorDecl> constructors, List<Method> methods)
    {
        this.extendedClasses = extendedClasses;
        this.implementedInterfaces = implementedInterfaces;
        this.attributes = attributes;
        this.constructors = constructors;
        this.methods = methods;
    }

    @Override
    public String toString() {
        String e = this.exprName != null ? this.exprName : "null";
        String ex = this.extendedClasses != null ? this.extendedClasses.toString() : "null";
        String i = this.implementedInterfaces != null ? this.implementedInterfaces.toString() : "null";
        String a = this.attributes != null ? this.attributes.toString() : "null";
        String c = this.constructors != null ? this.constructors.toString() : "null";
        String m = this.methods != null ? this.methods.toString() : "null";

        return "expression name : ".concat(e)
                .concat(", extended classes : ").concat(ex)
                .concat(", implemented interfaces : ").concat(i)
                .concat(", attributes : ").concat(a)
                .concat(", constructors : ").concat(c)
                .concat(", methods : ").concat(m);
    }
}
