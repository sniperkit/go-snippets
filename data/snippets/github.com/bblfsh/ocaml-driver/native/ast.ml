(* ocamlbuild -tag debug -use-ocamlfind ast.native *)

let _ =
  let cc = open_in "expr.ml" in
  let ll = Lexing.from_channel cc in
  let ast = Parse.implementation ll in
  let json = Astjson.structure_to_yojson ast in
  print_endline (Yojson.Safe.to_string json)
