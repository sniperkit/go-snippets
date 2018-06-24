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
train.data <- read.csv(args[1])
test.data <- read.csv(args[2])


formula <- class ~ .
fit <- rpart(formula, data=train.data, method="anova")
pred <- predict(fit, test.data)
cat(sum((pred-test.data$class)^2)/length(test.data))
