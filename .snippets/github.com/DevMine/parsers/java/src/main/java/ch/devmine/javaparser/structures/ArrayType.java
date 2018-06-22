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
public class ArrayType {

    @Getter @Setter
    private List<Integer> dimensions;

    @Getter @Setter @SerializedName(value = "element_type")
    private Expr elementType;

    public ArrayType() {}

    public ArrayType(List<Integer> dimensions, Expr elementType) {
        this.dimensions = dimensions;
        this.elementType = elementType;
    }

    @Override
    public String toString() {
        String l  = this.dimensions != null ? this.dimensions.toString() : "null";
        String e = this.elementType != null ? this.elementType.toString() : "null";

        return "length : ".concat(l)
                .concat(", element type : ").concat(e);
    }
}
