# -*- coding: utf-8 -*-
# Minimal & personal bookmark system
#
# Copyright (C) 2018-present Jeremies Pérez Morata
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

import click
import json
import shelve

_FIELDS = ['name', 'tags', 'desc']


@click.group()
def cli():
    pass


def load_db(db_path):
    db = {}
    with open(db_path) as f:
        db = json.loads(f.read())
    return db


def save_db(db, db_path):
    content = json.dumps(db, sort_keys=True, indent=2)
    with open(db_path, "w") as f:
        f.write(content)


def check_fields(entry):
    for field in _FIELDS:
        if field not in entry:
            entry[field] = None
    return entry


def find_entry (db, url):
    if url in db.keys():
        return check_fields(db[url])

    return check_fields({})


def get_input(field_name, field_data, is_list=False):
    str_values = ""
    if field_data is not None:
        if is_list:
            if type(field_data) is list and len(field_data) > 0:
                str_values = "[{0}]".format(" ".join(field_data))
        else:
            if field_data.strip() != "":
                str_values = "[{0}]".format(field_data)

    value = input("> {0}{1}: ".format(field_name.capitalize(), str_values)).strip()

    if value == "":
        value = field_data
    else:
        value = value

    if is_list and type(field_data) is not list:
        value = value.split()

    return value

@cli.command()
@click.option('--db-path', default="docs/db.json", help='Json DB File')
@click.argument('url')
def add(db_path, url):
    print("> Url:", url)

    db = load_db(db_path)

    entry = find_entry(db, url)

    entry["name"] = get_input("name", entry["name"], False)
    entry["tags"] = get_input("tags", entry["tags"], True)
    entry["desc"] = get_input("desc", entry["desc"], False)

    db[url] = entry

    save_db(db, db_path)

"""
@cli.command()
@click.option('--db-path', default="db.json", help='Json DB File')
@click.option('--out-path', default="./docs", help='Web OUT dir')
def generate(db_path, out_path):
    db = load_db(db_path)
"""
