/*
 * Copyright (c) 2015, b3log.org
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

function Doc(name, content, version, debug) {
    this.name = name;
    this.content = content;
    this.version = version;
    //this.dmp = new diff_match_patch();
    if (debug) {
        this.debug = true;
    } else {
        this.debug = false;
    }
    return this;
}

function OpenDoc(fileName) {
    var doc = undefined;

    var request = newRequest();
    request.fileName = fileName;

    $.ajax({
        async: false,
        url: "/doc/open",
        type: "POST",
        data: JSON.stringify(request),
        success: function (data) {
            if (data.succ) {
                var d = data.doc;
                doc = new Doc(fileName, d.content, d.version);
            } else {
                $('#dialogAlert').dialog("open", data.msg);
            }
        },
        error: function (XMLHttpRequest, textStatus, errorThrown) {
            // TODO
        }
    });

    return doc;
}

// commit doc to server.
Doc.prototype.commit = function () {
    var doc = this;
    var result = undefined;
    var request = newRequest();
    var file = {
        name: this.name,
        version: this.version,
        content: this.content
    };
    request.file = file;

    $.ajax({
        async: false,
        url: "/doc/commit",
        type: "POST",
        data: JSON.stringify(request),
        success: function (data) {
            result = data;
            if (data.succ) {
                doc.version = data.output.version;
            } else {
                $('#dialogAlert').dialog("open", data.msg);
            }
        },
        error: function (XMLHttpRequest, textStatus, errorThrown) {
            // TODO
        }
    });
    return result;
};

// pull from server.
Doc.prototype.pull = function () {
    var doc = this;
    var result = undefined;
    var request = newRequest();
    var file = {
        name: doc.name,
        version: doc.version
    };
    request.file = file;

    $.ajax({
        async: false,
        url: "/doc/fetch",
        type: "POST",
        data: JSON.stringify(request),
        success: function (data) {
            result = data;
            if (data.succ) {
                var length = data.patchss.length;
                for (var i = 0; i < length; i++) {
                    var patchsStr = data.patchss[i];
                    var patches = doc.dmp.patch_fromText(patchsStr);
                    var outputs = doc.dmp.patch_apply(patches, doc.content);
                    var result = outputs[1];
                    console.log(patches);
                    for (var i = 0; i < result.length; i++) {
                        if (!result[i]) {
                            console.log("result:" + result);
                        }
                    }
                    doc.content = outputs[0];
                }
                doc.version = data.version;
            } else {
                $('#dialogAlert').dialog("open", data.msg);
            }
        },
        error: function (XMLHttpRequest, textStatus, errorThrown) {
            // TODO
        }
    });
    return result;
};

// set content.
Doc.prototype.setContent = function (content) {
    this.content = content;
};