#!/usr/bin/python2.7
from __future__ import print_function
from mysql_regexp import *
from table_fields import *
import re


def parse_mysql(path):
    # Create the files to write the data
    table_file = {}

    for table in table_regexp.keys():
        table_file[table] = open("dataset/tables/" + table, "w")

    line_number = 0
    last_print = 0
    mysql = open(path)
    for l in mysql:
        # Print some feedback
        line_number += 1
        if (mysql.tell() - last_print) / (1024 * 1024) > 200:
            last_print = mysql.tell()
            print(last_print / (1024 * 1024), " megabytes read")
        try:
            if l.startswith("INSERT INTO"):
                table = l.split()[2][1:-1]
                if table in table_regexp:
                    for match in re.finditer(table_regexp[table], l):
                        fields = [match.group(f) for f in table_fields[table]]
                        print(",".join(fields), file=table_file[table])
        except Exception as e:
            print("error at line ", line_number, e, file=sys.stderr)


if __name__ == "__main__":
    parse_mysql("dataset/mysql")
