function overload(a: int): void {}
function overload(a: string): void {}
overload(42);
overload("foo");
