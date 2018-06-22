function foo<T, V>(arg: T): V {}
let x: string = foo<number, string>(42, "bar");
let auto = foo(21, "baz");

function arrT<T>(arg: T[]): void {}
function arrTAlt<T>(arg: Array<T>): void {}
let y = arrT([1, 2, 3])
let z = arrTAlt([4, 5, 6])
