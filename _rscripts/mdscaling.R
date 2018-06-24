#!/usr/bin/Rscript
library(MASS)

args <- commandArgs(trailingOnly=TRUE)

if (length(args) < 2 ) { 
			cat("Please provide the similarities csv and the number of dimensions.\n") 
	quit()
}

#library(MASS)
d <- as.matrix(read.csv(args[1],header=FALSE))
k <- as.numeric(args[2])
#MAXITERATIONS <- 1000
#TOLERANCE <- 1e-4
#EPSILON <- 1e-10
#fit <- isoMDS(d, k=k, trace=FALSE, tol = TOLERANCE, maxit = MAXITERATIONS, p=2)
#cat(fit$stress)
#cat("\n")
#for(i in 1:nrow(fit$points)) {
#		for(j in 1:ncol(fit$points)) {
#				cat(fit$points[i,j])
#				cat(" ")
#		}
#		cat("\n")
#}
#

cmd <- cmdscale(d, k=k, eig=TRUE)
cat(cmd$GOF[2])
cat("\n")

# use only if cmdscaling only
#cat(cmd$GOF[2])
#cat("\n")
#for(i in 1:nrow(cmd$points)) {
#		for(j in 1:ncol(cmd$points)) {
#				cat(cmd$points[i,j])
#				cat(" ")
#		}
#		cat("\n")
#}



fit <-sammon(d,y=cmd$points, trace=FALSE, k=k)
cat(fit$stress)
cat("\n")
for(i in 1:nrow(fit$points)) {
		for(j in 1:ncol(fit$points)) {
				cat(fit$points[i,j])
				cat(" ")
		}
		cat("\n")
}
