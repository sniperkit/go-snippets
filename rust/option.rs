fn find(haystack: &str, needle: char) -> Option<usize> {
    for (offset, c) in haystack.char_indices() {
        if c == needle {
            return Some(offset);
        }
    }
    None
}

fn main() {
    let file_name = "foobar.rs";

    match find(file_name, '.') {
        None => println!("No file extension found."),
        Some(i) => println!("File extension: {}", &file_name[i+1..]),
    }
}

fn option() {
    enum Option<T> {
        None,
        Some(T),
    }

    impl <T> Option<T> {
        fn unwrap(self) -> T {
            match self {
                Option::Some(val) => val,
                Option::None => panic!("Called 'Option::unwrap()' on a 'None' value"),
            }
        }
    }
}

fn extension_explicit(file_name: &str) -> Option<&str> {
    match find(file_name, '.') {
        None => None,
        Some(i) => Some(&file_name[i+1..]),
    }
}
