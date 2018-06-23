interface Foo {
    bar: string;
    baz?: number;
}

function useFoo(arg: Foo) {}

class extendsFoo implements Foo {
}
