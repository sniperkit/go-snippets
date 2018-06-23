extern crate futures;
extern crate hyper;
extern crate tokio_core;

use std::io::{self, Write};
use futures::{Future, Stream};
use hyper::Client;
use tokio_core::reactor::Core;

fn main() {

    let mut core = Core::new().unwrap();
    let client = Client::new(&core.handle());
    let uri = "http://www.163.com".parse().unwrap();

    let work = client.get(uri).map(|res| {
        println!("Response {}", res.status());

        res.body().for_each(|chunk| {
            io::stdout().write_all(&chunk).map(|_|()).map_err(From::from)
        })
    });

    core.run(work).unwrap();
}
