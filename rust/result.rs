/*
enum Result<T, E> {
    Ok(T),
    Err(E),
}

type Option<T> = Result<T, ()>;

impl <T, E: ::std::fmt::Debug> Result<T, E> {
    fn unwrap(self) -> T {
        match self {
            Result::Ok(val) => val,
            Result::Err(err) => panic!("Called 'Result::unwrap()' on an Err value: {:?}", err),
        }
    }
}
*/

fn double_number(number_str: &str) -> i32 {
    2 * number_str.parse::<i32>().unwrap()
}

fn main() {
     match pdouble_number("10") {
        Ok(m) => assert_eq!(m, 20),
        Err(err) => println!("Error: {:?}", err),
    }

   let n: i32 = double_number("10");
    assert_eq!(n, 2000);
}

use std::num::ParseIntError;

fn pdouble_number(number_str: &str) -> Result<i32, ParseIntError> {
    match number_str.parse::<i32> () {
        Ok(n) => Ok(2 * n),
        Err(err) => Err(err),
    }
}

fn ppdouble_number(number_str: &str) -> Result(i32, ParseIntError> {
    number_str.parse::<i32>().map(|n| 2 * n)
}
