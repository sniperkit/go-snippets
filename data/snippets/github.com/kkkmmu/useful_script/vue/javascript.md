1.  JavaScript can execute not only in the browser, but also on the server, or actually on any device where there exists a special program called the JavaScript engine.
2. In-browser JavaScript is able to:
    Add new HTML to the page, change the existing content, modify styles.
    React to user actions, run on mouse clicks, pointer movements, key presses.
    Send requests over the network to remote servers, download and upload files (so-called AJAX and COMET technologies).
    Get and set cookies, ask questions to the visitor, show messages.
    Remember the data on the client-side (“local storage”).
3. A single <script> tag can’t have both the src attribute and the code inside.
4. It’s recommended to put semicolons between statements even if they are separated by newlines. This rule is widely adopted by the community. Let’s note once again – it is possible to leave out semicolons most of the time. But it’s safer – especially for a beginner – to use them.
5. One-line comments start with two forward slash characters //.
6. Multiline comments start with a forward slash and an asterisk /* and end with an asterisk and a forward slash */.
7. the dollar sign '$' and the underscore '_' can also be used in names. They are regular symbols, just like letters, without any special meaning.
8. Normally, we need to define a variable before using it. But in the old times, it was technically possible to create a variable by a mere assignment of the value, without let. This still works now if we don’t put use strict. The behavior is kept for compatibility with old scripts.
9. Variables declared using const are called “constants”. They cannot be changed. An attempt to do it would cause an error.
10. var – is an old-school variable declaration. Normally we don’t use it at all, but we’ll cover subtle differences from let in the chapter The old "var", just in case you need them.
11. A variable in JavaScript can contain any data. A variable can at one moment be a string and later receive a numeric value.  Programming languages that allow such things are called “dynamically typed”, meaning that there are data types, but variables are not bound to any of them.
12. Double and single quotes are “simple” quotes. There’s no difference between them in JavaScript.
13. Backticks are “extended functionality” quotes. They allow us to embed variables and expressions into a string by wrapping them in ${…}.
14. If a variable is declared, but not assigned, then its value is exactly undefined.
15. The typeof operator returns the type of the argument. It’s useful when we want to process values of different types differently, or just want to make a quick check.
16. Math is a built-in object that provides mathematical operations. 
17. A strict equality operator === checks the equality without type conversion.
18. If a function does not return a value, it is the same as if it returns undefined. An empty return is also the same as return undefined.
19. In JavaScript, a function is a value, so we can deal with it as a value. 
20. An object can be created with figure brackets {…} with an optional list of properties. A property is a “key: value” pair, where key is a string (also called a “property name”), and value can be anything.
21. An empty object (“empty cabinet”) can be created using one of two syntaxes:
    let user = new Object(); // "object constructor" syntax
    let user = {};  // "object literal" syntax
22. To remove a property, we can use delete operator.
23. We can use square brackets in an object literal. That’s called computed properties.
24. The meaning of a computed property is simple: [fruit] means that the property name should be taken from fruit.
25. Most of the time, when property names are known and simple, the dot is used. And if we need something more complex, then we switch to square brackets.
26. A variable cannot have a name equal to one of language-reserved words like “for”, “let”, “return” etc. But for an object property, there’s no such restriction. 
27. In real code we often use existing variables as values for property names.
28. The use-case of making a property from a variable is so common, that there’s a special property value shorthand to make it shorter. Instead of name:name we can just write name. We can use both normal properties and shorthands in the same object.
29. We can use a special operator "in" to check for the existence of a property. Please note that on the left side of in there must be a property name. That’s usually a quoted string. If we omit quotes, that would mean a variable containing the actual name to be tested. F
30. To walk over all keys of an object, there exists a special form of the loop: for..in. 
31. One of the fundamental differences of objects vs primitives is that they are stored and copied “by reference”. Primitive values: strings, numbers, booleans – are assigned/copied “as a whole value”.
32. Two objects are equal only if they are the same object.
33. The JSON (JavaScript Object Notation) is a general format to represent values and objects. It is described as in RFC 4627 standard. Initially it was made for JavaScript, but many other languages have libraries to handle it as well. So it’s easy to use JSON for data exchange when the client uses JavaScript and the server is written on Ruby/PHP/Java/Whatever.

    JavaScript provides methods:

        JSON.stringify to convert objects into JSON.
        JSON.parse to convert JSON back into an object.
