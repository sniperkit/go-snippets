fn guess(n: i32) -> bool {
    if n < 1 || n > 10 {
        panic!("Invalid number: {}", n);
    }

    n == 5
}

fn double() {
    use std::env;

    let mut argv = env::args();
    let arg: String = argv.nth(1).unwrap();
    let n: i32 = arg.parse().unwrap();
    println!("{}", 2 * n);
}

fn main() {
    eextension();
    double();
    guess(11);
}

fn find(haystack: &str, needle: char) -> Option<usize> {
    for (offset, c) in haystack.char_indices() {
        if c == needle {
            return Some(offset);
        }
    }
    None
}

fn eextension() {
    let file_name = "foobar.rs";
    match find(file_name, '.') {
        None => println!("No file extension found."),
        Some(i) => println!("File extension: {}", &file_name[i+1..]),
    }
}

/*
   enum Option<T> {
   None,
   Some(T),
   }

   impl <T> Option<T> {
   fn unwrap(self) -> T {
   match self {
   Option::Some(val) => val,
   Option::None => panic!("called 'Option::unwrap' on a 'None' value!"),
   }
   }
   }
   */

fn file_path_ext_explicit(file_path: &str) -> Option<&str> {
    match file_name(file_path) {
        None => None,
        Some(name) => match extension(name) {
            None => None,
            Some(ext) => Some(ext),
        }
    }
}

fn file_name(file_path: &str) -> Option<&str> {
    unimplemented!()
}

fn extension(file_name: &str) -> Option<&str> {
    find(file_name, '.').map(|i| &file_name[i+1..])
}

fn and_then<F, T, A>(option: Option<T>, f: F) -> Option<A> where F: FnOnce(T) -> Option<A> {
    match option {
        None => None,
        Some(value) => f(value),
    }
}

fn file_path_ext(file_path: &str) -> Option<&str> {
    file_name(file_path).and_then(extension)
}
