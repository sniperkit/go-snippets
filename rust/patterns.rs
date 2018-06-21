
fn main() {
    let x = 1;
    match x {
        1 => println!("one"),
        2 => println!("two"),
        3 => println!("Three"),
        _ => println!("anything!"),
    }

    let x = 1;
    match x {
        y => println!("x is {}, y is {}", x, y),
    }

    let x = 10;
    let c = 'c';

    match c {
        x => println!("x is {}, c is {}", x, c),
    }
    println!("x is {}", x);

    let x = 1000;

    match x {
        1 | 2 => println!("one or two"),
        3 | 1000 => println!("three or 1000"),
        _ => println!("anything"),
    }

    destructuring();
    destructuring1();
    destructuring2();
    ignoremove();
    optionaltuple();
    reference();
    matchrange();
    bindings();
    structbindings();
    guards();
}

fn destructuring() {
    struct Point{
        x: i32,
        y: i32,
    }

    let origin = Point{x: 0, y: 0};
    match origin {
        Point {x , y} => println!("{}, {}", x, y),
    }
}

fn destructuring1() {
    struct Point{ 
        x: i32,
        y: i32,
    }

    let origin = Point{x: 0, y: 0};
    match origin {
        Point{ x: x1, y: y1} => println!("{}, {}", x1, y1),
    }
}

fn destructuring2() {
    struct Point{
        x: i32,
        y: i32,
    }

    let point = Point { x: 3, y: 4};
    match point {
        Point{x, ..} => println!("x is {}", x),
    }
    match point {
        Point{y, ..} => println!("y is {}", y),
    }

}

fn ignoremove() {
    let tuple: (u32, String) = (5, String::from("five"));
    let (x, _s) = tuple;

    //println!("Tuple is : {:?}", tuple);

    let tuple = (6, String::from("six"));

    let (x, _) = tuple;

    println!("Tuple is : {:?}", tuple);
}

fn optionaltuple() {
    enum OptionalTuple {
        Value(i32, i32, i32),
        Missing,
    }

    let x = OptionalTuple::Value(5, 3, 4);

    match x {
        OptionalTuple::Value(..) => println!("Got a tuple"),
        OptionalTuple::Missing => println!("no such suck"),
    }
}

fn reference() {
    let x = 5;

    match x {
        ref r => println!("Got a reference to {}", r),
    }

    let mut x = 6;
    match x {
        ref mut mr => println!("Got a mutable reference to {}", mr),
    }
}

fn matchrange() {
    let x = 'b';

    match x {
        'a' ... 'j' => println!("early letter"),
        'k' ... 'z' => println!("late letter"),
        _ => println!("Someting else"),
    }
}

fn bindings() {
    let x = 9;

    match x { 
        e @ 1...5 => println!("got a range element {}", e),
        _ => println!("anything"),
    }

    match x {
        e @ 1...5 | e @ 8...10 => println!("got a range element: {}", e),
        _ => println!("anything"),
    }
}

fn structbindings() {
    #[derive(Debug)]
    struct Person{
        name: Option<String>,
    }

    let name = "Steve".to_string();
    let x : Option<Person> = Some(Person {name: Some(name)});

    match x {
        Some(Person{name: ref a @ Some(_), ..}) => println!("{:?}", a),
        _ => {},
    }
}

fn guards() {
    enum OptionalInt {
        Value(i32),
        Missing,
    }

    let x = OptionalInt::Value(5);
    match x {
        OptionalInt::Value(i) if i > 5 => println!("Got a int larger than 5"),
        OptionalInt::Value(..) => println!("Got a int"),
        OptionalInt::Missing => println!("No such suck"),
    }
}
