#!/usr/bin/gnuplot
# size gnuplot script creates x-y plots that evaluate the impact of 
# dataset size to the similarity between datasets

if(!exists("inputFile") || !exists("outputFile")) {
	print("The following params must be defined: inputFile, outputFile");
	exit(0);
}


if(!exists("titleComment")) {
	titleComment="" 
} else {
	titleComment=" (".titleComment.")"
}

set title "Similarity vs dataset size".titleComment
set terminal postscript eps size 7,5.0 enhanced color font 'Arial,34'
set output outputFile
set xlabel "Dataset size"
set ylabel "Similarity"
set logscale x
set xrange[1000:]
set yrange [0:1]
set grid
set key bottom right

#set logscale x

set style line 1 lt 1 pt 7 lc rgb "#000000" lw 4 ps 4
set style line 2 lt 1 pt 2 lc rgb "#333333" lw 4 ps 4
set style line 3 lt 1 pt 4 lc rgb "#777777" lw 4 ps 4
set style line 4 lt 1 pt 10 lc rgb "#999999" lw 4 ps 4

#plot for [f in inputFile] f u 2:4 w lp t system("basename ".f) ls 1
#plot for [i=1:words(inputFile)] word(inputFile,i) u 2:3 w lp t system("basename ".word(inputFile,i)) ls i
plot inputFile u 1:2 w lp t col ls 1, inputFile u 1:3 w lp t col ls 2
system("epstopdf ".outputFile." && rm ".outputFile)
