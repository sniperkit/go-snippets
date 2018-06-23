fn main() {
    use std::fs::File;

    fn get_metadata() -> std::io::Result<()> {
        let mut f = File::open("foo.txt")?;
        let metadata = f.metadata()?;
        println!("{:?}", metadata);
        Ok(())
    }

    get_metadata().unwrap();
}
