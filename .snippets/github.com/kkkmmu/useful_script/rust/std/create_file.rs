use std::fs::File;
use std::io::prelude::*;

fn main() {
    create_file_unwrap();
    create_file_no_match();
    create_file_match();
}

fn create_file_unwrap() {
    let mut f = File::create("file_unwrap.txt").unwrap();
    f.write_all(b"Hello");
}

fn create_file_no_match() -> std::io::Result<()> {
    let mut f = File::create("file_no_match.txt")?;
   f.write_all(b"file no match")?;

   Ok(())
}

fn create_file_match() -> std::io::Result<()> {
    let mut f = File::create("file_match.txt");
    match f {
        Ok(mut f) => f.write_all(b"Hello"),
        Err(e) => {
            println!("Error when create file: {}", e);
            Ok(())
        },
    }
}
