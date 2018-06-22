/*
 * Copyright 2014-2015 The DevMine Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ch.devmine.javaparser.structures;


import com.google.gson.annotations.SerializedName;
import java.util.ArrayList;
import java.util.List;
import lombok.Getter;
import lombok.Setter;

/**
 *
 * @author lweingart
 */
public class Package {

    @Getter @Setter
    private List<String> doc;

    @Getter @Setter
    private String name;

    @Getter @Setter
    private String path;

    @Getter @Setter @SerializedName(value = "source_files")
    private List<SourceFile> sourceFiles;

    @Getter @Setter
    private int loc;

    public Package() {
        this.sourceFiles = new ArrayList<>();
    }

    public Package(List<String> doc, String name, String path, List<SourceFile> sourceFiles, int loc) {
        this.doc = doc;
        this.name = name;
        this.path = path;
        this.sourceFiles = sourceFiles;
        this.loc = loc;
    }

    @Override
    public String toString() {
        String d = this.doc != null ? this.doc.toString() : "null";
        String n = this.name != null ? this.name : "null";
        String p = this.path != null ? this.path : "null";
        String s = this.sourceFiles != null ? this.sourceFiles.toString() : "null";

        return "doc : ".concat(d)
                .concat(", name : ").concat(n)
                .concat(", path : ").concat(p)
                .concat(", source files : ").concat(s);
    }
}
