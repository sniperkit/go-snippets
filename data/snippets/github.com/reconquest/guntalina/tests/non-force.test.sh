#!/bin/bash

init

config <<DATA
actions:
    action-foo:
        commands:
            - bash -c 'echo foo executed, forcing fail ; exit 2'
    action-bar:
        commands:
            - echo bar executed

rules:
    - masks:
        - /etc/file1
      workflow:
        - action-foo
        - action-bar

DATA

sources <<DATA
/etc/file1
DATA


tests_ensure run

tests_ensure assert_diff <<DATA
following actions will be executed:
action-foo
    bash -c 'echo foo executed, forcing fail ; exit 2'

action-bar
    echo bar executed

following commands will be executed:
bash -c 'echo foo executed, forcing fail ; exit 2'
echo bar executed

executing: bash -c 'echo foo executed, forcing fail ; exit 2'
foo executed, forcing fail

exit status 2

DATA
