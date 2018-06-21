fn basic() {
    let num = 5;
    //    let plus_num = |x| x + num;
    let plus_num = |x: i32| x + num;

    assert_eq!(10, plus_num(5));
}

fn bcopy() {
    let mut num = 5;
    {
        let mut add_num = |x: i32| num += x;
        add_num(5);
    }

    assert_eq!(10, num);
}

//Another way to think about move closures:
//  They give a closure its own stack frame.
//Without move, a closure may be tied to the stack frame that created it, while a move closure is
//self-contained. 
//*****This means that you cannot generally return a non-move closure from a function.
fn bmove() {
    let mut num = 5;
    {
        let mut add_num = move |x: i32| num += x;
        add_num(5);
    }

    assert_eq!(5, num);
}

//Using () to call a function like 'foo()', is an overloadable operator. From this, everything else
//clicks into place. In Rust, we use the trait system to overload operators. Calling functions is
//no different. We have thre seperate traits to overload with:
//  Fn (take &self)
//  FnMut (take &mut self)
//  FnOnce (take self)
//The ||{} syntax for closures is sugar for these three traits. Rust will generate a struct for the
//environment, impl the approprate trait, and then use it.
//


//Fn is a trait, so it can be used as a generic type bound
fn call_with_one<F>(some_closure: F) -> i32 where F: Fn(i32) -> i32 {
    some_closure(1)
}

fn call_with_one_dynamic(some_closure: &Fn(i32) -> i32) -> i32 {
    some_closure(1)
}

fn call_with_ref<F> (some_closure: F) -> i32 where F: for<'a> Fn(&'a i32) -> i32 {
    let value = 1;
    some_closure(&value)
}

fn main() {
    basic();
    bcopy();
    bmove();

    let answer = call_with_one(|x|x+2);
    assert_eq!(3, answer);
    let answer1 = call_with_one_dynamic(&|x|x+2);
    assert_eq!(3, answer1);

    let answer2 = call_with_ref(|x|x+2);
    assert_eq!(3, answer2);

    return_closure();
}

//In order to return somthing from a function, Rust need to know what size the return type is. An
//easy way to give something a size to take a reference to it, as reference have a known size.

fn return_closure() {
    fn factory() -> Box<Fn(i32) -> i32> {
        let num = 5;
        Box::new(move |x| x+num)
    }

    let f = factory();
    let answer = f(1);
    assert_eq!(6, answer);
}
