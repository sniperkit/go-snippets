#!/usr/bin/python
import csv
import sys

if len(sys.argv)<2:
    sys.exit()

with open(sys.argv[1], "rb") as f:
    reader = csv.reader(f)
    summary = 0.0
    for row in reader:
        try:
            summary += float(row[len(row)-1])
        except ValueError:
            pass

print summary
