use std::fs::File;
use std::io::BufReader;
use std::io::prelude::*;

fn main() {
    fn bufferred_read() -> std::io::Result<()> {
        let file = File::open("read_file.rs")?;
        let mut buf_reader = BufReader::new(file);
        let mut contents = String::new();
        buf_reader.read_to_string(&mut contents)?;
        //assert_eq!(contents, "Hello world!\n");
        Ok(())
    }

    fn read_file() -> std::io::Result<()> {
        let mut file = File::open("read_file.rs")?;
        let mut contents = String::new();

        file.read_to_string(&mut contents)?;
       // assert_eq!(contents, "Hello, world!");
        Ok(())
    }

    fn read_line() -> std::io::Result<()> {
        let mut f = File::open("read_file.rs")?;
        let mut reader = BufReader::new(f);

        let mut line = String::new();
        let len = reader.read_line(&mut line)?;

        println!("First line is {} and {} bytes long", line, len);

        Ok(())
    }

    bufferred_read().unwrap();
    read_file().unwrap();
    read_line().unwrap();
}
