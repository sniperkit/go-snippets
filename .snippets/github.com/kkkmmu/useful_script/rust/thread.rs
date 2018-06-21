use std::thread;

fn main() {
  let handle =  thread::spawn(|| {
        println!("Hello from a thread");
    });

  println!("{:?}", handle.join().unwrap());
}
