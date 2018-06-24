#!/usr/bin/Rscript
# author: ggian
# date: Tue Sep 20 11:33:03 EEST 2016
#
# Script used to construct a Artificial Neural Network for classification. 
# The script expects a training file and a test file (used for the evaluation).

#for i in $(head  -n 1 $TRAINFILE | tr ',' '\n' | grep --invert-match "class"); do
#	let CLASS_INDEX=CLASS_INDEX+1
#done
#let CLASS_INDEX=CLASS_INDEX+1



library(nnet)
args <- commandArgs(trailingOnly=TRUE)

if (length(args) < 2 ) { 
	cat("Please provide the training set and test sets.\n") 
	quit()
}

datafile <- args[1]
testfile <- args[2]
train.data <- read.csv(datafile)
test.data <- read.csv(testfile)
ideal <- class.ind(train.data$class)

#ann = nnet(train.data[,-grep("class", colnames(train.data))], ideal, linout=TRUE, size=5, softmax=TRUE, trace=FALSE, MaxNWts=100000)
ann = nnet(class ~., data=train.data, linout=TRUE, size=5, trace=FALSE, MaxNWts=100000)

m<-cbind(test.data[,"class"],predict(ann,test.data[,-grep("class", colnames(test.data))],type="raw"))
cat(sum((m[,1]-m[,2])^2)/length(test.data))
#cat(sum(m[,1]!=m[,2])/length(test.data[,"class"])
