#!/usr/bin/env python

######################################################################
#
#    sql-runner.py
#    Copyright (C) 2014  Nik Khilnani   nik@Khilnani.org
#
#    This program is free software; you can redistribute it and/or modify
#    it under the terms of the GNU General Public License as published by
#    the Free Software Foundation; either version 2 of the License, or
#    (at your option) any later version.
#
#    This program is distributed in the hope that it will be useful,
#    but WITHOUT ANY WARRANTY; without even the implied warranty of
#    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
#    GNU General Public License for more details.
#
#    You should have received a copy of the GNU General Public License along
#    with this program; if not, write to the Free Software Foundation, Inc.,
#    51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
#
#######################################################################

# Looking for something more sophisticated? See: https://github.com/verigak/colors/

BLACK   = '\033[90m'
RED     = '\033[91m'
GREEN   = '\033[92m'
YELLOW  = '\033[93m'
BLUE    = '\033[94m'
MAGENTA = '\033[95m'
CYAN    = '\033[96m'
WHITE   = '\033[97m'

END = '\033[0m'

def black( msg ):
  print BLACK + msg + END

def red( msg ):
  print RED + msg + END

def green( msg ):
  print GREEN + msg + END

def yellow( msg ):
  print YELLOW + msg + END

def blue( msg ):
  print BLUE + msg + END

def magenta( msg ):
  print MAGENTA + msg + END

def cyan( msg ):
  print CYAN + msg + END

def white( msg ):
  print WHITE + msg + END

if __name__ == "__main__":
  black("black")
  red("red")
  green("green")
  yellow("yellow")
  blue("blue")
  magenta("magenta")
  cyan("cyan")
  white("white")
