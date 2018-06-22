# Copyright 2015 Open Source Robotics Foundation, Inc.
# Copyright 2013 Willow Garage, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from collections import namedtuple
import re

from parse_cmake import list_utils

QuotedString = namedtuple('QuotedString', 'contents comments')
_Arg = namedtuple('Arg', 'contents comments')
_Command = namedtuple('Command', 'name body comment')
BlankLine = namedtuple('BlankLine', '')

BEGIN_BLOCK_COMMANDS = [
    'function',
    'macro',
    'if',
    'else',
    'elseif',
    'foreach',
    'while'
]
END_BLOCK_COMMANDS = [
    'endfunction',
    'endmacro',
    'endif',
    'else',
    'elseif',
    'endforeach',
    'endwhile'
]


class FormattingOptions():
    """Specifies the formatting options for pretty-printing CMakeLists.txt output.
       The default configuration aims to match the house style used
       by the CMake project itself. See https://github.com/Kitware/CMake
    """
    def __init__(self):
        self.indent = '  '
        self.max_line_width = 79


class File(list):
    """Top node of the syntax tree for a CMakeLists file."""

    def pretty_print(self, formatting_opts=FormattingOptions()):
        '''
        Returns the pretty-print string for tree
        with indentation given by the string tab.
        '''
        return '\n'.join(compose_lines(self, formatting_opts)) + '\n'

    def __str__(self):
        return self.pretty_print()

    def __repr__(self):
        return 'File(' + repr(list(self)) + ')'


class Comment(str):
    def __repr__(self):
        return 'Comment(' + str(self) + ')'


def Arg(contents, comments=None):
    return _Arg(contents, comments or [])


def Command(name, body, comment=None):
    return _Command(name, body, comment)


class CMakeParseError(Exception):
    pass


def prettify(s, formatting_opts=FormattingOptions()):
    """
    Returns the pretty-print of the contents of a CMakeLists file.
    """
    return parse(s).pretty_print(formatting_opts)


def parse(s, path='<string>'):
    '''
    Parses a string s in CMakeLists format whose
    contents are assumed to have come from the
    file at the given path.
    '''
    nums_toks = tokenize(s)
    nums_items = list(parse_file(nums_toks))
    nums_items = attach_comments_to_commands(nums_items)
    items = [item for _, item in nums_items]
    return File(items)


def strip_blanks(tree):
    return File([x for x in tree if not isinstance(x, BlankLine)])


def compose_lines(tree, formatting_opts):
    """
    Yields pretty-printed lines of a CMakeLists file.
    """
    tab = formatting_opts.indent
    level = 0
    for item in tree:
        if isinstance(item, (Comment, str)):
            yield level * tab + item
        elif isinstance(item, BlankLine):
            yield ''
        elif isinstance(item, _Command):
            name = item.name.lower()
            if name in END_BLOCK_COMMANDS:
                level -= 1

            for i, line in enumerate(command_to_lines(item, formatting_opts)):
                offset = 1 if i > 0 else 0
                line2 = (level + offset) * tab + line
                yield line2

            if name in BEGIN_BLOCK_COMMANDS:
                level += 1


def is_parameter_name_arg(name):
    return re.match('^[A-Z_]+$', name) and name not in ['ON', 'OFF']


def command_to_lines(cmd, formatting_opts, use_multiple_lines=False):
    class output:
        lines = []
        current_line = cmd.name.lower() + '('
        is_first_in_line = True

    def end_current_line():
        output.lines += [output.current_line]
        output.current_line = ''
        output.is_first_in_line = True

    for arg_index, arg in enumerate(cmd.body):
        # when formatting a command to multiple lines, try to start
        # new lines with parameter names
        #
        #   command(FOO arg
        #     OPTION value
        #     OPTION value)
        if arg_index > 0 and use_multiple_lines and is_parameter_name_arg(arg.contents):
            end_current_line()

        arg_str = arg_to_str(arg).strip()
        if len(output.current_line) + len(arg_str) > formatting_opts.max_line_width:
            if not use_multiple_lines:
                # if the command does not fit on a single line, re-enter the function
                # in multi-line formatting mode so that we can choose the best
                # points to break the line
                return command_to_lines(cmd, formatting_opts, use_multiple_lines=True)
            else:
                end_current_line()

        if output.is_first_in_line:
            output.is_first_in_line = False
        else:
            output.current_line += ' '

        output.current_line += arg_str
        if len(arg.comments) > 0:
            end_current_line()

    output.current_line += ')'

    if cmd.comment:
        output.current_line += ' ' + cmd.comment

    end_current_line()

    return output.lines


