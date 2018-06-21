//macro_rules声明一个宏，say_hello是宏名字，say_hello后边是一对大括号，大括号中间是宏的内容。
//宏的内容中=>前边的()表示这个宏没有输入参数。在整个宏的内部，我们调用的println宏。
macro_rules! say_hello {
    () => (
        println!("Hello");
        )
}

//宏的内容除了可以放在大括号中以外还可以放在()和[]中。
//这里要注意的是当宏内容放在大括号中时，末尾不用放“;”，而小括号和中括号后要放";".

macro_rules! say_hello1 (
    () => ( 
        println!("Hello one");
        )
 );

macro_rules! say_hello2 [
    () => (
        println!("Hello two");
        )
];

//宏的参数是包含在宏的内容中的。
//宏定义中宏名称后边大括号中的内容称为macro rule(宏规则)。一个宏可以包含任意多个宏规则。
//宏规则的格式是这样的:
//  () => {};
//  () => [];
//  () => ();
//每个宏规则有一个=>。
//  在=>左边的小括号是该规则的输入参数，该参数的形式是一个叫做matcher的匹配器。
//  在=>右边的的括号是该规则对应的宏内容，需要执行或者被替换，叫做transcriber(转换器)，可以是语句，表达式等。
//
//within a matcher, each metavariable has a fragment specifier, identifying which syntatic for it
//matches.
//  1. ident: an identifier. Examples: x; foo.
//  2. path:  a qualified name. Example: T::SpecialA.
//  3. expr:  an expresssion. Examples: 2 + 2; if true {1} else {0};.
//  4. ty:    a type. Examples: i32; Vec<(char, String)>; &T,
//  5. pat:   a pattern. Examples: Some(t); (17, 'a');_.
//  6. stmt:  a single statement. Example: let x = 3.
//  7. block: a brace-delimited sequence of statements and optionally an expression. Example:
//          { log(error, "hi"); return 12;}.
//  8. item:  an item. Exmaples: fn foo() {}; struct Bar;
//  9. meta:  a 'meta item", as found in attributes. Example: cfg(target_os="windows").
//  10. tt:   a single token tree.
macro_rules! print_hello {
    ($msg:expr) =>  (
        println!("Hello {}", $msg);
        )
}

//宏的输入不是通过参数列表传入的，而是通过match的方式进行查找或者匹配的。
macro_rules! foo {
    (x => $e:expr) => (println!("Mode X: {}", $e));
    (y => $e:expr) => (println!("Mode Y: {}", $e));
}


//宏名字的解析与函数略有不同，宏的定义必须出现在宏调用之前，即与C里面函数类似---函数定义或声明必须出现在函数调用之前，只不过Rust宏没有单纯的声明，所以宏调用之前需要先定义。
//宏调用虽然与函数调用很像，但是宏的名字与函数名字是处于不同的命名空间的。
fn bar(x: i32) -> i32 {
    x * x
}

macro_rules! bar {
    ($x:ident) => { println!("{:?}", $x) } 
}

//宏里面的变量都是以$开头的，其余的都是按字面去匹配，以$开头的变量都是用来表示语法元素的，为了限定匹配什么类型的语法元素，需要用指示符（designator）加以限定，就跟普通的变量绑定一样用冒号将变量和类型分开。当前宏支持一下几种指示符：
//  1. indet: 标识符，用来表示函数或变量名。
//  2. expr: 表达式。
//  3. block： 代码块，用花括号包起来的多个语句。
//  4. pat: 模式，普通模式匹配。
//  5. path: 路径，注意这里不是操作系统中的文件路径，而是双冒号分割的限定符。
//  6. tt： 单个语法树。
//  7. ty：类型，语义层面的类型，如i32,char。
//  8. item: 条目。
//  9. meta: 元条目。
//  10. stmt: 单条语句。
//
//
//
//  Rust的函数不能接受任意多个参数，但是宏可以,这是通过重复（repetition）来实现的。
//  除了重复之外，宏还支持递归，即在宏定义时调用其自身，类似于递归函数。

fn main() {
    say_hello!();
    say_hello1!();
    say_hello2!();
    print_hello!("liwei");

    foo!(x => 7);
    foo!(y => 3);
    let x = 10;
    println!("{:?}", bar(x));
    println!("{:?}", bar!(x));
}

