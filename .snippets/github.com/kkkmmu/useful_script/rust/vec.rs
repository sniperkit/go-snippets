fn main() {
    let v = vec![1, 2, 3, 4, 5, 6];
    println!("v is {:?}", v);
    println!("v[2] is {}", v[2]);

    let v1 = vec!(7, 8, 9, 10, 11);
    println!("v1 is {:?}", v1);

    for i in &v1 {
        println!("i is {}", i);
    }

    for i in v1 {
        println!("Take ownership of the vector and its element {}", i);
    }

    /* 
    for i in v1 {
        println!("Take ownership of the vector and its element {}", i);
    }
    */
    let v2 = vec![0;20];
    println!("v2 is {:?}", v2);

    let v3 = vec![1, 2, 3];
    match v3.get(7) {
        Some(x) => println!("Item 7 is {}", x),
        None => println!("Sorrry, this vector is too short!"),
    }
}
