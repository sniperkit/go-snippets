this.a = 3;
function f() { return this; };
{a: 3, function () { return this.a; }};
