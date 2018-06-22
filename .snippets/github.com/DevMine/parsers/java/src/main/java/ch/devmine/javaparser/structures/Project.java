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
public class Project {

    @Getter @Setter
    private String name;

    @Getter @Setter
    private List<Language> languages;

    @Getter @Setter
    private List<Package> packages;

    @Getter @Setter
    private int loc;

    public Project() {}

    public Project(String name, List<Language> languages, List<Package> packages, int loc) {
        this.name = name;
        this.languages = languages;
        this.packages = packages;
        this.loc = loc;
    }

    @Override
    public String toString() {
        String n = this.name != null ? this.name : "null";
        String l = this.languages != null ? this.languages.toString() : "null";
        String p = this.packages != null ? this.packages.toString() : "null";
        String lo = this.loc != 0 ? Integer.toString(this.loc) : "null";

        return "name : ".concat(n)
                .concat(", languages : ").concat(l)
                .concat("packages : ").concat(p)
                .concat(", loc : ").concat(lo);

    }

}
