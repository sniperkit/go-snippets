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

var menu = {
    init: function () {
        this._initShare();
        this._initCurrentEditors();
    },
    _initShare: function () {
        $(".menu .ico-share").hover(function () {
            $(".menu .share-panel").show();
        });

        $(".share-panel .font-ico").click(function () {
            var key = $(this).attr('class').split('-')[2];
            var url = "https://github.com/gophergala/CoditorX",
                    pic = 'http://coditorx.b3log.org/static/images/logo-bg.png';
            var urls = {};
            urls.email = "mailto:?subject=" + $('title').text()
                    + "&body=" + $('meta[name=description]').attr('content') + ' ' + url;

            var twitterShare = encodeURIComponent($('meta[name=description]').attr('content') + " " + url + " #golang");
            urls.twitter = "https://twitter.com/intent/tweet?status=" + twitterShare;

            urls.facebook = "https://www.facebook.com/sharer/sharer.php?u=" + url;
            urls.googleplus = "https://plus.google.com/share?url=" + url;

            var title = encodeURIComponent($('title').text() + '. \n' + $('meta[name=description]').attr('content')
                    + " #golang#");
            urls.weibo = "http://v.t.sina.com.cn/share/share.php?title=" + title + "&url=" + url + "&pic=" + pic;
            urls.tencent = "http://share.v.t.qq.com/index.php?c=share&a=index&title=" + title +
                    "&url=" + url + "&pic=" + pic;

            window.open(urls[key], "_blank", "top=100,left=200,width=648,height=618");
        });
    },
    exit: function () {
        var request = newRequest();
        $.ajax({
            type: 'POST',
            url: '/logout',
            data: JSON.stringify(request),
            dataType: "json",
            success: function (data) {
                if (data.succ) {
                    window.location.href = "/";
                } else {
                    $('#dialogAlert').dialog("open", data.msg);
                }
            }
        });
    },
    _initCurrentEditors: function () {
        setInterval(menu.listCursors, 30000);
        menu.listCursors()
    },
    listCursors: function () {
        var request = newRequest();
        if (!editor.currentFileName) {
            return;
        }

        var docName = editor.currentFileName;
        request.docName = docName;
        $.ajax({
            type: 'POST',
            url: '/doc/listCursors',
            data: JSON.stringify(request),
            dataType: "json",
            success: function (data) {
                if (!data.succ) {
                    $('#dialogAlert').dialog("open", data.msg);
                    return false;
                }
                var imgStr = "";
                for (var i = 0; i < data.cursors.length; i++) {
                    var cursor = data.cursors[i];
                    imgStr += '<img class="gravatar" onerror="this.src=\'/static/images/user-thumbnail.png\'" src="https://secure.gravatar.com/avatar/' + cursor.md5Email + '?s=30&d=https://symphony.b3log.org/images/user-thumbnail.png" title="' + cursor.username + '"/>';
                }
                $(".fn-left").html(imgStr);
            }
        });
    }
};