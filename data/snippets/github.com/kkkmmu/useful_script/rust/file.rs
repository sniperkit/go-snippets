use std::io::Write;

fn main() {
    let mut f = std::fs::File::create("foo.txt").expect("Couldn't create foo.tx");
    let buf = b"Hello file";
    let result = f.write(buf);

    open_or_create();
}

fn open_or_create() {
    use std::fs::File;
    use std::io::ErrorKind;

    let mut f = File::open("hello.txt");

    let mut f = match f {
        Ok(file) => file,
        Err(ref error) if error.kind() == ErrorKind::NotFound => {
            match File::create("hello.txt") {
                Ok(fc) => fc,
                Err(e) => {
                    panic!("Tried to create file but there was a problem: {:?}", e)
                },
            }
        },
        Err(error) => {
            panic!("There was a problem opening the file {:?}", error)
        },
    };

   let rs = f.write(b"Hello ni hao world!");

   match rs {
       Ok(_) => println!("Write to file sucess"),
       Err(e) => println!("Write to file with error: {:?}", e),
   }
}
