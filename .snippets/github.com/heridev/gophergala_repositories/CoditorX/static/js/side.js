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

var side = {
    init: function () {
        new Tabs({
            id: ".side",
            clickBefore: function (it) {
                var $side = $(".side"),
                        $main = $(".main");
                if ($(it).hasClass('current')) {
                    if ($side.width() === 0) {
                        $side.animate({
                            width: '20%'
                        }, function () {
                            $main.css({
                                left: "20%",
                                width: "80%"
                            });
                        });
                    } else {
                        $side.animate({
                            width: '0'
                        });
                        $main.css({
                            left: "41px",
                            width: ($(window).width() - 41) + "px"
                        });
                    }
                } else {
                    $side.animate({
                        width: '20%'
                    }, function () {
                        $main.css({
                            left: "20%",
                            width: "80%"
                        });
                    });
                }
            }
        });

        this._initFileList();
        this._initShareList();
        this._initDialog();
    },
    _mockData: function (count) {
        var data = [],
                types = ['html', 'htm', 'sql', 'properties', 'md', 'json', 'xml',
                    'jpg', 'jpeg', 'bmp', 'gif', 'png', 'svg', 'ico', 'go', 'css', 'js', 'txt', ''];

        for (var i = 0, max = count; i < max; i++) {
            var type = types[Math.ceil(Math.random() * 18)];
            data.push({
                'isShare': false,
                'id': i,
                'name': 'fileName' + i + '.' + type,
                'type': type
            });
        }

        return data;
    },
    _initFileList: function () {
        $("#files").parent().mouseup(function (event) {
            event.stopPropagation();

            if (event.button === 0) { // 左键
                $("#files").next().hide();
                $("#filesNewMenu").hide();
                return;
            }

            // event.button === 2 右键
            $("#files").next().hide();
            var $filesNewMenu = $("#filesNewMenu");
            $filesNewMenu.show().css({
                "left": event.clientX + "px",
                "top": (event.clientY) + "px"
            });

            return false;
        });

        var request = newRequest();

        $.ajax({
            url: "/files",
            type: "POST",
            cache: false,
            data: JSON.stringify(request),
            success: function (data) {
                if (!data.succ) {
                    $('#dialogAlert').dialog("open", data.msg);
                    return false;
                }

                var filesHTML = '<ul class="list">',
                        $files = $("#files");

                for (var i = 0, max = data.files.length; i < max; i++) {
                    filesHTML += '<li title="' + data.files[i].name + '" data-share="' + data.files[i].isShare
                            + '"><span class="ico-file ' + coditor.getClassBySuffix(data.files[i].type)
                            + '"></span> ' + data.files[i].name + '</li>';
                }
                $files.html(filesHTML + '</ul>');

                $files.find("li").mouseup(function (event) {
                    event.stopPropagation();

                    if (event.button === 0) { // 左键
                        $files.next().hide();
                        return;
                    }

                    // event.button === 2 右键
                    $("#filesNewMenu").hide();
                    $files.next().show().css({
                        "left": "38px",
                        "top": (event.target.offsetTop - $files.parent().scrollTop() + 22) + "px"
                    });
                    if ($(this).data('share')) {
                        $files.next().find('.share').hide();
                        $files.next().find('.unshare').show();
                    } else {
                        $files.next().find('.share').show();
                        $files.next().find('.unshare').hide();
                    }

                    $files.find("li").removeClass("current");
                    $(this).addClass('current');

                    return false;
                }).dblclick(function () {
                    side.open(coditor.workspace + coditor.pathSeparator + $.trim($(this).text()));
                    editor.codemirror.setOption("readOnly", false);
                });
            }
        });
    },
    _initShareList: function () {
        var request = newRequest();

        $.ajax({
            url: "/shares",
            type: "POST",
            data: JSON.stringify(request),
            success: function (data) {
                if (!data.succ) {
                    $('#dialogAlert').dialog("open", data.msg);
                    return false;
                }

                var filesHTML = '<ul class="list">',
                        $shareFiles = $("#shareFiles");

                for (var i = 0, max = data.shares.length; i < max; i++) {
                    var shareFile = data.shares[i];
                    var index = shareFile.docName.lastIndexOf("\\.");
                    var fileType = "";
                    if (index > 0) {
                        fileType = shareFile.docName.sub(index + 1);
                    }
                    filesHTML += '<li shareType="' + shareFile.shareType + '" docName="workspaces' + coditor.pathSeparator + shareFile.owner + coditor.pathSeparator + "workspace" + coditor.pathSeparator + shareFile.docName + '" title="' + '/' + shareFile.owner + '/' + shareFile.docName + '" ><span class="ico-file ' + coditor.getClassBySuffix(fileType)
                            + '"></span> ' + '/' + shareFile.owner + '/' + shareFile.docName + '</li>';
                }
                $shareFiles.html(filesHTML + '</ul>');

                $shareFiles.find("li").mouseup(function (event) {
                    event.stopPropagation();

                    if (event.button === 0) { // 左键
                        $shareFiles.next().hide();
                        return;
                    }

                    // event.button === 2 右键
                    $shareFiles.next().show().css({
                        "left": "38px",
                        "top": (event.target.offsetTop - $shareFiles.parent().scrollTop() + 22) + "px"
                    });

                    $shareFiles.find("li").removeClass("current");
                    $(this).addClass('current');

                    return;
                }).dblclick(function () {
                    side.open($.trim($(this).attr("docName")));
                    var shareType = $(this).attr("shareType");
                    if (shareType == 0) {
                        editor.codemirror.setOption("readOnly", true);
                    } else {
                        editor.codemirror.setOption("readOnly", false);
                    }
                });
                ;
            }
        });
    },
    _initDialog: function () {
        $("#dialogAlert").dialog({
            "modal": true,
            "height": 160,
            "width": 310,
            "title": 'Prompt',
            "hiddenOk": true,
            "cancelText": 'confirm',
            "afterOpen": function (msg) {
                $("#dialogAlert").html(msg);
            }
        });

        $(".dialog-prompt > input").keyup(function (event) {
            var $okBtn = $(this).closest(".dialog-main").find(".dialog-footer > button:eq(0)");
            if (event.which === 13 && !$okBtn.prop("disabled")) {
                $okBtn.click();
            }

            if ($.trim($(this).val()) === "") {
                $okBtn.prop("disabled", true);
            } else {
                $okBtn.prop("disabled", false);
            }
        });

        $("#dialogUnshareConfirm").dialog({
            "modal": true,
            "height": 36,
            "width": 260,
            "title": 'Unshare',
            "okText": 'Unshare',
            "cancelText": 'Cancel',
            "afterOpen": function () {
                $("#dialogUnshareConfirm > b").html('"' + $.trim($('#files li.current').text()) + '"');
            },
            "ok": function () {
                var request = newRequest();
                request["fileName"] = coditor.workspace + coditor.pathSeparator + $.trim($('#files li.current').text());
                request["editors"] = '';
                request["isPublic"] = 0;
                request["viewers"] = '';
                $.ajax({
                    type: 'POST',
                    url: '/share',
                    data: JSON.stringify(request),
                    async: false,
                    dataType: "json",
                    success: function (data) {
                        if (!data.succ) {
                            $('#dialogAlert').dialog("open", data.msg);
                            return false;
                        }

                        $("#files li.current").data("share", false);
                        $("#dialogUnshareConfirm").dialog('close');
                    }
                });
            }
        });


        $("#dialogRemoveConfirm").dialog({
            "modal": true,
            "height": 36,
            "width": 260,
            "title": 'Delete',
            "okText": 'Delete',
            "cancelText": 'Cancel',
            "afterOpen": function () {
                $("#dialogRemoveConfirm > b").html('"' + $.trim($('#files li.current').text()) + '"');
            },
            "ok": function () {
                var request = newRequest();
                request.name = $.trim($('#files li.current').text());

                $.ajax({
                    type: 'POST',
                    url: '/file/del',
                    data: JSON.stringify(request),
                    dataType: "json",
                    success: function (data) {
                        if (!data.succ) {
                            $('#dialogAlert').dialog("open", data.msg);
                            //return false;
                        }

                        $('#files li.current').remove();
                        $("#dialogRemoveConfirm").dialog("close");
                    }
                });
            }
        });

        $("#dialogShare").load('/share', function () {



            $("#dialogShare").dialog({
                "modal": true,
                "height": 190,
                "width": 600,
                "title": 'Share',
                "afterOpen": function () {
                    $("#dialogShare").find('input[type=text]').val('');
                    $("#dialogShare .fileName").val(coditor.workspace + coditor.pathSeparator + $.trim($("#files li.current").text()));
                    $("#dialogShare").find('input[type=checkbox]').prop('checked', false);
                    $("#dialogShare").find('.viewers').show();
                    var owner = coditor.sessionUsername;
                    var fileName = $("#files li.current").attr("title");
                    var publicViewUrl = "http://coditorx.b3log.org/" + owner + "/doc/" + fileName;
                    $("#dialogShare").find(".publicView").html(publicViewUrl);
                    $("#dialogShare").find(".publicView").attr("href", publicViewUrl);
                },
                "ok": function () {
                    var fileName = $("#dialogShare .fileName").val(),
                            editors = $("#dialogShare .editors").val(),
                            isPublic = 0,
                            viewers = '';
                    if ($("#dialogShare .isPublic").prop("checked")) {
                        isPublic = 1;
                    } else {
                        viewers = $("#dialogShare .viewersInp").val();
                    }

                    var request = newRequest();
                    request["fileName"] = fileName;
                    request["editors"] = editors;
                    request["isPublic"] = isPublic;
                    request["viewers"] = viewers;

                    $.ajax({
                        type: 'POST',
                        url: '/share',
                        data: JSON.stringify(request),
                        async: false,
                        dataType: "json",
                        success: function (data) {
                            if (!data.succ) {

                                $('#dialogAlert').dialog("open", data.msg);
                                $("#dialogShare").dialog('close');
                                return false;
                            }
                            $("#files li.current").data("share", true);
                            $("#dialogShare").dialog('close');
                        }
                    });
                }
            });

            $("#dialogShare").find('input[type=checkbox]').click(function () {
                if ($(this).prop('checked')) {
                    $("#dialogShare").find('.viewers').hide();
                } else {
                    $("#dialogShare").find('.viewers').show();
                }
            });
        });

        $("#dialogRenamePrompt").dialog({
            "modal": true,
            "height": 52,
            "width": 260,
            "title": 'Rename',
            "okText": 'Rename',
            "cancelText": 'Cancel',
            "afterOpen": function () {
                $("#dialogRenamePrompt").closest(".dialog-main").find(".dialog-footer > button:eq(0)").prop("disabled", true);
                $("#dialogRenamePrompt > input").val($.trim($('#files li.current').text())).select().focus();
            },
            "ok": function () {
                var name = $("#dialogRenamePrompt > input").val(),
                        request = newRequest();

                request.newName = name;
                request.oldName = $.trim($('#files li.current').text());

                $.ajax({
                    type: 'POST',
                    url: '/file/rename',
                    data: JSON.stringify(request),
                    dataType: "json",
                    success: function (data) {
                        if (!data.succ) {
                            $('#dialogAlert').dialog("open", data.msg);
                            $("#dialogRenamePrompt").dialog("close");
                            return false;
                        }

                        $("#dialogRenamePrompt").dialog("close");
                        side._initFileList();
                    }
                });
            }
        });

        $("#dialogNewFilePrompt").dialog({
            "modal": true,
            "height": 52,
            "width": 260,
            "title": 'Create File',
            "okText": 'Create',
            "cancelText": 'Cancel',
            "afterOpen": function () {
                $("#dialogNewFilePrompt > input").val('').focus();
                $("#dialogNewFilePrompt").closest(".dialog-main").find(".dialog-footer > button:eq(0)").prop("disabled", true);
            },
            "ok": function () {
                var request = newRequest(),
                        name = $("#dialogNewFilePrompt > input").val();
                request.name = name;

                var isOk = false;
                $.ajax({
                    async: false,
                    type: 'POST',
                    url: '/file/new',
                    data: JSON.stringify(request),
                    dataType: "json",
                    success: function (data) {
                        if (data.succ) {
                            side._initFileList();
                            isOk = true;
                        } else {
                            $('#dialogAlert').dialog("open", data.msg);
                            $("#dialogNewFilePrompt").dialog("close");
                        }
                    }
                });
                return isOk;
            }
        });
    },
    new : function () {
        $("#dialogNewFilePrompt").dialog("open");
    },
    remove: function () {
        $("#dialogRemoveConfirm").dialog("open");
    },
    share: function () {
        $("#dialogShare").dialog('open');
    },
    unshare: function () {
        $("#dialogUnshareConfirm").dialog("open");
    },
    rename: function () {
        $("#dialogRenamePrompt").dialog('open');
    },
    open: function (fileName) {
        $('.content').show();
        $('.welcome').hide();

        var $editor = $("#editor");
        editor.codemirror = CodeMirror.fromTextArea($editor[0], {
            autofocus: true,
            lineNumbers: true,
            theme: "3024-night"
        });

        var mode = CodeMirror.findModeByFileName(fileName);
        if (mode) {
            editor.codemirror.setOption("mode", mode.mode);
        }

        editor.codemirror.setSize('50%', $(".main").height() - $(".menu").height());
        editor.codemirror.on('changes', function (cm, changes) {
            if (changes && changes[0] && "setValue" === changes[0].origin) {
                return;
            }

            var request = newRequest();
            request.cmd = "commit";
            request.content = cm.getValue();
            request.docName = fileName;
            request.user = coditor.sessionUsername;
            request.cursor = cm.getCursor();
            request.color = coditor.color;

            editor.channel.send(JSON.stringify(request));

            $('.preview').html(markdown.toHTML(editor.codemirror.getValue()));
        });

        var request = newRequest();
        request.fileName = fileName;

        $.ajax({
            async: false,
            url: "/doc/open",
            type: "POST",
            data: JSON.stringify(request),
            success: function (data) {
                if (!data.succ) {
                    $('#dialogAlert').dialog("open", data.msg);
                    return false;
                }
                editor.codemirror.doc.setValue(data.doc.content);
                editor.currentFileName = fileName;

                $('.preview').html(markdown.toHTML(data.doc.content));
            }
        });

        var request = newRequest();
        request.docName = fileName;
        request.offset = 0;
        request.color = coditor.color;

        $.ajax({
            async: false,
            url: "/doc/setCursor",
            type: "POST",
            data: JSON.stringify(request),
            success: function (data) {
                if (!data.succ) {
                    $('#dialogAlert').dialog("open", data.msg);
                    return false;
                }
            }
        });

        menu.listCursors();
    },
    getDocURL: function () {
        var url = 'http://coditorx.b3log.org/' + userName + '/doc/' + $.trim($('#files li.current').text());
        $('#dialogAlert').dialog('open', '<a href="' + url + '" target="_blank">' + url + '</a>');
    }
};