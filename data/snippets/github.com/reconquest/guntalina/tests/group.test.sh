#!/bin/bash

init

config <<DATA
actions:
    unique-action-foo:
        commands:
            - echo foo executed

    unique-action-bar:
        commands:
            - echo bar executed

rules:
    - masks:
        - /etc/file1
      workflow:
        - unique-action-foo
      group: orgalorg
    - masks:
        - /etc/*
      workflow:
        - unique-action-bar
      group: orgalorg
DATA

sources <<DATA
/etc/file1
/etc/file2
DATA


tests_ensure run

tests_ensure assert_diff <<DATA
following actions will be executed:
unique-action-foo
    echo foo executed

following commands will be executed:
echo foo executed

executing: echo foo executed
foo executed

DATA
