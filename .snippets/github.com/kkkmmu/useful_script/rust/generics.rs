fn main() {
    enum Option<T> {
        Some(T),
        None,
    }


    let x : Option<i32> = Option::Some(5);
    let y : Option<f64> = Option::Some(1.0);

    enum Result<T, E>  {
        OK(T),
        Err(E),
    }

    let int_origin = Point{ x: 0, y: 0};
    let float_origin = Point{ x: 0.1, y: 0.1};

    let mut v = Vec::new();

    v.push(true);
    println!("{:?}", v);

    let v: Vec<bool> = Vec::new();
    println!("{:?}", v);

    let v = Vec::<bool>::new();
    println!("{:?}", v);
    /*
    take_anything(10);
    take_anything(10.00);
    take_anything("Hello world");
    */
}

fn take_anything<T> (x: T) {
    //println!("{:?}", x)
}

fn takes_one_of_the_same_thing<T> (x: T, y: T) {

}

fn takes_two_things<T, U>(x: T, y: U) {

}

struct Point<T> {
    x: T,
    y: T,
}
    
impl <T> Point<T> {
    fn swap(&mut self) {
        std::mem::swap(&mut self.x, &mut self.y);
    }
}
