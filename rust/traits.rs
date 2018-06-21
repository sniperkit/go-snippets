fn main() {
    basic();
    bound();
    struct_bound();
    trait_for_float();
    where_clause();
    where_clause_real();
    convert_to();
    default_methods();
    inheritance();
}

fn inheritance() {
    trait Foo {
        fn foo(&self);
    }

    trait FooBar : Foo {
        fn foobar(&self);
    }

    struct Baz;

    impl Foo for Baz {
        fn foo(&self) { println!("foo"); }
    }

    impl FooBar for Baz {
        fn foobar(&self) { println!("foobar");}
    }
}

fn default_methods() {
    trait Foo {
        fn is_valid(&self) -> bool;
        fn is_invalid(&self) -> bool { !self.is_valid() }
    }
    struct UseDefault;

    impl Foo for UseDefault {
        fn is_valid(&self) -> bool {
            println!("Called UseDefault.is_valid");
            true
        }
    }

    struct OverrideDefault;

    impl Foo for OverrideDefault {
        fn is_valid(&self) -> bool {
            println!("Called overrideDefault.is_valid.");
            true
        }

        fn is_invalid(&self) -> bool {
            println!("Called OverrideDefault.is_invalid.");
            true
        }
    }
}

fn convert_to() {
    trait ConvertTo<Output> {
        fn convert(&self) -> Output;
    }

    impl ConvertTo<i64> for i32 {
        fn convert(&self) -> i64 { *self as i64 }
    }

    fn convert_t_to_i64<T: ConvertTo<i64>>(x: T) -> i64 {
        x.convert()
    }

    fn convert_i32_to_t<T> (x: i32) -> T where i32: ConvertTo<T> {
        x.convert()
    }
}

fn where_clause_real() {
    use std::fmt::Debug;

    fn foo<T: Clone, K: Clone + Debug> (x: T, y: K) {
        x.clone();
        y.clone();
        println!("{:?}", y);
    }

    fn bar<T, K> (x: T, y: K) where T :Clone, K: Clone + Debug {
        x.clone();
        y.clone();
        println!("{:?}", y);
    }

    foo(1, 2);
    bar(3, 4);
}

fn where_clause() {
    use std::fmt::Debug;

    fn foo<T: Clone, K: Clone + Debug>(x: T, y: K) {
        x.clone();
        y.clone();
        println!("{:?}", y);
    }

    foo(1, 23);
}

fn trait_for_float() {
    trait ApproxEqual {
        fn approx_equal(&self, other: &Self) -> bool;
    }

    impl ApproxEqual for f32 {
        fn approx_equal(&self, other: &Self) -> bool {
            (self - other).abs() <= ::std::f32::EPSILON
        }
    }

    println!("{}", 1.0.approx_equal(&1.000000001));
}

fn struct_bound() {
    struct Rectangle<T> {
        x: T,
        y: T,
        width: T,
        height: T,
    }

    impl <T: PartialEq> Rectangle<T> {
        fn is_square(&self) -> bool {
            self.width == self.height
        }
    }

    let mut r = Rectangle {
        x: 0,
        y: 0,
        width: 47,
        height: 47,
    };

    assert!(r.is_square());
    r.height = 42;
    assert!(!r.is_square());
}

fn bound() {
    trait HasArea {
        fn area(&self) -> f64;
    }

    struct Circle {
        x: f64,
        y: f64,
        radius: f64,
    }

    impl HasArea for Circle {
        fn area(&self) -> f64 {
            std::f64::consts::PI * (self.radius * self.radius)
        }
    }

    struct Square {
        x: f64,
        y: f64,
        side: f64,
    }

    impl HasArea for Square {
        fn area(&self) -> f64 {
            self.side * self.side
        }
    }

    fn print_area<T: HasArea>(shape: T) {
        println!("This shape has an area of {}", shape.area());
    }

    let c = Circle {
        x: 0.0f64,
        y: 0.0f64,
        radius: 3.0f64,
    };

    let s = Square {
        x: 0.0f64,
        y: 0.0f64,
        side: 1.0f64,
    };

    print_area(c);
    print_area(s);
}

fn basic() {
    struct Circle{
        x: f64,
        y: f64,
        radius: f64,
    }

    trait HasArea {
        fn area(&self) -> f64;
        fn is_larger(&self, &Self) -> bool;
    }

    impl Circle {
        fn area(&self) -> f64 {
            std::f64::consts::PI * (self.radius * self.radius)
        }
    }

    impl HasArea for Circle {
        fn area(&self) -> f64 {
            std::f64::consts::PI * (self.radius * self.radius)
        }

        fn is_larger(&self, other: &Self) -> bool {
            self.area() > other.area()
        }
    }
}
