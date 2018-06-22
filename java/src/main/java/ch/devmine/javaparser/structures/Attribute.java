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
public class Attribute {

    @Getter @Setter
    private List<String> doc;

    @Getter @Setter
    private String name;

    @Getter @Setter
    private String type;

    @Getter @Setter
    private String value;

    @Getter @Setter @SerializedName(value = "is_pointer")
    private boolean isPointer;

    @Getter @Setter
    private String visibility;

    @Getter @Setter
    private boolean constant;

    @Getter @Setter @SerializedName(value = "static")
    private boolean statick;

    public Attribute() {}

    public Attribute(List<String> doc, String name, String type, String value, boolean isPointer,
                     String visibility, boolean constant, boolean statick)
    {
        this.doc = doc;
        this.name = name;
        this.type = type;
        this.value = value;
        this.isPointer = isPointer;
        this.visibility = visibility;
        this.constant = constant;
        this.statick = statick;
    }

    @Override
    public String toString() {
        String d = this.doc != null ? this.doc.toString() : "null";
        String n = this.name != null ? this.name : "null";
        String t = this.type != null ? this.type : "null";
        String v = this.value != null ? this.value : "null";
        String vi = this.visibility != null ? this.visibility : "null";
        String c = this.constant ? "true" : "false";
        String s = this.statick ? "true" : "false";

        return "doc : ".concat(d)
                .concat(", name : ").concat(n)
                .concat(", type : ").concat(t)
                .concat(", value : ").concat(v)
                .concat(", visibility : ").concat(vi)
                .concat(", constant : ").concat(c)
                .concat(", static : ").concat(s);
    }
}
