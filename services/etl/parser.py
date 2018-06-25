#!/usr/bin/env python3

# Parser for PEDSnet ETL conventions document

import json
import re
import sys
from logger import logger

title_re = re.compile(r'^\#[^\#]+', re.I)
table_name_re = re.compile(r'^\#\#\s*\d+\.\d+(?P<name>[ \w]*)', re.I)


class TableParser():
    def __init__(self, name, header_pattern, field_pattern):
        self.name = name
        self.header_re = re.compile(header_pattern, re.I)
        self.field_re = re.compile(field_pattern, re.X | re.I)

    def __repr__(self):
        return '<TableParser: {}>'.format(self.name)

    def matches_header(self, line):
        return self.header_re.match(line) is not None

    def parse_row(self, line):
        m = self.field_re.match(line)

        if not m:
            return

        name = m.group('name').lower().strip()
        req = m.group('req').lower().strip().startswith('yes')
        data_type = m.group('type').strip()
        desc = m.group('desc').strip()
        conv = m.group('conv')

        if conv:
            conv = conv.strip()
        else:
            conv = ''

        return (name, req, data_type, desc, conv)


table_parsers = (
    TableParser(
        name='v1',
        header_pattern=r'^field\s*\|\s*required',
        field_pattern=r'''^
            (?P<name>[^\|]+)      # name of the field
            \|(?P<req>[^\|]+)     # network required
            \|(?P<type>[^\|]+)    # data type
            \|(?P<desc>[^\|]+)    # description
            (\|(?P<conv>[^\|]+))? # optional conventions
        '''),

    TableParser(
        name='v2',
        header_pattern=r'^field\s*\|\s*(foreign key/|not null)',
        field_pattern=r'''^
            (?P<name>[^\|]+)      # name of the field
            \|(?P<fk>[^\|]+)      # foreign key constraint
            \|(?P<req>[^\|]+)     # network required
            \|(?P<type>[^\|]+)    # data type
            \|(?P<desc>[^\|]+)    # description
            (\|(?P<conv>[^\|]+))? # optional conventions
        '''),
)


def _get_table_name(l):
    return table_name_re.match(l).group('name').lower().strip().replace(' ', '_')  # noqa


class Parser():
    def __init__(self, fileobj):
        self.model = {
            'content': '',
            'tables': {}
        }

        self.fileobj = fileobj
        self.parsed = False
        self.line = ''

    def parse(self):
        if self.parsed:
            return self.model

        self.get_content()

        while self.get_next_table():
            pass

        self.parsed = True

        return self.model

    def add_content(self, content):
        logger.debug('[parser] added top-level content')

        self.model['content'] = content

    def add_table(self, table, content=''):
        if table in self.model['tables']:
            raise KeyError('table %s already added' % table)

        logger.debug('[parser] added table %s', table)

        self.model['tables'][table] = {
            'content': content,
            'fields': {}
        }

    def add_table_content(self, table, content):
        logger.debug('[parser] added table %s content', table)
        self.model['tables'][table]['content'] = content

    def add_field(self, table, field, required, data_type, desc, conv):
        if field in self.model['tables'][table]['fields']:
            raise KeyError('field %s already added for table %s', field, table)

        logger.debug('[parser] added field %s/%s', table, field)

        self.model['tables'][table]['fields'][field] = {
            'required': required,
            'data_type': data_type,
            'description': desc,
            'etl_conventions': conv
        }

    def get_content(self):
        # Skip everything up to and including the title line.
        for line in self.fileobj:
            if title_re.match(line):
                break

        # Everything up to the first table is model-level content.
        self.line = ''
        content = []

        # When the first table name is found, break and queue the line.
        for line in self.fileobj:
            if table_name_re.match(line):
                logger.debug('[parser] found table name: %s', line)
                self.line = line
                break

            content.append(line)

        self.add_content(''.join(content).strip())

    def get_next_table(self):
        '''
        Read from the currently cached line and up to the start of the next
        table (or EOF), and parse current table's content.
        '''

        # The document should now have the table name line cached.
        # If it doesn't, there are no more tables to read.
        if self.line == '':
            return False

        # Table name line looks something like this: "## 1.1 PERSON".
        # Parse out the actual table name.
        table_name = _get_table_name(self.line)
        self.add_table(table_name)

        # Everything after the table name line and up to the table header
        # is the content (documentation) of the table.
        content = []
        parser = None

        for line in self.fileobj:
            for parser in table_parsers:
                if parser.matches_header(line):
                    logger.debug('[parser] matched header for %s', table_name)
                    break
            else:
                parser = None

            if parser:
                break

            content.append(line)
        else:
            line = None

        # If a line is present, a table header matched. Skip the
        # line '--- | --- | --- | --- | ---' that follows.
        if line is not None:
            logger.debug('[parser] parsing fields for %s', table_name)
            next(self.fileobj)

            # Parse the rows in the table.
            for line in self.fileobj:
                field = parser.parse_row(line)

                # If a valid field is parsed, add and continue. Otherwise
                # set the line for the next iteration.
                if not field:
                    self.line = line
                    break

                self.add_field(table_name, *field)

        # Read the rest of the table contents, whatever comes after
        # the fields specification table, up to the next table.
        if not table_name_re.match(self.line):
            content.append(line)
            self.line = ''

            for line in self.fileobj:
                if table_name_re.match(line):
                    self.line = line
                    break

                content.append(line)

        self.add_table_content(table_name, ''.join(content).strip())

        return True


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
