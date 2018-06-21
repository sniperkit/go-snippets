fn main() {
    println!("Hello world!");
    hello("liwei");
    hello("God!");
    println!("x + y = {}", sum(1,2));

    let x = 11;
    let z = if x == 10 {
        100
    } else {
        200
    };

    println!("z is {}", z);
}

fn hello(name :&'static str) {
    println!("Hello {}!", name);
}

fn sum(x :u32, y :u32) -> u32 {
    x + y
}

