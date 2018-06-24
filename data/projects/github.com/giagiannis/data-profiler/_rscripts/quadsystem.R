#!/usr/bin/Rscript
args <- commandArgs(trailingOnly=TRUE)
if (length(args) < 1 ) { 
    cat("Please provide the points along with their distances\n") 
	quit()
}

library(nleqslv)

dataset <- read.csv(args[1])
n <- nrow(dataset) # number of points and dimensionality of space
coords <- as.matrix(dataset[,1:n])
distances <- as.matrix(dataset[n+1])

func <- function(x) {
		y <- numeric(n)
		for(i in 1:n) y[i] <- sqrt(sum((x-coords[i,])^2)) - distances[i]
		y
}

# start from a random points in space
noOfPoints <- 100
maxDist <- max(distances)
xstart=matrix(runif(noOfPoints*n,max=max(coords)+maxDist,min=min(coords)-maxDist),nrow=noOfPoints,ncol=n)

# create plot before estimating the roots
pdf(paste0(args[1],as.numeric(Sys.time()), ".pdf"))
plot(coords, 
	 xlim=c(min(coords)-maxDist, max(coords)+maxDist), 
	 ylim=c(min(coords)-maxDist, max(coords)+maxDist),
	 main=args[1])
theta <- seq(0,2*pi, length=50)
for(i in 1:n) lines(x = distances[i]*cos(theta)+coords[i,1], y=distances[i]*sin(theta)+coords[i,2])
grid(NULL,NULL, col="black")
points(xstart, pch=3)

zeros <- searchZeros(xstart,func,digits=2)

# no solution was found
if(is.null(zeros)){
		cat("-1")
		quit()
}
solutions <- nrow(zeros$x)
cat(solutions)
cat("\n")
for(i in 1:nrow(zeros$x)) {
		for(j in 1:ncol(zeros$x)) {
				cat(zeros$x[i,j])
				cat(" ")
		}
		cat("\n")
}
cat(zeros$xfnorm,sep=" ")

# plot the solutions
points(zeros$x, pch=7)

# used only for one solution
#nlq <- nleqslv(xstart, dslnex)

#cat(nlq$x)
#cat("\n")
#cat(nlq$termcd)
