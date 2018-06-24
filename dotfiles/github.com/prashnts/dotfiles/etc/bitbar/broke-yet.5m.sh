#!/bin/sh
export PATH=$PATH:/usr/local/bin

balance=$(redis-cli get broke-yet | gnumfmt --format='%3.1f' --to=si)

echo "₹$balance" "| font=Input color=#242934"
