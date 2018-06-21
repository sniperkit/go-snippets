# Selectors
CSS has its own terminology to describe the CSS language. Previously in this tutorial, you created a line in your stylesheet like this:

strong {
    color: red;
}
# In CSS terminology, this entire line is a rule. This rule starts with strong, which is a selector (or a selector list). It selects which elements in the DOM the rule applies to.

# More details
## The part inside the curly braces is the declaration.

# The keyword color is a property, and red is a value.

# The semicolon after the property-value pair separates it from other property-value pairs in the same declaration.

This tutorial refers to a selector like strong as a tag selector. The CSS Specification refers to it as a type selector.

This page of the tutorial explains more about the selectors that you can use in CSS rules.

# In addition to tag names, you can use attribute values in selectors. This allows your rules to be more specific.

Two attributes have special status for CSS. They are class and id.

# Class selectors

Use the class attribute in an element to assign the element to a named class. It is up to you what name you choose for the class. Multiple elements in a document can have the same class value.

# ID selectors

Use the id attribute in an element to assign an ID to the element. It is up to you what name you choose for the ID. The ID name must be unique in the document.

## Example
This HTML tag has both a class attribute and an id attribute:
    <p class="key" id="principal">
     The id value, principal, must be unique in the document, but other tags in the document can have the same class name, key.

     In a CSS stylesheet, this rule makes all the elements with class key green. (They might not all be <p> elements.)

.key {
color: green;
}
This rule makes the one element with the id principal bold:

#principal {
    font-weight: bolder;
}
# Attribute Selectors

You are not restricted to the two special attributes, class and id. You can specify other attributes by using square brackets. Inside the brackets you put the attribute name, optionally followed by a matching operator and a value. Additionally, matching can be made case-insensitive by appending an " i" after the value, but not many browsers support this feature yet. Examples:

## [disabled]
Selects all elements with a "disabled" attribute.
## [type='button']
Selects elements with a "button" type.
## [class~=key]
Selects elements with the class "key" (but not e.g. "keyed", "monkey", "buckeye"). Functionally equivalent to .key.
## [lang|=es]
Selects elements specified as Spanish. This includes "es" and "es-MX" but not "eu-ES" (which is Basque).
## [title*="example" i]
Selects elements whose title contains "example", ignoring case. In browsers that don't support the "i" flag, this selector probably won't match any element.
## a[href^="https://"]
Specifies what the attribute's value should start with; in this case, it selects secure links.
## img[src$=".png"]
Selects elements who's value ends with the provided string. Indirectly selects PNG images; any images that are PNGs but whose URL doesn't end in ".png" (e.g. `src="some-image.png?_=cachebusterhash"`) won't be selected.
Pseudo-classes selectors

# A CSS pseudo-class is a keyword added to selectors that specifies a special state of the element to be selected. For example :hover will apply a style when the user hovers over the element specified by the selector.

# Pseudo-classes, together with pseudo-elements, let you apply a style to an element not only in relation to the content of the document tree, but also in relation to external factors like the history of the navigator (:visited, for example), the status of its content (like :checked on some form elements), or the position of the mouse (like :hover which lets you know if the mouse is over an element or not). To see a complete list of selectors, visit CSS3 Selectors working spec.

Syntax
selector:pseudo-class {
property: value;
}
List of pseudo-classes

:link
:visited
:active
:hover
:focus
:first-child
:last-child
:nth-child
:nth-last-child
:nth-of-type
:first-of-type
:last-of-type
:empty
:target
:checked
:enabled
:disabled
Selector lists

A rule can be shared by multiple selectors by using a comma (,) to separate selectors.

Example
In the example, both of elements that have the "content-1" class and elements that have the "content-2" class will display bold text.
.content-1, .content-2 {
    font-weight: bold;
}
Information: Specificity
Multiple rules may have selectors that each match the same element. If a property is given in only one rule, there is no conflict and the property is set on the element. If more than one rule applies to an element and sets the same property, then CSS gives priority to the rule that has the more specific selector. An ID selector is more specific than a class, pseudo-class or attribute selector, which in turn are more specific than a tag or pseudo-element selector.

More details
You can also combine selectors, making a more specific selector. For example, the selector .key selects all elements that have the class name key. The selector p.key selects only <p> elements that have the class name key.

If the stylesheet has conflicting rules and they are equally specific, then CSS gives priority to the rule that is later in the stylesheet.

When you have a problem with conflicting rules, try to resolve it by making one of the rules more specific, so that it has priority. If you cannot do that, try moving one of the rules nearer the end of the stylesheet so that it has priority.

Information: Selectors based on relationships
CSS has some ways to select elements based on the relationships between elements. You can use these to make selectors that are more specific.

Common selectors based on relationships
    Selector    Selects
A E Any E element that is a descendant of an A element (that is: a child, or a child of a child, etc.)
    A > E   Any E element that is a child (i.e. direct descendant) of an A element
    E:first-child   Any E element that is the first child of its parent
