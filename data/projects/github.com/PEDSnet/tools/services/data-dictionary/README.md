# Annotated data dictionary service

Service that exposes a single endpoint for fetching and parsing the PCORnet annotated data dictionary.

## Run

```
main.py [--token=<token>] [--host=<host>] [--port=<port>] [--debug]
```

Options:

- `host`
- `port`
- `token`, alternately the `GITHUB_AUTH_TOKEN` environment variable can be set.


# Parser for annotated data dictionary

Parses an annotated data dictionary excel spreadsheet and extracts the site comments and implementation status.

## Installation

To install the prerequisites, run
```
pip3 install -r requirements.txt 
```

## Run
```
./parser data_dictionary.xlsx
```
See `tests/expected.json` for example output.

Note: when running the parser, you may see the following openpyxl warning:
```
UserWarning: Discarded range with reserved name
```
This is expected and is not a concern. 
When openpyxl parses an excel workbook, it can only extract names that refer to cell ranges. However, the workbook, depending on its origin, may include other types of names, usually printer settings. These names start with '_xlmn' and openpyxl discards them.
To check what names were discarded, you can unzip the excel spreadsheet into a new directory, go into that directory and run
```
grep -ri _xlnm .
```

