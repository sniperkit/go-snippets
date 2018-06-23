use std::{thread, time};
use std::time::Duration;

fn main() {
    park_thread();
    one_handred_thread();
    thread::sleep(time::Duration::from_millis(100));
}

fn park_thread() {
    let parked_thread = thread::Builder::new().spawn(||{
        println!("Parking thread");
        thread::park();
        println!("Thread unparked");
    }).unwrap();

    thread::sleep(Duration::from_millis(10));
    println!("Unpark the thread");
    parked_thread.thread().unpark();
    parked_thread.join().unwrap();
}

fn one_handred_thread() {
    let mut i = 0;
    loop {
        thread::Builder::new().spawn(move ||{
            println!("thread {}", i)
        });

        i+= 1;
        
        if i > 99 {
            break;
        }
    }
}
