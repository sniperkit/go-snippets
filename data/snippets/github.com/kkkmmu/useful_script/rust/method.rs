#[derive(Debug)]
struct Circle {
    x: f64,
    y: f64,
    radius: f64,
}

impl Circle {
    fn area(&self) -> f64 {
        std::f64::consts::PI * (self.radius * self.radius)
    }

    fn grow(&self, increment: f64) -> Circle {
        Circle{ x: self.x, y: self.y, radius: self.radius + increment}
    }
}

impl Circle {
    fn new(x: f64, y: f64, radius: f64) -> Circle {
        Circle {
            x: x,
            y: y,
            radius: radius,
        }
    }
}

impl Circle {
    fn reference(&self) {
        println!("Taking self by reference");
    }

    fn mutable_reference(&mut self) {
        println!("Taking self by mutable reference!");
    }

    fn takes_ownership(self) {
        println!("Taking ownership of self!");
    }
}

struct CircleBuilder {
    x: f64,
    y: f64,
    radius: f64,
}

impl CircleBuilder {
    fn new() -> CircleBuilder {
        CircleBuilder { x: 0.0, y: 0.0, radius: 1.0, }
    }

    fn x(&mut self, coordinate: f64) -> &mut CircleBuilder {
        self.x = coordinate;
        self
    }

    fn y(&mut self, coordinate: f64) -> &mut CircleBuilder {
        self.y = coordinate;
        self
    }

    fn radius(&mut self, radius: f64) -> &mut CircleBuilder {
        self.radius = radius;
        self
    }

    fn finalize(&self) -> Circle {
        Circle{x: self.x, y: self.y, radius: self.radius }
    }
}

fn main() {
    let c = Circle{ x: 0.0, y: 0.0, radius: 3.0};
    println!("{}", c.area());

    println!("{:?}", c);
    let d = c.grow(2.0).area();
    println!("{}", d);
    //println!(c.reference());
    //println!(c.mutable_reference());
    //println!(c.takes_ownership());
    println!("{:?}", c);

    let c = Circle::new(0.0, 0.0, 11.0);
    println!("{:?}", c);

    let c = CircleBuilder::new().x(1.0).y(2.0).radius(3.0).finalize();
    println!("{:?}", c);
}
