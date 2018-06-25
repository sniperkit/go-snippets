#!/usr/bin/env python3

# Parser for DQA secondary reports

import csv
import json
import logging
import sys

logger = logging.getLogger(__name__)
logger.setLevel(logging.INFO)

handler = logging.StreamHandler()
handler.setLevel(logging.DEBUG)
logger.addHandler(handler)

# PEDSnet-v2-Colorado-ETLv2 has data model version listed as '2'
VERSION_MAP = {'v1': '1.0.0',
               'v2': '2.0.0',
               '2': '2.0.0'}

# take care of the known typos.
# Not sure if there's a better way to do this.
STATUS_MAP = {'soltution proposed': 'solution proposed',
              'solutoin proposed': 'solution proposed'}


class Parser():
    def __init__(self, file_obj,
                 site, table,
                 field_totals, table_totals,
                 site_table_totals, site_totals, expanded_site_totals):
        self.file_obj = file_obj
        self.site = site
        self.table = table
        self.field_totals = field_totals
        self.table_totals = table_totals
        self.site_table_totals = site_table_totals
        self.site_totals = site_totals
        self.expanded_site_totals = expanded_site_totals
        self.parsed = False

    def parse(self):
        if self.parsed:
            return self.field_totals

        reader = csv.reader(self.file_obj)

        # skip the header
        next(reader)

        for rec in reader:
            if not rec[7]:
                continue

            field = rec[5].strip().lower()
            code = rec[7].strip().upper()
            finding = rec[9].strip()
            prevalence = rec[10].strip().lower() or 'not specified'
            rank = rec[11].strip().lower()
            response = rec[12].strip()
            cause = rec[13].strip()
            status = rec[14].strip().lower() or 'not specified'

            if status in STATUS_MAP:
                status = STATUS_MAP[status]

            # goal and description are specified in rec[6] and rec[8]
            # respectively, but it's better to pull them from
            # Data-Quality/Dictionary/DCC_DQA_Dictionary.csv
            # to ensure consistensy and no typos

            record = {
                'finding': finding,
                'prevalence': prevalence,
                'rank': rank,
                'response': response,
                'cause': cause
            }

            site_record = record.copy()
            site_record['site'] = self.site

            if status not in self.field_totals[self.table][field][code]:
                self.field_totals[self.table][field][code][status] = [site_record]
            else:
                self.field_totals[self.table][field][code][status].append(site_record)

            field_record = record.copy()
            field_record['field'] = field

            if self.site not in self.table_totals[self.table][code][status]:
                self.table_totals[self.table][code][status][self.site] = [field_record]
            else:
                self.table_totals[self.table][code][status][self.site].append(field_record)

            if status not in self.site_table_totals[self.table][code]:
                self.site_table_totals[self.table][code][status] = [field_record]
            else:
                self.site_table_totals[self.table][code][status].append(field_record)

            if status not in self.site_totals:
                self.site_totals[status] = 1
            else:
                self.site_totals[status] += 1

            if status not in self.expanded_site_totals[self.table]:
                self.expanded_site_totals[self.table][status] = 1
            else:
                self.expanded_site_totals[self.table][status] += 1

    @staticmethod
    def get_version(file_obj):
        reader = csv.reader(file_obj)

        try:
            # first line is the header, skip it
            next(reader)

            version = next(reader)[1]
            if version in VERSION_MAP:
                version = VERSION_MAP[version]

            return version

        except Exception:
            return None


def main(fileobj):
    output = Parser(fileobj, 'unknown site', 'unknown_table').parse()
    json.dump(output, sys.stdout, indent=4, sort_keys=True)


if __name__ == '__main__':
    # No filename provided, read from stdin.
    if len(sys.argv) == 1:
        main(sys.stdin)
        sys.exit(0)

    with open(sys.argv[1]) as f:
        main(f)
