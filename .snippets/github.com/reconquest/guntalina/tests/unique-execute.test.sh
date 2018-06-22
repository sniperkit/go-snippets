#!/bin/bash

init

config <<DATA
actions:
    reload-software:
        commands:
            - echo executed

rules:
    - masks:
        - /etc/file1
      workflow:
        - reload-software
    - masks:
        - /etc/*
      workflow:
        - reload-software
DATA

sources <<DATA
/etc/file1
/etc/file2
DATA


tests_ensure run

tests_ensure assert_diff <<DATA
following actions will be executed:
reload-software
    echo executed

following commands will be executed:
echo executed

executing: echo executed
executed

DATA
