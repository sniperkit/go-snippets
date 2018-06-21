fn main() {
    let x = (1, "hello");
    let y: (i32, &str) = (2, "world");

    println!("x is {:?}" , x);
    println!("y is {:?}", y);

    let (a, b) = x;
    println!("a is {}, b is {}", a, b);
    println!("y.0 is {}, y.1 is {}", y.0, y.1);
}
