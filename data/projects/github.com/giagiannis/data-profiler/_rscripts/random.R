#!/usr/bin/Rscript
# 
# Author: ggian
# Last updated: Tue Oct 11 12:45:57 EEST 2016
# This script creates a random n^2-dimensional vector. To be used for 
# sanity checking of the analysis methods.

args <- commandArgs(trailingOnly=TRUE)
datafile <- args[1]
dataset <- read.csv(datafile)
cols = ncol(dataset)
fit <- princomp(dataset, cor=TRUE)


# prints the number of the eigenvalues/vectors
cat(runif(cols,0,1))
