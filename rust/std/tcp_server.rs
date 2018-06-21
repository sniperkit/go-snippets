use std::io;
use std::net::{TcpListener, TcpStream};

fn main() {
    fn handle_client(stream: TcpStream) {
        println!("{:?}", stream)
    }

    fn start_server() -> std::io::Result<()>{
        let listener = TcpListener::bind("127.0.0.1:8080").unwrap();
        for stream in listener.incoming() {
            handle_client(stream?);
        }
        Ok(())
    }

    start_server().unwrap();
}


