//1. trait object link &Foo and Box<Foo> are normal values that store a value of any type that
//implenments the given trait, where the precise type can only be kown at runtime.
//2. A trait object can be obtained from a pointer to a concrete type that implements the trait by
//   casting it : 
//      *&x as &Foo 
//  or coercing:
//      &x as an arguemnt to a function that takes &Foo
//3. This operation can be seen as "erasing" the compiler's knowledge about the specific type of
//   the pointer, and hence trait objects are sometimes referred to as "type erasure".
//

trait Foo {
    fn method(&self) -> String;
}

//The method of the trait can be called on a trait object via a special record of function pointers
//traditionally called a 'vtable' (created and managed by the compiler).
impl Foo for u8 {
    fn method(&self) -> String {
        format!("u8: {}", *self)
    }
}

impl Foo for String {
    fn method(&self) -> String {
        format!("String: {}", *self)
    }
}

//A function that takes a trait object is not speciallized to each of the types that implements the
//trait. Only one copy is generated, often (but not always) resulting in less code bloat.
//Howerever, this comes at the cost of requiring slower virtual function calls, and effectively
//inhibiting any chance of inlining and related optimizations from occuring.
fn do_something(x: &Foo) {
    x.method();
}


//Not every trait can be used to make a trait object. Only trait that are objec-safe can be made
//into trait objects. A trait is ojbect-safe if both of these are true:
//  1. the trait does not require that Self: Sized.
//  2. all of its mehtods are object-safe.
//
//So waht makes a method object-safe? Each method must require that: 
//  Self: Sized.
//or all of the fowllowing:
//  1. Must not have any type parameters.
//  2. Must not use Self.
fn main() {
    let x = 5u8;
    let y = "Hello".to_string();

    do_something(&x as &Foo);
    do_something(&y as &Foo);

    do_something(&x);
    do_something(&y);

}
