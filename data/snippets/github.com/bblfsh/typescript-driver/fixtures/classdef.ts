class Foo {
    member: string;
    constructor(message: string) {
        this.member = message;
    }
    greet() {
        return "Hello, " + this.member;
    }
}

let greeter = new Greeter("world");
