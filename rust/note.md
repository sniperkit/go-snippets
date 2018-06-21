1. Reference must not live longer than the resource they refer to. Rust will check the sopes of your reference to ensure that this is true.
2. Only thinks relating to references need lifetime.
3. Mutability: the ability to change something.
4. Mutability is a property of either a borrow(&mut) or a binding (let mut). This means that, for example, you cannot have a struct with some fields mutable and some immutable.
5. A trait is a language feature that tells the Rust compiler about functionality a type must provide.
