fn main() {
    let a = [0, 1, 2, 3, 4, 5, 6, 7];
    let complete = &a[..];
    let middle = &a[1..4];

    println!("Complete is {:?}", complete);
    println!("middle is {:?}", middle);
}
