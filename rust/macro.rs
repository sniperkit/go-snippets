macro_rules! say_hello {
    () => {
        println!("Hello");
    }
}

macro_rules! create_function {
    ($fn_name: ident) => (
        fn $fn_name() {
            println!("You called {:?}()",
            stringify!($fn_name))
        }
        )
}

create_function!(foo);
create_function!(bar);

macro_rules! print_result{
    ($expression: expr) => (
        println!("{:?} = {:?}", stringify!($expression), $expression);
    )
}

//overload

macro_rules! test {
    ($left: expr; and $right: expr) => (
        println!("{:?} and {:?} is {:?}", stringify!($left), stringify!($right), $left && $right)
        );

    ($left: expr; or $right: expr) => (
        println!("{:?} or {:?} is {:?}", stringify!($left), stringify!($right), $left || $right)
        );
}

fn main() {
    say_hello!();

    foo();
    bar();

    print_result!(1u32 + 1);

    print_result!({
        let x = 1u32;

        x * x + 2 * x + 1
    });

    test!(1i32 + 1 == 2i32; and 2i32 * 2 == 3i32);
    test!(true; or false);
}
