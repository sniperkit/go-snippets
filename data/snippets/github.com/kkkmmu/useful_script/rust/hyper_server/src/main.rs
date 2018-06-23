extern crate hyper;
extern crate futures;

use futures::future::Future;
use hyper::header::ContentLength;
use hyper::server::{Http, Request, Response, Service};
use hyper::{Method, StatusCode};
use std::ascii::AsciiExt;
use futures::Stream;
use hyper::Chunk;

struct HelloWorld;

const PHRASE: &'static str = "Hello, Welcome to rust server!";

impl Service for HelloWorld {
    type Request = Request;
    type Response = Response;
    type Error = hyper::Error;

    type Future = Box<Future<Item=Self::Response, Error=Self::Error>>;

    fn call(&self, _req: Request) -> Self::Future {
        Box::new(futures::future::ok(
                Response::new()
                .with_header(ContentLength(PHRASE.len() as u64))
                .with_body(PHRASE)
                ))
    }
}

struct Echo;

fn to_uppercase(chunk: Chunk) -> Chunk {
    let uppered = chunk.iter().map(|byte| byte.to_ascii_uppercase()).collect::<Vec<u8>>();
    Chunk::from(uppered)
}

impl Service for Echo {
    type Request = Request;
    type Response = Response;
    type Error = hyper::Error;

    type Future = Box<Future<Item=Self::Response, Error=Self::Error>>;


    fn call(&self, req: Request) -> Self::Future {
        let mut response = Response::new();

        match (req.method(), req.path()) {
            (&Method::Get, "/") => {
                response.set_body("Try to Posting data to /echo");
            },
            (&Method::Post, "/echo") => {
                response.set_body(req.body());
            },
            _ => {
                response.set_status(StatusCode::NotFound);
            },
        };

        Box::new(futures::future::ok(response))
    }
}


fn main() {
    let addr = "127.0.0.1:8080".parse().unwrap();

    //let server = Http::new().bind(&addr, ||Ok(HelloWorld)).unwrap();
    let server = Http::new().bind(&addr, ||Ok(Echo)).unwrap();
    server.run().unwrap();
}

