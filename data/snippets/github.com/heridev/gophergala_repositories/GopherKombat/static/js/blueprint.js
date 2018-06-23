/*global $*/
/*global document*/
/*global GK*/
/*global ace*/
/*global console*/

$(document).ready(function() {
    var editor = ace.edit("editor");
    editor.setTheme("ace/theme/tomorrow_night_eighties");
    editor.getSession().setMode("ace/mode/golang");
    $(".submitCode").click(function() {
        var code = editor.getSession().getValue(),
            data,
            url = "/blueprint/submit";
        data = {
            code: code
        };
        GK.requestAgent().doPOST(url, data, function(resp) {
            console.log(resp);
            $(".validation").html(resp.message);
            $(".validation").show();
        });
        
    });
    GK.requestAgent().doGET("/blueprint/get", function(resp) {
        var code = resp.code;
        editor.setValue(code);
        $(".message").html(resp.message);
    });
});