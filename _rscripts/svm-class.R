#!/usr/bin/Rscript
# author: ggian
# date: Mon Oct 24 11:40:15 EEST 2016
#
# Script used to create a linear model used for regression


library(e1071)
args <- commandArgs(trailingOnly=TRUE)
if (length(args) < 2 ) { 
	cat("Please provide the training set and test sets.\n") 
	quit()
}

datafile <- args[1]
testfile <- args[2]
train.data <- read.csv(datafile)
test.data <- read.csv(testfile)
#head(train.data)
#head(test.data)

fit <- svm(train.data[,-grep("class",colnames(train.data))], as.factor(train.data$class))
pred <- predict(fit, test.data[,-grep("class",colnames(test.data))])
#summary(fit)
#table(pred,test.data$class)

m<-cbind(as.factor(test.data$class),pred)
cat(sum(m[,1]!=m[,2])/length(test.data$class))
