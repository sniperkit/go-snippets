"use strict";

// The trigger for this function is defined at the bottom of this script
function load() {
    // Get the locally stored settings.  If nothing is stored
    // locally, the current values that are defined by sdes.config
    // will be used.  Take a look at the js/config.js file to see
    // the default sdes.config values.
    chrome.storage.local.get(sdes.config, update);

    var saveButton    = document.getElementById("save-button"),
        saveErrorBody = document.getElementById("save-error"),
        newRuleButton = document.getElementById("new-rule-button"),
        restoreLink   = document.getElementById("restore"),
        xFrameOptions = ["DENY", "SAMEORIGIN"],
        htmlUtil      = new sdes.utils.html();

    var hostTypes = [
        { id: "", label: ""}, 
        { id: "github", label: "github.com", api: "https://api.github.com" },
        { id: "gitlab", label: "gitlab.com", api: "https://gitlab.com/api/v3" },
        { id: "github-ent", label: "Github Enterprise", api: "BASE_URL/api/v3" },
        { id: "gitlab-le", label: "GitLab CE/EE", api: "BASE_URL/api/v3" }
    ];

    var hostToInfo = mapHostToInfo();

    restoreLink.onclick = restoreDefault;

    function mapHostToInfo() {
        var hostToInfo = {};

        for ( var i = 0; i < hostTypes.length; i++ ) {
            var info = hostTypes[i];

            hostToInfo[info.id] = info;
        }

        return hostToInfo;
    }

    function update(localConfig) {
        var idToInput    = {},
            ruleToInputs = {};

        for ( var key in localConfig ) {
            if ( key === "routingRules" )
                renderRules(key, localConfig[key]);
            else
                set(key, localConfig[key]);
        }

        function renderRules(renderTo, rules) {
            var renderToBody = document.getElementById(renderTo);

            if ( renderToBody === null )
                throw("No element with the id '"+renderTo+"' found");

            $(renderToBody).html("");

            for ( var i = 0; i < rules.length; i++ )
                addRule("rule-"+i, rules[i], false);

            setupEvents();
            check();

            newRuleButton.onclick = clickedAdd; 

            function addRule(id, rule, expand) {
                var direction = expand ? "down" : "right",
                    display   = expand ? "" : "display:none",

                    html =
                        "<div class=row>"+
                            "<input id="+id+"-matches type=text "+
                                "placeholder='URL pattern' "+
                                "value='"+rule.matches+"' "+
                                "style='width:480px;'>"+
                            "<span class='options'>"+
                                "<span id="+id+"-option-down "+
                                    "class='octicon octicon-arrow-down option' "+
                                    "style='margin-right:10px;font-size:20px;'></span>"+
                                "<span id="+id+"-option-up "+
                                    "class='octicon octicon-arrow-up option' "+
                                    "style='margin-right:10px;font-size:20px;'></span>"+
                                "<span id="+id+"-option-delete "+
                                    "class='octicon octicon-x option' "+
                                    "style='font-weight:bold;font-size:17px;'></span>"+
                            "</sapn>"+
                        "</div>"+
                        "<div style='clear:both'></div>"+
                        "<div class='block block-header'>"+
                            "<span id="+id+"-host-triangle "+
                                "class='triangle octicon octicon-triangle-"+direction+"'></span>"+
                            "<span id="+id+"-host-title "+
                                "class='block-name'>Settings</span>"+
                        "</div>"+
                        "<div class='block' style='"+display+"'>"+
                            getHostSettings()+
                            getGitSenseSettings()+
                        "</div>";

                var div = 
                    htmlUtil.createDiv({
                        id: id,
                        html: html,
                        cls: "rule"
                    });

                renderToBody.appendChild(div);

                renderToBody.style.borderBottom = "1px solid #666";

                function getHostTypes() {
                    var html = "<select id="+id+"-host-type>";

                    for ( var i = 0; i < hostTypes.length; i++ ) {
                        var type     = hostTypes[i],
                            selected = type.id === rule.host.type ? "selected" : "";

                        html += "<option value=\""+type.id+"\" "+selected+">"+type.label+"</option>";
                    }

                    html += "</select>";

                    return html;
                }

                function getXFrameOptions() {
                    var html = "<select id="+id+"-host-xFrameOptions style='width:112px'>";

                    for ( var i = 0; i < xFrameOptions.length; i++ ) {
                        var option   = xFrameOptions[i],
                            selected = option === rule.host.xFrameOptions ? "selected" : "";

                        html += "<option "+selected+">"+option+"</option>";
                    }

                    html += "</select>";

                    return html;
                }

                function getHostSettings() {
                    var baseUrl = rule.host.baseUrl;

                    var html =
                        "<div class='block-header'>"+
                            "<span id="+id+"-host-title "+
                                "style='font-weight:bold;'>"+
                                "github.com, gitlab.com, GitHub Enterprise, and GitLab CE/EE"+
                            "</span>"+
                        "</div>"+
                        "<div class='block-body'>"+
                            "<table>"+
                                "<tr>"+
                                    "<td class=field-cell>Type</td>"+
                                    "<td class=field-cell>Base URL</td>"+
                                    "<td class=field-cell>X-Frame-Options</td>"+
                                    "<td class=field-cell style='display:none;'>Access token</td>"+
                                "</tr>"+
                                "<tr>"+
                                    "<td class=field-cell>"+getHostTypes()+"</td>"+
                                    "<td class=field-cell>"+
                                        "<input type=text "+
                                            "id="+id+"-host-baseUrl "+
                                            "value='"+baseUrl+"' "+
                                            "style='width:270px'>"+
                                    "</td>"+
                                    "<td class=field-cell>"+getXFrameOptions()+"</td>"+
                                    "<td class=field-cell style='padding-right:0px;'>"+
                                        "<input type=password "+
                                            "id="+id+"-host-secret "+
                                            "value='"+rule.host.secret+"' "+
                                            "placeholder=Optional "+
                                            "style='width:213px;display:none;'>"+
                                    "</td>"+
                                "</tr>"+
                            "</table>"+
                        "</div>";

                    return html;
                }

                function getGitSenseSettings() {
                    var html =
                        "<div class='block-header'>"+
                            "<span id="+id+"-gitsense-title "+
                                "style='font-weight:bold;'>GitSense</span>"+
                        "</div>"+
                        "<div class='block-body'>"+
                            "<table>"+
                                "<tr>"+
                                    "<td class=field-cell>Base URL</td>"+
                                    "<td class=field-cell style='display:none;'>Access token</td>"+
                                "</tr>"+
                                "<tr>"+
                                    "<td class=field-cell>"+
                                        "<input type=text "+
                                            "id="+id+"-gitsense-baseUrl "+
                                            "value='"+rule.gitsense.baseUrl+"' "+
                                            "style='width:550px'>"+
                                    "</td>"+
                                    "<td class=field-cell style='padding-right:0px;'>"+
                                        "<input type=password "+
                                            "id="+id+"-gitsense-secret "+
                                            "value='"+rule.gitsense.secret+"' "+
                                            "placeholder=Optional "+
                                            "style='width:350px;display:none;'>"+
                                    "</td>"+
                                "</tr>"+
                            "</table>"+
                            //"<div style='padding-top:10px;padding-bottom:5px;'>"+
                            //    "Commit decorator"+
                            //"</div>"+
                            //"<input type=text "+
                            //    "id="+id+"-gitsense-commitDecorator "+
                            //    "value='"+rule.gitsense.commitDecorator+"' "+
                            //    "placeholder=Optional "+
                            //    "style='width:809px;padding-right:0px;'>"+
                        "</div>";
                    return html;
                }
            }

            function setupEvents() {
                var elems   = document.getElementsByClassName("block"),
                    options = ["up", "down", "delete"];

                for ( var i = 0; i < elems.length; i++ ) {
                    var header = elems[i++];
                    var body   = elems[i];

                    setupBlock(header, body);
                }

                for ( var i = 0; i < rules.length; i++ ) {
                    var rule = rules[i],
                        id   = "rule-"+i;

                    if ( ruleToInputs[id] !== undefined )
                        continue;

                    ruleToInputs[id] = getInputs(id, rule);

                    for ( var j = 0; j < options.length; j++ ) {
                        var thisId = id+"-option-"+options[j];

                        document.getElementById(thisId).onclick = clickedOption;
                    }
                }

                function getInputs(id, rule) {
                    var inputs = {};

                    for ( var key in rule ) {
                        switch( key ) {
                            case "matches":
                                var thisId = id+"-"+key,
                                    input  = document.getElementById(thisId);

                                idToInput[thisId] = input;

                                inputs[key]        = input;
                                input.onkeyup      = check;
                                input.currentValue = rule[key];
                                
                                continue;
                            case "host":
                                for ( var hostKey in rule[key] ) {

                                    // This values gets dynamically defined when the rules
                                    // are saved
                                    if ( hostKey === "api" )
                                        continue;

                                    var hostId = id+"-host-"+hostKey,
                                        input  = document.getElementById(hostId);

                                    inputs[key+"-"+hostKey] = input;

                                    idToInput[hostId] = input;

                                    if ( hostKey === "type" || hostKey === "xFrameOptions" )
                                        input.onchange = check;
                                    else
                                        input.onkeyup = check;

                                    input.currentValue = rule[key][hostKey];
                                }

                                continue;
                            case "gitsense":
                                for ( var gitsenseKey in rule[key] ) {
                                    var gitsenseId = id+"-gitsense-"+gitsenseKey,
                                        input      = document.getElementById(gitsenseId);

                                    idToInput[gitsenseId] = input;

                                    inputs[key+"-"+gitsenseKey] = input;

                                    input.onkeyup      = check;
                                    input.currentValue = rule[key][gitsenseKey];
                                }

                                break;
                            default:
                                throw("Unrecognized key '"+key+"'");
                        }
                    }
                    
                    return inputs;
                }

                function setupBlock(header, body) {
                    var triangle = header.children[0],
                        title    = header.children[1];

                    title.onclick = function() {
                        if ( triangle.className.match(/down/) ) {
                            body.style.display = "none";
                            triangle.setAttribute("class", triangle.className.replace(/down/, "right"));
                        } else {
                            body.style.display = "block";
                            triangle.setAttribute("class", triangle.className.replace(/right/, "down"));
                        }
                    }
                }

                function clickedOption() {
                    var temp     = this.id.split("-"),
                        thisIdx  = parseInt(temp[1]),
                        thisId   = temp[0]+"-"+temp[1],
                        thisBody = document.getElementById(thisId),
                        action   = temp.pop(),
                        rule     = rules[thisIdx],
                        bodies   = document.getElementsByClassName("rule");

                    if ( action === "delete" ) {
                        if ( bodies.length === 1 )
                            renderToBody.style.borderBottom = "0px";

                        rule.deleted = true;
                        thisBody.parentNode.removeChild(thisBody);
                        check();

                        return;
                    }

                    if ( bodies.length === 1 )
                        return;

                    var pos;

                    for ( var i = 0; i < bodies.length; i++ ) {
                        var body = bodies[i];
                       
                        if ( thisId === body.id )
                            pos = i; 
                    }

                    var above = pos === 0 ? null : bodies[pos - 1],
                        below = pos === bodies.length - 1 ? null : bodies[pos + 1],
                        me    = $(thisBody).detach();

                    switch (action) {
                        case "up":
                            if ( above === null ) 
                                me.appendTo(renderToBody);
                            else 
                                $(me).insertBefore(above);

                            break;
                        case "down":
                            if ( below === null ) 
                                $(me).insertBefore(above);
                            else
                                me.appendTo(renderToBody);

                            break;
                        default:
                            throw("Unrecognized action '"+action+"'");
                    }

                    check();
                }
            }

            function set(key, value) {
                var input = document.getElementById(key+"-input");

                // input == null means the value that is defined by
                // sdes.config.<key> can't be updated
               
                if ( input === null )
                    return;

                input.value = 
                    value.length === 0 ? 
                        "" : 
                        JSON.stringify(value, null, 2);

                input.onkeyup = check;

                if ( updatedTextArea )
                    return;

                input.autocorrect = input.autocomplete = input.autocapitalize = input.spellcheck = false; 
            }

            function check() {
                for ( var i = 0; i < rules.length; i++ ) {
                    var id     = "rule-"+i,
                        rule   = rules[i];

                    if ( rule.deleted ) {
                        enableSave();
                        return;
                    }

                    for ( var key in rule ) {
                        switch( key ) {
                            case "matches":
                                var thisId = id+"-"+key,
                                    input  = idToInput[thisId];

                                if ( input.currentValue === clean(input.value) )
                                    continue;

                                enableSave();
                                return;
                            case "host":
                                for ( var hostKey in rule[key] ) {
                                    if ( hostKey === "api" )
                                        continue;

                                    var hostId = id+"-host-"+hostKey,
                                        input  = idToInput[hostId];

                                    if ( input.currentValue === clean(input.value) )
                                        continue;

                                    enableSave();
                                    return;
                                }

                                continue;
                            case "gitsense":
                                for ( var gitsenseKey in rule[key] ) {
                                    var gitsenseId = id+"-gitsense-"+gitsenseKey,
                                        input      = idToInput[gitsenseId];

                                    if ( input.currentValue === clean(input.value) )
                                        continue;

                                    enableSave();
                                    return;
                                } 

                                continue;
                            default:
                                throw("Unrecognized key '"+key+"'");
                        }
                    }
                }
   
                // See if the order has changed
                var bodies = document.getElementsByClassName("rule");

                for ( var i = 0; i < bodies.length; i++ ) {
                    var idx = parseInt(bodies[i].id.split("-")[1]);

                    if ( i === idx )
                        continue;

                    enableSave();
                    return;
                }

                saveButton.style.cursor = "default";
                saveButton.disabled = true;
                saveButton.onclick = null;
                saveButton.onclick  = clickedSave;

                function enableSave() {
                    saveButton.style.cursor = "pointer";
                    saveButton.disabled = false;
                    saveButton.onclick  = clickedSave;
                }
            }

            function clickedSave() {
                $(saveErrorBody).hide();

                var newRules = getPageRules();

                if ( newRules === null )
                    return;

                for ( var i = 0; i < newRules.length; i++ ) {
                    var rule = newRules[i],
                        host = rule.host;

                    rule.gitsense.baseUrl = new URL(rule.gitsense.baseUrl).origin;
                    
                    host.baseUrl = new URL(host.baseUrl).origin;
                    host.api     = hostToInfo[host.type].api.replace(/BASE_URL/, host.baseUrl);
                }

                var newConfig = { routingRules: newRules };

                chrome.permissions.getAll(function(all) {
                    var newOrigins  = mapOrigins(newRules),
                        currOrigins = {},
                        currPerms   = {},
                        addOrigins  = [],
                        rmOrigins   = [];

                    for ( var i = 0; i < all.origins.length; i++ ) {
                        var origin = all.origins[i];

                        currOrigins[origin] = "";
                    }

                    for ( var i = 0; i < all.origins.permissions; i++ ) {
                        var perm = all.permissions[i];

                        if ( perm === "storage" )
                            continue;

                        currPerms[perm] = "";
                    }

                    for ( var origin in newOrigins ) {
                        if ( currOrigins[origin] === undefined )
                            addOrigins.push(origin);
                    }

                    if ( addOrigins.length === 0 ) {
                        save();
                        return;
                    }

                    if ( addOrigins.length !== 0 ) {
                        chrome.permissions.request(
                            { origins: addOrigins },
                            function(granted) {
                                if ( granted )
                                    save();
                                else
                                    return;
                            }
                        );
                    }
                });

                function save() {
                    var statusBody = document.getElementById("save-status");

                    statusBody.textContent = "Saving...";

                    chrome.storage.local.set(
                        newConfig,
                        function() {
                            statusBody.textContent = "Saved";
            
                            setTimeout(
                                function() { 
                                    statusBody.textContent   = ""; 
                                    statusBody.style.display = "none";
                                    //saveButton.disabled      = true;
                                    update(newConfig);
                                }, 
                                750
                            ); 
                        }
                    );
                }

                function getPageRules() {
                    var bodies   = document.getElementsByClassName("rule"),
                        newRules = [],
                        errors   = [];

                    for ( var i = 0; i < bodies.length; i++ ) {
                        var body    = bodies[i],
                            id      = body.id,
                            idx     = parseInt(id.split("-")[1]),
                            rule    = rules[idx],
                            ruleNum = i+1,
                            newRule = {};

                        for ( var key in rule ) {
                            switch( key ) {
                                case "matches":
                                    var thisId = id+"-"+key,
                                        input  = idToInput[thisId],
                                        value  = clean(input.value);

                                    newRule[key] = value;

                                    if ( value === "" ) {
                                        showError(input);
                                        errors.push(
                                            "Missing required URL match pattern "+
                                            "in rule "+ruleNum
                                        );

                                        continue;
                                    }

                                    if ( ! value.match(/^https*:\/\//) ) {
                                        showError(input);
                                        errors.push("Invalid matching URL in rule "+ruleNum);
                                        continue;
                                    }
                                
                                    showGood(input);
                                    continue;
                                case "host":
                                    newRule[key] = {};

                                    for ( var hostKey in rule[key] ) {
                                        if ( hostKey === "api" )
                                            continue;

                                        var hostId = id+"-host-"+hostKey,
                                            input  = idToInput[hostId],
                                            value  = clean(input.value);

                                        newRule[key][hostKey] = value;

                                        if ( 
                                            value === "" &&
                                            (hostKey === "type" || hostKey === "baseUrl")
                                        ) {
                                            showError(input);

                                            errors.push(
                                                "Missing required "+key+" "+hostKey+" value "+
                                                "in rule "+ruleNum
                                            );

                                            continue;
                                        }

                                        if ( hostKey === "baseUrl" && ! value.match(/^https*:\/\//) ) {
                                            showError(input);

                                            errors.push(
                                                "Invalid "+key+" base URL in rule "+ruleNum
                                            );

                                            continue;
                                        }

                                        showGood(input);
                                    }

                                    continue;
                                case "gitsense":
                                    newRule[key] = {};

                                    for ( var gitsenseKey in rule[key] ) {
                                        var gitsenseId = id+"-gitsense-"+gitsenseKey,
                                            input      = idToInput[gitsenseId],
                                            value      = clean(input.value);

                                        newRule[key][gitsenseKey] = value;

                                        if ( 
                                            value === "" &&
                                            (gitsenseKey === "baseUrl")
                                        ) {
                                            showError(input);

                                            errors.push(
                                                "Missing required "+key+" "+gitsenseKey+" value "+
                                                "in rule "+ruleNum
                                            ); 

                                            continue;
                                        }

                                        if ( gitsenseKey === "baseUrl" && ! value.match(/^https*:\/\//) ) {
                                            showError(input);

                                            errors.push(
                                                "Invalid "+key+" base URL in rule "+ruleNum
                                            );

                                            continue;
                                        }

                                        showGood(input);
                                    } 

                                    continue;
                                default:
                                    throw("Unrecognized key '"+key+"'");
                            }
                        }

                        newRules.push(newRule);
                    }

                    if ( errors.length === 0 )
                        return newRules;

                    showErrorMsg(errors);
                    return null;
                }

                function showError(input, message) {
                    input.style.border = "1px solid red";
                    input.style.backgroundColor = "rgb(255, 234, 234)";
                }

                function showGood(input) {
                    input.style.border = null;
                    input.style.backgroundColor = null;
                }

                function showErrorMsg(errors) {
                    $(saveErrorBody).html("");
                    $(saveErrorBody).show();

                    var html = 
                        "<strong>Page rule errors</strong>"+
                        "<ul>";

                    for ( var i = 0; i < errors.length; i++ )
                        html += "<li>"+errors[i]+"</li>";

                    html += "</ul>";

                    $(saveErrorBody).html(html);
                }

                function mapOrigins(rules) {
                    var origins = {};

                    for ( var i = 0; i < rules.length; i++ ) {
                        var rule = rules[i],
                            urls = [ rule.matches, rule.host.baseUrl, rule.gitsense.baseUrl ];

                        for ( var j = 0; j < urls.length; j++ ) {
                            var a = document.createElement("a");
                            a.href = urls[j];
                            origins[a.origin+"/*"] = null;
                        }
                    }

                    return origins;
                }
            }

            function clickedAdd() {
                var rule = {
                    matches: "",
                    gitsense: {
                        baseUrl: "",
                        secret: ""
                    },
                    host: {
                        type: "",
                        api: "",
                        baseUrl: "",
                        secret: ""
                    }
                },

                id = "rule-"+rules.length;

                rules.push(rule);

                addRule(id, rule, true);

                setupEvents();
            }
        }

        function clean(str) {
            return str.replace(/\s/g, "");
        }
    }

    function restoreDefault() {
        if ( ! window.confirm("Are you sure") ) 
            return;

        var keys = [];

        for ( var key in sdes.config )
            keys.push(key);
   
        chrome.storage.local.remove(
            keys,
            function() {
                $(saveErrorBody).hide();
                update(sdes.config);
            }
        ) 
    }
}

document.addEventListener("DOMContentLoaded", load);
