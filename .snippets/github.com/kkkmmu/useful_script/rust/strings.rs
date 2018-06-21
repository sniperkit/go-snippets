fn main() {
    let mut s = "Hello".to_string();
    println!("{}", s);

    s.push_str(", liwei!");
    println!("{}", s);

    takes_slices(&s);

    for b in s.as_bytes() {
        println!("{}", b);
    }

    for c in s.chars() {
        println!("{}", c);
    }

    println!("{:?}", s.chars().nth(1));
    println!("{:?}", s.chars().nth(2));
    println!("{:?}", s.chars().nth(3));
}

fn takes_slices(slice :&str) {
    println!("Got : {}", slice);
}

