let message = "liwei";
alert("Hello world" + message);
alert(Infinity);
alert("abc"/2)
function showMessage(from, text = "no text given") {
      alert( from + ": " + text );
}

showMessage("Ann"); // Ann: no text given

let user = {
    name: "John"
};

let descriptor = Object.getOwnPropertyDescriptor(user, 'name');

alert(JSON.stringify(descriptor, null, 2));
