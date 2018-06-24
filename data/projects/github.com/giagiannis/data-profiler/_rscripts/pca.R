#!/usr/bin/Rscript
# 
# Author: ggian
# Last updated: Tue Oct  4 16:02:00 EEST 2016
# This script analyzes an input dataset, and outputs the eigenvectors
# of the dataset, scaled by their respective variance (as determined by 
# the respective eigenvalue)

args <- commandArgs(trailingOnly=TRUE)
dataset <- read.csv(args[1])
cols = ncol(dataset)
fit <- prcomp(dataset)
variances=c()
sumOfVariances = sum(as.numeric(fit$sdev)^2)
for (i in 1:(cols)) {
	variances[i] = (fit$sdev[i]^2)/sumOfVariances
}


for (i in 1:(cols*cols)) {
	if(i %% cols == 1) {
		mul=1
		if (fit$rotation[i]<0) {
			mul=-1
			#mul=1
		}
	}
	cat((as.numeric(variances[(i-1)%/%cols + 1])*fit$rotation[i]*mul))
	cat("\t")
}


#
#
#fit <- princomp(dataset, cor=TRUE)
#
#
## prints the number of the eigenvalues/vectors
##cat(length(fit$sd)^2)
##cat("\n")
#
## prints the eigenvalues
#
#variances=c()
#sumOfVariances = sum(as.numeric(fit$sd)^2)
#for (i in 1:(cols)) {
#	variances[i] = (fit$sd[i]^2)/sumOfVariances
#}
#
## prints the eigenvectors
#for (i in 1:(cols*cols)) {
#	if(i %% cols == 1) {
#		mul=1
#		if (fit$loadings[i]<0) {
#			mul=-1
#			#mul=1
#		}
#	}
#	cat((as.numeric(variances[(i-1)%/%cols + 1])*fit$loadings[i]*mul))
#	cat(" ")
#}
