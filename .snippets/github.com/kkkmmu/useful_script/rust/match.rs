enum Message {
    Quit,
    ChangeColor(i32, i32, i32),
    Move{x: i32, y: i32},
    Write(String),
}

fn quit() {
    println!("Quiting...........!");
}

fn change_color(r: i32, g: i32, b:i32) {
    println!("Changing Color........!");
}

fn move_cursor(x: i32, y: i32) {
    println!("Moving cursor.........!");
}

fn process_message(msg: Message) {
    match msg {
        Message::Quit => quit(),
        Message::ChangeColor(r, g, b) => change_color(r, g, b),
        Message::Move{x, y: new_name_for_y} => move_cursor(x, new_name_for_y),
        Message::Write(s) => println!("{}", s),
    };
}

fn main() {
    let ms = Message::Quit;
    process_message(ms);
    let ms = Message::ChangeColor(10, 11, 12);
    process_message(ms);
    let ms = Message::Move{x:100, y:1000};
    process_message(ms);

    let ms = Message::Write(String::from("liwei"));
    process_message(ms);
}
