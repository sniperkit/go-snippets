#!/usr/bin/Rscript
# author: ggian
# date: Tue Sep 20 11:54:23 EEST 2016
# 
# Script used to calculate a CART.
# Limitation: both the trainsets and the testsets MUST enumerate their classes
# starting from 1 (not 0) & their class attribute must be named "class" 
# and be positioned as the last column on the csv files.


library(rpart)
args <- commandArgs(trailingOnly=TRUE)
if (length(args) < 2 ) { 
	cat("Please provide the training set and test sets.\n") 
	quit()
}
datafile <- args[1]
testfile <- args[2]
train.data <- read.csv(datafile)
test.data <- read.csv(testfile)
cls.idx <- grep("class", colnames(train.data))


formula <- class ~ .
fit <- rpart(formula, train.data, method="class")
m<-cbind(test.data[,cls.idx], predict(fit, test.data, type="class"))
#print(m)
cat(sum(m[,1]!=m[,2])/length(test.data[,cls.idx]))
