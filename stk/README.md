STK
======

**stk** is a command line tool that captures stderr and search for a reason on http://stackoverflow.com

```
$ stk mysql -u root -ptest
```

Will produce:
```
Warning: Using a password on the command line interface can be insecure.

Similar from Stackoverflow:
Bash Script Mysql Warning: Using a password on the command line interface can be insecure

Accepted solution:
Windows Server 2003 and later provide the where.exe program which does some of what which does, though it matches all types of files, not just executable commands.  (It does not match built-in shell commands like cd.)  It will even accept wildcards, so where nt* finds all files in your %PATH% and current directory whose names start with nt.

Try where /? for help.

Note that Windows PowerShell defines where as an alias for the Where-Object cmdlet, so if you want where.exe, you need to type the full name instead of omitting the .exe extension.


URL: http://stackoverflow.com/questions/23762575/bash-script-mysql-warning-using-a-password-on-the-command-line-interface-can-be
```


A project made for Gopher Gala by [Igor Vasilcovsky](https://github.com/vasilcovsky) and [Brian Oluwo](https://github.com/broluwo) using Go lang.

##Example Usage
![Usage](https://raw.githubusercontent.com/gophergala/stk/master/content/stk.gif?token=AChZnLq3CjjS9NXpaElipGZqqr6n5C6Uks5UzWD7wA%3D%3D)


##Dependencies (outside the std lib)
* [Kingpin](https://github.com/alecthomas/kingpin) - A Go (golang) command line and flag parser
* [Builder] (https://github.com/lann/builder) - fluent immutable builders for Go
