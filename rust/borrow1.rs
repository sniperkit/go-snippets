fn main() {
    fn sum_vec(v: &Vec<i32>) -> i32 {
        v.iter().fold(0, |a, &b| a + b)
    }

    fn foo(v1: &Vec<i32>, v2: &Vec<i32>) -> i32 {
        let s1 = sum_vec(v1);
        let s2 = sum_vec(v2);

        s1 + s2
    }

    let v1 = vec![1, 3, 4];
    let v2 = vec![6, 7, 9];

    let answer = foo(&v1, &v2);

    println!("v1 is {:?}, v2 is {:?}", v1, v2);
    println!("answer is {}", answer);
}
