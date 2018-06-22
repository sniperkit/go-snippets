template parameter assumptions:
- `.*` means to match any characters there, except what follows immediately, and include the contents/counts/etc as for filter analysis
- multi-line match is assumed in places where the param is specified between a beginning and end character match

base assumptions across all languages that do not need to be defined in templates:
- white space is any blank line that contains no characters or only white space characters
- loc is anything that is not a comment or white space

gaps/XXX:

- cases where there are mixed languages in a single file (maybe ignore forever?)
- finding end of a def in python, match on next line that begins with non-white-space w/o including it
