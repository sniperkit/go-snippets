#!/bin/bash

init

config <<DATA
actions:
    do-something-1:
        commands:
            - echo command_1

    do-something-2:
        commands:
            - echo command_2
rules:
    - masks:
        - /etc/something/something1*.conf
      workflow:
        - do-something-1

    - masks:
        - /etc/something/something2*.conf
      workflow:
        - do-something-2
DATA

sources <<DATA
/etc/something/something1.conf
/etc/something/something2_with_suffix.conf
DATA


tests_ensure run

tests_ensure assert_diff <<DATA
following actions will be executed:
do-something-1
    echo command_1

do-something-2
    echo command_2

following commands will be executed:
echo command_1
echo command_2

executing: echo command_1
command_1

executing: echo command_2
command_2
DATA
