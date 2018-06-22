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
public class SourceFile {

    @Getter @Setter
    private String path;

    @Getter @Setter
    private Language language;

    @Getter @Setter
    private List<String> imports;

    @Getter @Setter
    private List<InterfaceDecl> interfaces;

    @Getter @Setter
    private List<ClassDecl> classes;

    @Getter @Setter
    private List<EnumDecl> enums;

    @Getter @Setter
    private int loc;

    public SourceFile() {}

    public SourceFile(String path, Language language, List<String> imports, List<InterfaceDecl> interfaces,
                      List<ClassDecl> classes, List<EnumDecl> enums, int loc)
    {
        this.path = path;
        this.language = language;
        this.imports = imports;
        this.interfaces = interfaces;
        this.classes = classes;
        this.enums = enums;
        this.loc = loc;
    }

    @Override
    public String toString() {
        String p = this.path != null ? this.path : "null";
        String l = this.language != null ? this.language.toString() : "null";
        String i = this.imports != null ? this.imports.toString() : "null";
        String in = this.interfaces != null ? this.interfaces.toString() : "null";
        String c = this.classes != null ? this.classes.toString() : "null";
        String e = this.enums != null ? this.enums.toString() : "null";
        String lo = this.loc != 0 ? Integer.toString(this.loc) : "null";

        return "path : ".concat(p)
                .concat(", language : ").concat(l)
                .concat(", imports : ").concat(i)
                .concat(", interfaces : ").concat(in)
                .concat(", classes : ").concat(c)
                .concat(", enums : ").concat(e)
                .concat(", loc : ").concat(lo);
    }
}
