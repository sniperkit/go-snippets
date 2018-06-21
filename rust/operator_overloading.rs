use std::ops::Add;

#[derive(Debug)]
struct Point {
    x: i32,
    y: i32,
}

impl Add for Point {
    type Output = Point;

    fn add(self, other: Point) -> Point {
        Point {x: self.x + other.x, y: self.y + other.y }
    }
}

fn main() {
    let p1 = Point {x: 1, y: 0};
    let p2 = Point {x: 3, y: 5};

    let p3 = p2 + p1;

    println!("{:?}", p3);

    traits();
}

fn traits() {
    use std::ops::Mul;

    trait HashArea<T> {
        fn area(&self) -> T;
    }

    struct Square<T> {
        x: T,
        y: T,
        side: T,
    }

    impl<T> HashArea<T> for Square<T> where T: Mul<Output=T> + Copy {
        fn area(&self) -> T {
            self.side * self.side
        }
    }

    let s = Square {
        x: 0.0f64,
        y: 0.0f64,
        side: 12.0f64,
    };

    println!("Area of s: {}", s.area())

}
