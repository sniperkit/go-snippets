open Yojson.Basic
open Yojson.Basic.Util

type response = {
  status: string;
  errors: string list;
  ast: Astjson.structure;
}
[@@deriving yojson]

let print_response resp =
  resp |> response_to_yojson |> Yojson.Safe.to_string |> print_endline

let _ =
  while true do
    try
      let line_json = from_string (input_line stdin) in
      let encoding = line_json |> member "encoding" |> to_string
      and content = line_json |> member "content" |> to_string in
      let src = match encoding with
        | "UTF8" -> content
        | "Base64" -> B64.decode content
        | _ -> invalid_arg encoding
      in
      let ast = Parse.implementation (Lexing.from_string src) in
      print_response {status = "ok"; errors = []; ast}
    with
    | End_of_file -> ()
    | e ->
      let error = Printexc.to_string e in
      print_response {status = "fatal"; errors = [error]; ast = []};
  done
