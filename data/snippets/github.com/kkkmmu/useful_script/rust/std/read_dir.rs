use std::io;
use std::fs::{self, DirEntry};
use std::path::Path;

fn main()  {
    fn printdir(dir :&DirEntry) {
        println!("{:?} {:?}", dir, dir.file_name());
        let ft = dir.file_type().unwrap();
        let mt = dir.metadata().unwrap();
        println!("{:?} {:?} {:?}", dir.path(), ft, mt)
    }
    fn visit_dirs(dir: &Path, cb: &Fn(&DirEntry)) -> io::Result<()> {
        if dir.is_dir() {
            for entry in fs::read_dir(dir)?{
                let entry = entry?;
                let path = entry.path();
                if path.is_dir() {
                    visit_dirs(&path, cb)?;
                } else {
                    cb(&entry);
                }
            }
        }
        Ok(())
    }

    visit_dirs(Path::new(".."), &printdir).unwrap()
}
