fn main() {
    let x = 17;
    {
        let y = 25;
        println!("x is {}, y is {}", x, y);
    }
    println!("x is {}, y is {}", x, y);

}

