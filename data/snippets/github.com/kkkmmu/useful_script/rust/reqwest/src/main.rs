extern crate reqwest;

use std::io::Read;

fn main() {
    let mut resp = reqwest::get("https://www.github.com").unwrap();

    assert!(resp.status().is_success());

    let mut content = String::new();
    resp.read_to_string(&mut content).unwrap();

    println!("{:?}", content);
}
