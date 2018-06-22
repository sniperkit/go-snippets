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

var coditor = {
    conf: undefined,
    sessionId: undefined,
    sessionUsername: undefined,
    color: undefined,
    workspace: undefined,
    i18n: undefined,
    pathSeparator: undefined,
    init: function () {
        var height = $(".main").height() - $(".menu").height();
        $(".preview").height(height - 20);
        $(".welcome").height(height);
        $(".welcome").css('padding-top', (height - 80) / 2 + 'px');

        $(window).resize(function () {
            var height = $(".main").height() - $(".menu").height();
            if (editor.codemirror) {
                editor.codemirror.setSize('50%', height);
            }
            $(".preview").height(height - 20);
            $(".welcome").height(height);
            $(".welcome").css('padding-top', (height - 80) / 2 + 'px');
        });

        // 点击隐藏弹出层
        $("body").bind("mouseup", function (event) {
            $(".frame").hide();
        });

        // 禁止鼠标右键菜单
        document.oncontextmenu = function () {
            return false;
        };

        this.conf = conf;
        this.sessionId = sessionId;
        this.sessionUsername = sessionUsername;
        this.color = color;
        this.i18n = i18n;
        this.workspace = workspace;
        this.pathSeparator = pathSeparator;
    },
    getClassBySuffix: function (suffix) {
        var iconSkin = "ico-file-other";
        switch (suffix) {
            case "html":
            case "htm":
                iconSkin = "ico-file-html";
                break;
            case "go":
                iconSkin = "ico-file-go";
                break;
            case "css":
                iconSkin = "ico-file-css";
                break;
            case "txt":
                iconSkin = "ico-file-text";
                break;
            case "sql":
                iconSkin = "ico-file-sql";
                break;
            case "properties":
                iconSkin = "ico-file-pro";
                break;
            case "md":
                iconSkin = "ico-file-md";
                break;
            case "js", "json":
                iconSkin = "ico-file-js";
                break;
            case "xml":
                iconSkin = "ico-file-xml";
                break;
            case "jpg":
            case "jpeg":
            case "bmp":
            case "gif":
            case "png":
            case "svg":
            case "ico":
                iconSkin = "ico-file-img";
                break;
        }

        return iconSkin;
    }
};

$(document).ready(function () {
    menu.init();
    side.init();
    coditor.init();
    session.init();
    notification.init();
    editor.init();
});