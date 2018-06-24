#!/bin/bash

eval go build main.go buffer.go parser.go graph.go settings.go
eval ./main >> result.txt
eval sort -k1 -n result.txt >> resultSort2.txt
eval clang++ -std=c++11 -O2 main.cpp -o mainD
eval ./mainD < resultSort2.txt >> resultSortD.txt
eval go build analyze.go
eval ./analyze < resultSortD.txt
