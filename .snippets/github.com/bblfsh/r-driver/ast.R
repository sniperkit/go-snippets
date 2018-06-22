library(stringr)
suppressMessages(library(jsonlite))

tree = function(x) {
  type = typeof(x)
  value = NULL
  children = NULL
  
  if (is.atomic(x) && length(x) == 1) {
    value = deparse(x)[1]
  } else if (is.name(x)) {
    value = as.character(x)
    if (value == "") {
      value = "<empty>"
    }
  } else if (is.call(x)) {
    value = "call"
    children = lapply(as.list(x), tree)
  } else if (is.pairlist(x)) {
    language = "language"
    value = "list"
    branches = format(names(x))
    children = list()
    for (i in seq_along(x)) {
      children[[branches[i]]] = tree(x[[i]])
    }
  } else {
    value = format(x)
  }

  list(type=type, value=value, children=children)
}

stdin = file("stdin")
ast = tryCatch({
  lapply(parse(text=readLines(stdin)), tree)
}, finally={
  close(stdin)
})
print(toJSON(ast, auto_unbox=TRUE,
             pretty=Sys.getenv("PRETTY", "FALSE") == "TRUE"))

