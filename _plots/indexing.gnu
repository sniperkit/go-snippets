#!/usr/bin/gnuplot
# indexing gnuplot script visualizes the results of the indexing util

if(!exists("inputFile") || !exists("outputFile")) {
	print("The following params must be defined: inputFile, outputFile");
	exit(0);
}


if(!exists("titleComment")) {
	titleComment="" 
} else {
	titleComment=" (".titleComment.")"
}

set title "Distance between estimated and optimal points ".titleComment
set terminal postscript eps size 7,5.0 enhanced color font 'Arial,34'
set output outputFile
set xlabel "k"
set ylabel "Distance"
#set logscale x
#set xrange[1000:]
#set yrange [0:1]
#set logscale x
set grid
set key top right

set style line 1 lt 1 pt 7 lc rgb "#000000" lw 4 ps 4
set style line 2 lt 1 pt 2 lc rgb "#333333" lw 4 ps 4
set style line 3 lt 1 pt 4 lc rgb "#777777" lw 4 ps 4
set style line 4 lt 1 pt 10 lc rgb "#999999" lw 4 ps 4
set style line 5 lt 1 pt 13 lc rgb "#aaaaaa" lw 4 ps 4

plot for [i=1:words(inputFile)] word(inputFile,i) u 1:2 w lp t system("basename ".word(inputFile,i)) ls i
system("epstopdf  --outfile=".outputFile."-dst-.pdf ".outputFile." && rm ".outputFile)

set title "Stress factor".titleComment
set terminal postscript eps size 7,5.0 enhanced color font 'Arial,34'
set output outputFile
set xlabel "k"
set ylabel "Stress"
#set logscale x
#set xrange[1000:]
#set yrange [0:1]
#set logscale x
set grid
set key top right

set style line 1 lt 1 pt 7 lc rgb "#000000" lw 4 ps 4
set style line 2 lt 1 pt 2 lc rgb "#333333" lw 4 ps 4
set style line 3 lt 1 pt 4 lc rgb "#777777" lw 4 ps 4
set style line 4 lt 1 pt 10 lc rgb "#999999" lw 4 ps 4
set style line 5 lt 1 pt 13 lc rgb "#aaaaaa" lw 4 ps 4

plot for [i=1:words(inputFile)] word(inputFile,i) u 1:3 w lp t system("basename ".word(inputFile,i)) ls i
system("epstopdf  --outfile=".outputFile."-stress.pdf ".outputFile." && rm ".outputFile)

set title "Indexing time".titleComment
set terminal postscript eps size 7,5.0 enhanced color font 'Arial,34'
set output outputFile
set xlabel "k"
set ylabel "Time (sec)"
#set logscale x
#set xrange[1000:]
#set yrange [0:1]
#set logscale x
set grid
set key top right

set style line 1 lt 1 pt 7 lc rgb "#000000" lw 4 ps 4
set style line 2 lt 1 pt 2 lc rgb "#333333" lw 4 ps 4
set style line 3 lt 1 pt 4 lc rgb "#777777" lw 4 ps 4
set style line 4 lt 1 pt 10 lc rgb "#999999" lw 4 ps 4
set style line 5 lt 1 pt 13 lc rgb "#aaaaaa" lw 4 ps 4

plot for [i=1:words(inputFile)] word(inputFile,i) u 1:4 w lp t system("basename ".word(inputFile,i)) ls i
system("epstopdf  --outfile=".outputFile."-time.pdf ".outputFile." && rm ".outputFile)
