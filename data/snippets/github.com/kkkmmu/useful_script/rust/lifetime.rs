fn main() {
    /*
    fn skip_prefix(line: &str, prefix: &str) -> &str {

    }
    */

    fn skip_prefix<'a, 'b>(line : &'a str, prefix: &'b str) -> &'a str{
        line
    }

    let line = "lang:en=Hello World!";
    let lang = "en";

    let v ;
    {
        let p = format!("lang:{}=", lang);
        v = skip_prefix(line, p.as_str());
    }

    println!("{}", v);
}

/*
 * Here we have a function skip_prefix which takes two &str references as parameters and returns a
 * single &str reference. We call it by passing in references to line and p: Two variables with
 * different lifetimes. Now the safety of the println!-line depends on whether the reference
 * returned by skip_prefix function references the still living line or the already dropped p
 * string.
 *
 * Because of the above ambiguity, Rust will refuse to compile the example code. To get it to
 * compile we need to tell the compiler more about the lifetimes of the references. This can be
 * done by making the lifetimes explicit in the function declaration:
 *
 * fn skip_prefix<'a, 'b>(line: &'a str, prefix: &'b str) -> &'a str {
 *     // ...
 *     }
 * *
 */

/*
 * Note It's important to understand that lifetime annotations are descriptive, not prescriptive.
 * This means that how long a reference is valid is determined by the code, not by the annotations.
 * The annotations, however, give information about lifetimes to the compiler that uses them to
 * check the validity of references. The compiler can do so without annotations in simple cases,
 * but needs the programmer's support in complex scenarios.
 *
 * */
