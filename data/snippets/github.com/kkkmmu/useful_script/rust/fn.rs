fn print_number(x :i32) {
    println!("x is {}", x);
}

fn print_sum(x :i32, y :i32) {
    println!("x + y is {}", x + y);
}

fn add_one(x :i32) -> i32 {
    x + 1
}

fn radd_one(x :i32) -> i32 {
    return x + 1;
}

fn diverges() -> ! {
    panic!("This function never returns!");
}

fn main() {

    let f: fn(i32) -> i32 = add_one;
    let g = add_one;
    println!("f(1): is {}", f(1));
    println!("g(2): is {}", g(2));
}
