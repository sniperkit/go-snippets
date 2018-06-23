struct HasDrop;

impl Drop for HasDrop {
    fn drop(&mut self) {
        println!("Dropping");
    }
}

fn firework () {
    struct Firework {
        strength : i32,
    }

    impl Drop for Firework {
        fn drop(&mut self) {
            println!("BOOM times {}!!!!", self.strength);
        }
    }

    let fircracker = Firework{strength: 1};
    let tnt = Firework{ strength: 1000};
}

struct Droppable {
    name: &'static str,
}

impl Drop for Droppable {
    fn drop(&mut self) {
        println!("> Dropping {} ", self.name);
    }
}

fn main() {
    let _a = Droppable{name: "a"};

    {
        let _b = Droppable{name: "b"};
        {
            let _c = Droppable{name: "c"};
            let _d = Droppable{name: "d"};
            println!("Exiting block B");
        }
        println!("Just exited block B");
        println!("Exiting block A");
    }
    println!("Just exited block A");

    drop(_a);

    println!("end of the main function");

    let x = HasDrop;

    firework();

}
