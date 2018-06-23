abstract class Foo {
    abstract makeFoo(): number;
    implemented(): void {}
}

class Bar extends Foo {
    makeFoo(): number { return 42; }
}
