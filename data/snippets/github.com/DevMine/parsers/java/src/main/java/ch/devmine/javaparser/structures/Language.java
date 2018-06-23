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
public class Language {

    @Getter @Setter
    private String language;

    @Getter @Setter
    private List<String> paradigms;

    public Language() {}

    public Language(String lang, List<String> paradigms) {
        this.language = lang;
        this.paradigms = paradigms;
    }

    @Override
    public String toString() {
        String l = (this.language != null) ? this.language : "null";
        String p = (this.paradigms != null) ? this.paradigms.toString() : "null";

        return "language : ".concat(l)
                .concat(", paradigms : ").concat(p);
    }
}