def arg_to_str(arg):
    comment_part = '  ' + '\n'.join(arg.comments) + '\n' if arg.comments else ''
    return arg.contents + comment_part


def parse_file(toks):
    '''
    Yields line number ranges and top-level elements of the syntax tree for
    a CMakeLists file, given a generator of tokens from the file.

    toks must really be a generator, not a list, for this to work.
    '''
    prev_type = 'newline'
    for line_num, (typ, tok_contents) in toks:
        if typ == 'comment':
            yield ([line_num], Comment(tok_contents))
        elif typ == 'newline' and prev_type == 'newline':
            yield ([line_num], BlankLine())
        elif typ == 'word':
            line_nums, cmd = parse_command(line_num, tok_contents, toks)
            yield (line_nums, cmd)
        prev_type = typ


def attach_comments_to_commands(nodes):
    return list_utils.merge_pairs(nodes, command_then_comment, attach_comment_to_command)


def command_then_comment(a, b):
    line_nums_a, thing_a = a
    line_nums_b, thing_b = b
    return (isinstance(thing_a, _Command) and
            isinstance(thing_b, Comment) and
            set(line_nums_a).intersection(line_nums_b))


def attach_comment_to_command(lnums_command, lnums_comment):
    command_lines, command = lnums_command
    _, comment = lnums_comment
    return command_lines, Command(command.name, command.body[:], comment)


def parse_command(start_line_num, command_name, toks):
    cmd = Command(name=command_name, body=[], comment=None)
    expect('left paren', toks)
    for line_num, (typ, tok_contents) in toks:
        if typ == 'right paren':
            line_nums = range(start_line_num, line_num + 1)
            return line_nums, cmd
        elif typ == 'left paren':
            pass
            # raise ValueError('Unexpected left paren at line %s' % line_num)
        elif typ in ('word', 'string'):
            cmd.body.append(Arg(tok_contents, []))
        elif typ == 'comment':
            c = tok_contents
            if cmd.body:
                cmd.body[-1].comments.append(c)
            else:
                cmd.comments.append(c)
    msg = 'File ended while processing command "%s" started at line %s' % (
        command_name, start_line_num)
    raise CMakeParseError(msg)


def expect(expected_type, toks):
    line_num, (typ, tok_contents) = next(toks)
    if typ != expected_type:
        msg = 'Expected a %s, but got "%s" at line %s' % (
            expected_type, tok_contents, line_num)
        raise CMakeParseError(msg)

# http://stackoverflow.com/questions/691148/pythonic-way-to-implement-a-tokenizer
# TODO: Handle multiline strings.
scanner = re.Scanner([
    (r'#.*', lambda scanner, token: ("comment", token)),
    (r'"[^"]*"', lambda scanner, token: ("string", token)),
    (r"\(", lambda scanner, token: ("left paren", token)),
    (r"\)", lambda scanner, token: ("right paren", token)),
    (r'[^ \t\r\n()#"]+', lambda scanner, token: ("word", token)),
    (r'\n', lambda scanner, token: ("newline", token)),
    (r"\s+", None),  # skip other whitespace
])


def tokenize(s):
    """
    Yields pairs of the form (line_num, (token_type, token_contents))
    given a string containing the contents of a CMakeLists file.
    """
    toks, remainder = scanner.scan(s)
    line_num = 1
    if remainder != '':
        msg = 'Unrecognized tokens at line %s: %s' % (line_num, remainder)
        raise ValueError(msg)
    for tok_type, tok_contents in toks:
        yield line_num, (tok_type, tok_contents.strip())
        line_num += tok_contents.count('\n')