B + E   Any E element that is the next sibling of a B element (that is: the next child of the same parent)
    You can combine these to express complex relationships.

    You can also use the symbol * (asterisk) to mean "any element".

    Example
    An HTML table has an id attribute, but its rows and cells do not have individual identifiers:

    <table id="data-table-1">
    ...
    <tr>
    <td>Prefix</td>
    <td>0001</td>
    <td>default</td>
    </tr>
    ...
    These rules make the first cell in each row underline, and the sibling of first cell in each row strikethroughted (in example 2.nd cell) . They only affect one specific table in the document:

#data-table-1 td:first-child {text-decoration: underline;}
#data-table-1 td:first-child + td {text-decoration: line-through;}
    The result looks like:

    Prefix  0001    default
    More details
    In the usual way, if you make a selector more specific, then you increase its priority.

    If you use these techniques, you avoid the need to specify class or id attributes on so many tags in your document. Instead, CSS does the work.

    In large designs where speed is important, you can make your stylesheets more efficient by avoiding complex rules that depend on relationships between elements.

    For more examples about tables, see Tables in the CSS Reference page.

    Action: Using class and ID selectors
    Edit your HTML file, and duplicate the paragraph by copying and pasting it.
    Then add id and class attributes to the first copy, and an id attribute to the second copy as shown below. Alternatively, copy and paste the entire file again:
    <!doctype html>
    <html>
    <head>
    <meta charset="UTF-8">
    <title>Sample document</title>
    <link rel="stylesheet" href="style1.css">
    </head>
    <body>
    <p id="first">
    <strong class="carrot">C</strong>ascading
    <strong class="spinach">S</strong>tyle
    <strong class="spinach">S</strong>heets
    </p>
    <p id="second">
    <strong>C</strong>ascading
    <strong>S</strong>tyle
    <strong>S</strong>heets
    </p>
    </body>
    </html>
    Now edit your CSS file. Replace the entire contents with:
    strong { color: red; }
    .carrot { color: orange; }
    .spinach { color: green; }
#first { font-style: italic; }
    Save the files and refresh your browser to see the result:
    Cascading Style Sheets
    Cascading Style Sheets
    You can try rearranging the lines in your CSS file to show that the order has no effect.

    The class selectors .carrot and .spinach have priority over the tag selector strong.

    The ID selector #first has priority over the class and tag selectors.

    Challenges
    Without changing your HTML file, add a single rule to your CSS file that keeps all the initial letters that same color as they are now, but makes all the other text in the second paragraph blue:
    Cascading Style Sheets
    Cascading Style Sheets
    Now change the rule you have just added (without changing anything else), to make the first paragraph blue too:
    Cascading Style Sheets
    Cascading Style Sheets
    See a solution for the challenge.
    Action: Using pseudo-classes selectors
    Create an HTML file with following content:
    <!doctype html>
    <html>
    <head>
    <meta charset="UTF-8">
    <title>Sample document</title>
    <link rel="stylesheet" href="style1.css">
    </head>
    <body>
    <p>Go to our <a class="homepage" href="http://www.example.com/" title="Home page">Home page</a>.</p>
    </body>
    </html>
    Now edit your CSS file. Replace the entire contents with:
    a.homepage:link, a.homepage:visited {
padding: 1px 10px 1px 10px;
color: #fff;
background: #555;
            border-radius: 3px;
border: 1px outset rgba(50,50,50,.5);
        font-family: georgia, serif;
        font-size: 14px;
        font-style: italic;
        text-decoration: none;
    }

a.homepage:hover, a.homepage:focus, a.homepage:active {
    background-color: #666;
}
Save the files and refresh your browser to see the result (put the mouse over the following link to see the effect):
    Go to our Home page
    Action: Using selectors based on relationships and pseudo-classes
    With selectors based on relationships and pseudo-classes you can create complex cascade algorithms. This is a common technique used, for example, in order to create pure-CSS dropdown menus (that is only CSS, without using JavaScript). The essence of this technique is the creation of a rule like the following:

    div.menu-bar ul ul {
display: none;
    }

div.menu-bar li:hover > ul {
display: block;
}
to be applied to an HTML structure like the following:

<div class="menu-bar">
<ul>
<li>
<a href="example.html">Menu</a>
<ul>
<li>
<a href="example.html">Link</a>
</li>
<li>
<a class="menu-nav" href="example.html">Submenu</a>
<ul>
<li>
<a class="menu-nav" href="example.html">Submenu</a>
<ul>
<li><a href="example.html">Link</a></li>
<li><a href="example.html">Link</a></li>
<li><a href="example.html">Link</a></li>
<li><a href="example.html">Link</a></li>
</ul>
</li>
<li><a href="example.html">Link</a></li>
</ul>
</li>
</ul>
</li>
</ul>
</div>
See our complete CSS-based dropdown menu example for a possible cue.

What next?
Your sample stylesheet is starting to look dense and complicated. The next section describes ways to make CSS easier to read.
