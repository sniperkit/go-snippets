#!/usr/bin/Rscript
# author: ggian
# date: Tue Sep 20 11:33:03 EEST 2016
#
# Script used to construct a Artificial Neural Network for classification. 
# The script expects a training file and a test file (used for the evaluation).
library(neuralnet)
args <- commandArgs(trailingOnly=TRUE)

if (length(args) < 2 ) { 
	cat("Please provide the training set and test sets.\n") 
	quit()
}

layers <- c(5,5)
datafile <- args[1]
testfile <- args[2]
train.data <- read.csv(datafile)
test.data <- read.csv(testfile)
atts<-attributes(train.data[1:length(train.data)-1])
form <- "class ~"
for(x in atts$names) form <- paste0(form,"+", x)
fit <- neuralnet(as.formula(form), data = train.data, hidden=layers, stepmax=1e6)
cmp <- compute(fit, test.data)
for (i in 1:length(cmp$net.result)) {
		cat(paste0(cmp$net.result[i],"\n"))
}
