sdes.utils.html = function() {
    var _this   = this,
        varUtil = new sdes.utils.variable();
   
    this.createBreak = function() {
        return createBasicElement("br");
    }

    this.createHeader3 = function(params) {
        return createBasicElement("h3", params);
    }

    this.createDiv = function(params) {
        return createBasicElement("div", params);
    }
    
    this.createPre = function(params) {
        return createBasicElement("pre", params);
    }
    
    this.createSpan = function(params) {
        return createBasicElement("span", params);
    }

    this.createNav = function(params) {
        return createBasicElement("nav", params);
    }

     this.createButton = function(params) {
        return createBasicElement("button", params);
    }

    this.createParagraph = function(params) {
        return createBasicElement("p", params)
    };

    this.createOrderedList = function(params) {
        return createTableElement("ol", params)
    }

    this.createUnorderedList = function(params) {
        return createTableElement("ul", params)
    }

    this.createList = function(params) {
        return createTableElement("li", params)
    }

    this.createTable = function(params) {
        return createTableElement("table", params)
    };
    
    this.createTableBody = function(params) {
        return createTableElement("tbody", params)
    }
    
    this.createTableRow = function(params) {
        return createTableElement("tr",params)
    };
    
    this.createTableHeader = function(params) {
        return createTableElement("th",params)
    };
    
    this.createTableCell = function(params) {
        return createTableElement("td", params)
    };
    
    this.createSup = function(params) {
        return createBasicElement("sup", params)
    };
    
    this.createIcon = function(params) {
        var icon = document.createElement("i"),
            style;
    
        if ( varUtil.isNoU(params) )
            return icon;
    
        if ( ! varUtil.isNoU(params.id) )
            icon.id = params.id;
    
        if ( ! varUtil.isNoU(params.cls) )
            icon.setAttribute("class", params.cls);
    
        if ( varUtil.isNoU(params.style) )
            return icon;
    
        for ( style in params.style )
            icon.style[style] = params.style[style];
    
        return icon;
    }
    
    this.createOptions = function(options) {
        var i,
            option,
            key,
            domOption,
            domOptions = [];
    
        for (i = 0; i < options.length; i++) {
            option = options[i];
            domOption = document.createElement("option");
    
            if ( option.text !== undefined )
                domOption.text = option.text;
    
            if ( option.selected )
                domOption.selected = true;
    
            domOptions.push(domOption);
        }
    
        return domOptions;
    }
    
    this.createSelect = function(options, params) {
        var select = document.createElement("select"),
            style,
            option;
    
        for ( i = 0; i < options.length; i++ ) {
            option = options[i];
            select.add(option);
        }
    
        if ( varUtil.isNoU(params) )
            return select;
    
        if ( varUtil.isNoU(params.style) )
            return select;
    
        for ( style in params.style ) {
            select.style[style] = params.style[style];
        }
    
        return select;
    }
    
    this.createTextInput = function(params) {
        var input = document.createElement("input"),
            style;
    
        input.type = "text";
    
        if ( varUtil.isNoU(params) )
            return input;
    
        if ( ! varUtil.isNoU(params.value) )
            input.value = params.value;
    
        if ( ! varUtil.isNoU(params.placeholder) )
            input.placeholder = params.placeholder;

        if ( ! varUtil.isNoU(params.cls) )
            input.setAttribute("class", params.cls);

        if ( ! varUtil.isNoU(params.disabled) )
            input.disabled = params.disabled;
    
        if ( varUtil.isNoU(params.style) )
            return input;
    
        for ( style in params.style )
            input.style[style] = params.style[style];
    
        return input;
    }
    
    this.createTextNode = function(text) {
        return document.createTextNode(text);
    }
    
    this.getSpan = function(html) {
        var span = document.createElement("span");
        $(span).html(html);
    
        return span;
    }
    
    this.createClearDiv = function() {
        var div = document.createElement("div");
        div.style.clear = "both";
        return div;
    }
    
    this.createLink = function(params) {
        if ( varUtil.isNoU(params) )
            params = {};
    
        var link = document.createElement("a"),
            option,
            style,
            node;
    
        if ( varUtil.isNoU(params) )
            return link;

        if ( ! varUtil.isNoU(params.text) ) {
            node = document.createTextNode(params.text);
            link.appendChild(node);
        } else if ( ! varUtil.isNoU(params.html) ) {
            $(link).html(params.html);
        }

        if ( ! varUtil.isNoU(params.append) ) {
            if ( params.append.length !== undefined ) {
                for ( var i = 0; i < params.append.length; i++ )
                    link.appendChild(params.append[i]);
            } else {
                link.appendChild(params.append);
            }
        }
    
        if ( ! varUtil.isNoU(params.id) )
            link.id = params.id;
    
        if ( ! varUtil.isNoU(params.cls) )
            link.setAttribute("class", params.cls);
    
        if ( ! varUtil.isNoU(params.target) )
            link.setAttribute("target", params.target);
    
        if ( ! varUtil.isNoU(params.href) )
            link.href = params.href;

        if ( ! varUtil.isNoU(params.ariaSelected) )
            link.setAttribute("aria-selected", params.ariaSelected);

        if ( ! varUtil.isNoU(params.target) )
            link.setAttribute("target", params.target);
    
        if ( varUtil.isNoU(params.style) )
            return link;

        for ( style in params.style )
            link.style[style] = params.style[style];
    
        return link;
    }
    
    this.createMessageTable = function(params) {
        var table =
                _this.createTable({
                    style: {
    
                    }
                });
    }
    
    function createTableElement(type, params) {
        var elem = document.createElement(type),
            style;
    
        if ( varUtil.isNoU(params) )
            return elem;
    
        if ( ! varUtil.isNoU(params.id) )
            elem.id = params.id;
    
        if ( ! varUtil.isNoU(params.html) )
            $(elem).html(params.html);
    
        if ( ! varUtil.isNoU(params.text) )
            $(elem).text(params.text);
    
        if ( ! varUtil.isNoU(params.cls) )
            elem.setAttribute("class", params.cls);

        if ( type === "td" && ! varUtil.isNoU(params.rowSpan) )
            elem.rowSpan = params.rowSpan;
    
        if ( type === "td" && ! varUtil.isNoU(params.colSpan) )
            elem.colSpan = params.colSpan;

        if ( ! varUtil.isNoU(params.append) ) {
            if ( params.append.length !== undefined ) {
                for ( var i = 0; i < params.append.length; i++ )
                    elem.appendChild(params.append[i]);
            } else {
                elem.appendChild(params.append);
            }
        }
    
        if ( varUtil.isNoU(params.style) )
            return elem;
    
        for ( style in params.style )
            elem.style[style] = params.style[style];
    
        return elem;
    }
    
    function createBasicElement(type, params) {
        var elem = document.createElement(type),
            style;
    
        if ( varUtil.isNoU(params) )
            return elem;
    
        if ( ! varUtil.isNoU(params.id) )
            elem.id = params.id;

        if ( ! varUtil.isNoU(params.ariaLabel) )
            elem.setAttribute("aria-label", params.ariaLabel);
    
        if ( ! varUtil.isNoU(params.role) )
            elem.setAttribute("role", params.role);

        if ( ! varUtil.isNoU(params.title) )
            elem.setAttribute("title", params.title);
    
        if ( ! varUtil.isNoU(params.html) ) {
            if ( typeof(params.html) === "object" )
                elem.appendChild(params.html);
            else
                $(elem).html(params.html);
        }
    
        if ( ! varUtil.isNoU(params.text) )
            $(elem).text(params.text);
    
        if ( ! varUtil.isNoU(params.cls) )
            elem.setAttribute("class", params.cls);

        if ( ! varUtil.isNoU(params.pointerEvents) )
            elem.setAttribute("pointer-events", params.pointerEvents);

        if ( ! varUtil.isNoU(params.append) ) {
            if ( params.append.length !== undefined ) {
                for ( var i = 0; i < params.append.length; i++ )
                    elem.appendChild(params.append[i]);
            } else {
                elem.appendChild(params.append);
            }
        }
    
        if ( varUtil.isNoU(params.style) )
            return elem;
    
        for ( style in params.style )
            elem.style[style] = params.style[style];
    
        return elem;
    }
    
    this.createCheckbox = function(params) {
        var checkbox = document.createElement("input"),
            style;
    
        checkbox.type = "checkbox";
    
        if ( ! varUtil.isNoU(params.id) )
            checkbox.id = params.id;
    
        if ( ! varUtil.isNoU(params.checked) )
            checkbox.checked = params.checked;
    
        if ( ! varUtil.isNoU(params.disabled) )
            checkbox.disabled = params.disabled;
    
        if ( ! varUtil.isNoU(params.title) )
            checkbox.title = params.title;
    
        if ( varUtil.isNoU(params.style) )
            return checkbox;
    
        for ( style in params.style )
            checkbox.style[style] = params.style[style];
    
        return checkbox;
    }
    
    this.createImage = function(params) {
        var img = document.createElement("img"),
            style;
    
        if ( ! varUtil.isNoU(params.id) )
            img.id = params.id;
    
        if ( ! varUtil.isNoU(params.src) )
            img.src = params.src;
    
        if ( ! varUtil.isNoU(params.cls) )
            img.setAttribute("class", params.cls);

        if ( ! varUtil.isNoU(params.title) )
            img.title = params.title;
    
        if ( varUtil.isNoU(params.style) )
            return img;
    
        for ( style in params.style )
            img.style[style] = params.style[style];
    
        return img;
    }
    
    this.createSvg = function(params) {
        return createBasicElement("svg", params);
    }
    
    this.getImage = function(id, src, cls, styles) {
        var img = document.createElement("img"),
            style;
    
        if ( ! varUtil.isNoU(params.id) )
            img.id = id;
    
        if ( ! varUtil.isNoU(params.src) )
            img.src = src;
    
        if ( ! varUtil.isNoU(params.cls) )
            img.setAttribute("class",cls);
    
        if ( varUtil.isNoU(params.style) )
            return img;
    
        for ( style in styles )
            img.style[style] = styles[style];
    
        return img;
    }
}
