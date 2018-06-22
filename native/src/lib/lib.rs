#![feature(rustc_private)]
extern crate syntax;
extern crate rustc_errors;
extern crate serialize;

use std::rc::Rc;
use std::cell::RefCell;
use syntax::ast::Crate;
use syntax::codemap::CodeMap;
use syntax::parse::{ParseSess, parse_crate_from_source_str};
use rustc_errors::{Level, Handler};
use rustc_errors::emitter::Emitter;
use rustc_errors::diagnostic::Diagnostic;
use rustc_errors::diagnostic_builder::DiagnosticBuilder;
use serialize::{Encoder, Encodable};

/// A diagnostics emitter that will add all the diagnostics to its internal
/// list of diagnostics.
/// Only diagnostics that are errors will be added.
struct DiagnosticsEmitter {
    diagnostics: Rc<RefCell<Vec<Diagnostic>>>,
}

impl DiagnosticsEmitter {
    /// Creates a new emitter that will add the emitted diagnostics to the
    /// passed list of diagnostics.
    pub fn new(diagnostics: Rc<RefCell<Vec<Diagnostic>>>) -> DiagnosticsEmitter {
        DiagnosticsEmitter {
            diagnostics: diagnostics,
        }
    }
}

impl Emitter for DiagnosticsEmitter {
    fn emit(&mut self, db: &DiagnosticBuilder) {
        let d = db.clone().into_diagnostic();
        match d.level {
            Level::Bug |
            Level::Fatal |
            Level::PhaseFatal |
            Level::Error => {
                self.diagnostics.borrow_mut().push(d);
            }
            _ => {}
        }
    }
}

/// Representation of the parsed AST with its status, the diagnostics gathered
/// during the parsing.
#[derive(Debug)]
pub struct ParsedAST {
    pub ast: Option<Crate>,
    pub status: ParseStatus,
    pub errors: Vec<Diagnostic>,
}

/// The status of the parsing. A status `Ok` will mean the pasing went well,
/// `Error` means that compilation succeeded but error diagnostics were emitted
/// and `Fatal` means no ast was generated, only error diagnostics.
#[derive(Debug,Eq,PartialEq)]
pub enum ParseStatus {
    Ok,
    Error,
    Fatal
}

impl Encodable for ParseStatus {
    fn encode<S: Encoder>(&self, s: &mut S) -> Result<(), S::Error> {
        match *self {
            ParseStatus::Ok => s.emit_str("ok"),
            ParseStatus::Error => s.emit_str("error"),
            ParseStatus::Fatal => s.emit_str("fatal"),
        }
    }
}

impl ParsedAST {
    fn new(krate: Option<Crate>, errors: Vec<Diagnostic>) -> ParsedAST {
        ParsedAST{
            ast: krate.clone(),
            status: match krate {
                Some(_) => if errors.len() > 0 {
                    ParseStatus::Error
                } else {
                    ParseStatus::Ok
                },
                None => ParseStatus::Fatal,
            },
            errors: errors,
        }
    }
}

/// Parses the given source and returns a sturcture with the parsed AST, all
/// the diagnostics emitted during the parsing.
pub fn parse_source(source: String) -> ParsedAST {
    let diagnostics: Rc<RefCell<_>> = Rc::new(RefCell::new(Vec::new()));
    let emitter = DiagnosticsEmitter::new(diagnostics.clone());
    let sh = Handler::with_emitter(false, false, Box::new(emitter));
    let ps = ParseSess::with_span_handler(sh, Rc::new(CodeMap::new()));

    let result = parse_crate_from_source_str("stdin".into(), source, &ps);
    match result {
        Ok(krate) => {
            let errors = diagnostics.borrow();
            ParsedAST::new(Some(krate), (*errors).to_vec())
        },
        Err(mut d) => {
            d.emit();
            let errors = diagnostics.borrow();
            ParsedAST::new(None, (*errors).to_vec())
        },
    }
}

#[cfg(test)]
fn to_ident(item: &syntax::ast::Item) -> String {
    format!("{}", item.ident.name)
}

#[test]
fn test_parse_source_ok() {
    let result = parse_source("pub fn meaning_of_life() -> i32 {\n\t42\n}".into());
    assert_eq!(result.status, ParseStatus::Ok);
    assert_eq!(result.errors.len(), 0);

    let items = result.ast.unwrap().module.items;
    assert_eq!(items.len(), 1);
    assert_eq!(to_ident(&items.get(0).unwrap()), "meaning_of_life");
}

#[test]
fn test_parse_source_error() {
    let result = parse_source("pub fn meaning_of_life() -> i32 {\n\tthis->is_not->C\n} pub fn meaning_of_foo() -> f64 { this->is->not->C }".into());
    assert_eq!(result.status, ParseStatus::Error);
    assert_eq!(result.errors.len(), 2);

    let items = result.ast.unwrap().module.items;
    assert_eq!(items.len(), 2);
    assert_eq!(to_ident(&items.get(0).unwrap()), "meaning_of_life");
    assert_eq!(to_ident(&items.get(1).unwrap()), "meaning_of_foo");
}

#[test]
fn test_parse_source_fatal() {
    let result = parse_source("pub fn meaning_of_life() i32 {\n\t42\n}".into());
    assert_eq!(result.status, ParseStatus::Fatal);
    assert_eq!(result.errors.len(), 1);
    assert_eq!(result.ast.is_none(), true);
}
