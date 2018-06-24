#!/bin/bash
msg=$(date +%Y%m%d%H%M%S)
git add .
git commit -m $msg
git push origin master