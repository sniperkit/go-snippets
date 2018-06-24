#!/usr/bin/Rscript
# 
# Author: ggian
# Last updated: Thu Jan 26 10:40:03 EET 2017
# This script prints the ranks of an input dataset (used for analysis)

args <- commandArgs(trailingOnly=TRUE)
dataset <- read.csv(args[1])

s <- seq(1,nrow(dataset),1)
rs <- rev(s)
r <- rank(dataset)
cat(sqrt(sum((r-s)^2)))
