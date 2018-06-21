use std::fs::File;
use std::io::prelude::*;

fn main() {

    fn open_file() -> std::io::Result<()> {
        let mut file = File::open("foo.txt")?;
        let mut contents = String::new();

        file.read_to_string(&mut contents)?;
        assert_eq!(contents, "Hello, world!");
        Ok(())
    }

    fn open_file_match() -> std::io::Result<()> {
        let f = File::open("foot.txt") ;
        let mut contents = String::new();
        match f {
            Ok(mut f) => {
                match f.read_to_string(&mut contents) {
                    Ok(_) => {
                        println!("{}", contents);
                        Ok(())
                    },
                    Err(e) => Err(e),
                }
            },
            Err(e) => {
                println!("Error happened when open file: {}", e);
                Err(e)
            },
        }
    }

    open_file().unwrap();
    open_file_match().unwrap();

}

