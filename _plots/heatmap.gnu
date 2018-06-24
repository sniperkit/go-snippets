#!/usr/bin/gnuplot
# heatmap script plots the output of the heatmap utility

if(!exists("inputFile") || !exists("outputFile") || !exists("nrdatasets")) {
	print("The following params must be defined: inputFile, outputFile, nrdatasets");
	exit(0);
}

if(!exists("titleComment")) {
	titleComment=""
} else {
	titleComment=" (".titleComment.")"
}
#set palette gray
set title "Similarity Matrix heatmap ".titleComment
print outputFile
set terminal postscript eps size 7,5.0 enhanced color font 'Arial,34'
set output outputFile
set xlabel "Dataset index"
set ylabel "Dataset index"
set xrange [-0.5:nrdatasets-0.5]
set yrange [-0.5:nrdatasets-0.5]
plot inputFile u 2:1:3 w image
system("epstopdf ".outputFile." && rm ".outputFile)
