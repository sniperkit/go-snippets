fn main() {
    let mut x = 10;

    loop {
        if x < 0 {
            break;
        }

        println!(" x is {}", x);
        x -= 1;
    }

    while x < 10 {
        x += 1;
        println!(" x is {}", x);
    }


    for i in 0..x {
        println!("i is {}", i);
    }

    for (idx, value) in (0..x).enumerate() {
        println!("i is {}, v is {}", idx, value);
    }

    'outer: for x in 0..10 {
        'inner: for y in 0..10 {
            if x % 2 == 0 { continue 'outer;}
            if y % 2 == 0 { continue 'inner;}
            println!("x is {}, y is {}", x, y);
        }
    }
}
