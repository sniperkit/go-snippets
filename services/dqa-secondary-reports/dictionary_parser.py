#!/usr/bin/env python3

# Parser for the dictionary of DQA status codes
# listed in Data-Quality/Dictionary/DCC_DQA_Dictionary.csv

import csv
import json
import sys


class Parser():
    def __init__(self, file_obj):
        self.file_obj = file_obj
        self.results = {}
        self.parsed = False

    def parse(self):
        if self.parsed:
            return self.results

        reader = csv.reader(self.file_obj)

        # skip the header
        next(reader)

        for rec in reader:
            goal = rec[0].strip().lower()
            code = rec[1].strip().upper()
            desc = rec[2]

            self.results[code] = {
                'goal': goal,
                'desc': desc
            }

        return self.results


def main(fileobj):
    output = Parser(fileobj).parse()
    json.dump(output, sys.stdout, indent=4, sort_keys=True)


if __name__ == '__main__':
    # No filename provided, read from stdin.
    if len(sys.argv) == 1:
        main(sys.stdin)
        sys.exit(0)

    with open(sys.argv[1]) as f:
        main(f)
