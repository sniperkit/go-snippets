fn main() {
    plus_one();
    plus_two();
    callable();
}

fn plus_one() {
    let po = | x: i32| x + 1;
    assert_eq!(2, po(1));
}

fn plus_two() {

    let pt = |x| {
        let mut result: i32 = x;
        result += 1;
        result += 1;
        result
    };

    assert_eq!(4, pt(2));
}

fn callable() {
    fn call_with_one<F> (some_closure: F) -> i32 where F: Fn(i32) -> i32 {
        some_closure(1)
    }

    let answer = call_with_one(|x| x + 2);
    assert_eq!(3, answer);

    fn call_with_one1(some_closure: &Fn(i32) -> i32) -> i32 {
        some_closure(1)
    }

    let answer = call_with_one1(&|x| x + 2);

    assert_eq!(3, answer);

    fn call_with_ref<F> (some_closure: F) -> i32 where F: Fn(&i32) -> i32 {
        let value = 0;
        some_closure(&value)
    }
}
