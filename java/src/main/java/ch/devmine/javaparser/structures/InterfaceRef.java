/*
 * Copyright 2014-2015 The DevMine Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ch.devmine.javaparser.structures;

import com.google.gson.annotations.SerializedName;
import lombok.Getter;
import lombok.Setter;

/**
 *
 * @author lweingart
 */
public class InterfaceRef {

    @Getter @Setter
    private String namespace;

    @Getter @Setter @SerializedName(value = "interface_name")
    private String interfaceName;

    public InterfaceRef() {}

    public InterfaceRef(String namespace, String interfaceName) {
        this.namespace = namespace;
        this.interfaceName = interfaceName;
    }

    @Override
    public String toString() {
        String n = (this.namespace != null) ? this.namespace : "null";
        String i = (this.interfaceName != null) ? this.interfaceName : "null";

        return "namespace : ".concat(n)
                .concat(", interface name : ").concat(i);
    }
}
