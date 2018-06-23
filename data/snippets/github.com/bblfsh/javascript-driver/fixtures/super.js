class A { a() {} }
class B extends A { constructor() { super(); } }
class C extends A { b() { super.a(); } }
