#![allow(dead_code)]
#![feature(rustc_private)]
extern crate syntax;
extern crate syntax_pos;
extern crate rustc_errors;
extern crate rustparse;
extern crate serialize;
extern crate serialize as rustc_serialize;

use syntax_pos::MultiSpan;
use rustc_errors::Level;
use rustc_errors::snippet::Style;
use syntax::ast::Crate;
use rustc_errors::diagnostic::Diagnostic;
use std::io::{BufRead, stdin};
use serialize::json::as_json;
use rustparse::{parse_source, ParsedAST, ParseStatus};
use rustc_serialize::json::decode;

/// Input received via stdin from the caller of the driver.
#[derive(RustcDecodable, RustcEncodable)]
struct ParseInput {
    content: String,
}

/// Result that is sent via stdout to the consumer of the driver.
#[derive(RustcEncodable)]
struct ParseOutput { 
    ast: Option<Crate>,
    errors: Vec<Diagnostic>,
    status: ParseStatus,
}

impl ParseOutput {
    fn from_ast(ast: ParsedAST) -> ParseOutput {
        ParseOutput{
            ast: ast.ast,
            errors: ast.errors,
            status: ast.status,
        }
    }

    fn from_error(msg: String) -> ParseOutput {
        ParseOutput{
            ast: None,
            errors: vec![Diagnostic{
                level: Level::Fatal,
                message: vec![(msg, Style::Level(Level::Fatal))],
                span: MultiSpan::new(),
                children: vec![],
                code: None,
            }],
            status: ParseStatus::Fatal,
        }
    }
}

fn report_error(msg: String) {
    println!("{}", as_json(&ParseOutput::from_error(msg)));
}

fn handle_request(content: String) -> Result<ParseOutput, String> {
    let input: ParseInput = match decode(content.as_str()) {
        Ok(i) => i,
        Err(_) => {
            return Err("error decoding input from json".into());
        },
    };

    Ok(ParseOutput::from_ast(parse_source(input.content)))
}

fn main() {
    let stdin = stdin();
    let mut handle = stdin.lock();
    loop {
        let mut line = String::new();
        if handle.read_line(&mut line).is_err() {
            continue;
        }

        if line.len() == 0 {
            return;
        }
        
        match handle_request(line) {
            Ok(output) => println!("{}", as_json(&output)),
            Err(err) => report_error(err),
        }
    }
}

#[test]
fn test_handle_request() {
    let result = handle_request(format!("{}", as_json(&ParseInput{
        content: r#"
        #[derive(Debug, Clone, Copy)]
        pub struct Point {
            x: i64,
            y: i64,
        }

        impl Point {
            pub fn new(x: i64, y: i64) -> Point {
                Point{x: x, y: y}
            }
        }
        "#.into(),
    })));
    assert!(result.is_ok());

    let output = result.unwrap();
    println!("{}", as_json(&output));
    assert_eq!(output.errors.len(), 0);
    assert!(output.ast.is_some());
}

#[test]
fn test_handle_request_fail() {
    let result = handle_request("wow such json, very well-formed".into());
    assert!(result.is_err());
}
