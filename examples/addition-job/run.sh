#!/bin/bash

# Capture command line arguments
INPUT=$1
OUTPUT=$2

echo ''
echo '----------------------------------------------------'
echo 'Calling job with arguments ' ${INPUT} ${OUTPUT}
SCRIPT=my_alg.py

python ${SCRIPT} ${INPUT} ${OUTPUT}
rc=$?
echo 'Done calling job - wrapper finished'
echo '----------------------------------------------------'
echo ''
exit ${rc}
