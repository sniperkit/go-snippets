#!/usr/bin/gnuplot
# clustering script plots the output of the clustering util

if(!exists("inputFile") || !exists("outputFile")) {
	print("The following params must be defined: inputFile (list of datasets), outputFile");
	exit(0);
}

set title "Error degradation for different clusters"
set terminal postscript eps size 7,5.0 enhanced color font 'Arial,34'
#set output outputFile
set xlabel "Number of clusters"
set grid

#set logscale x

set style line 1 lt 1 pt 7 lc rgb "#000000" lw 4 ps 4
set style line 2 lt 1 pt 2 lc rgb "#333333" lw 4 ps 4
set style line 3 lt 1 pt 4 lc rgb "#777777" lw 4 ps 4
set style line 4 lt 1 pt 10 lc rgb "#999999" lw 4 ps 4

#plot for [f in inputFile] f u 2:4 w lp t system("basename ".f) ls 1
set output outputFile."-avg"
set ylabel "Average error difference"
plot for [i=1:words(inputFile)] word(inputFile,i) u 2:3 w lp t system("basename ".word(inputFile,i)) ls i
system("epstopdf ".outputFile."-avg && rm ".outputFile."-avg")

set output outputFile."-max"
set ylabel "Max error difference"
set key bottom left
plot for [i=1:words(inputFile)] word(inputFile,i) u 2:4 w lp t system("basename ".word(inputFile,i)) ls i
system("epstopdf ".outputFile."-max && rm ".outputFile."-max")

set output outputFile."-median"
set key top right
set ylabel "Median error difference"
plot for [i=1:words(inputFile)] word(inputFile,i) u 2:5 w lp t system("basename ".word(inputFile,i)) ls i
system("epstopdf ".outputFile."-median && rm ".outputFile."-median")
