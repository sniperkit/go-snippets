#!/usr/bin/gnuplot
# script used to plot a number of columns from a given file

if(!exists("inputFile") || !exists("outputFile") || !exists("cols") || !exists("xtitle") || !exist("ytitle"))  {
	print("The following params must be defined: inputFile, outputFile, cols, xtitle, ytitle");
	exit(0);
}


set xlabel xtitle
set ylabel ytitle

if (exists("titleComment")) {
	set title ytitle." vs ".xtitle." (".titleComment.")" 
} else {
	set title ytitle." vs ".xtitle 
}

set terminal postscript eps size 7,5.0 enhanced color font 'Arial,34'
set output outputFile
set grid

set style line 1 lt 1 pt 7 lc rgb "#000000" lw 4 ps 4
set style line 2 lt 1 pt 2 lc rgb "#333333" lw 4 ps 4
set style line 3 lt 1 pt 4 lc rgb "#777777" lw 4 ps 4
set style line 4 lt 1 pt 10 lc rgb "#999999" lw 4 ps 4
set style line 5 lt 1 pt 13 lc rgb "#aaaaaa" lw 4 ps 4
set key top right

#plot for [i=1:words(cols)] for [j=1:words(inputFile)] word(inputFile,j) u 1:(column(word(cols,i)+0)) w lp t col(word(cols,i)) ls (i+j-1)
plot for [i=1:words(cols)] for [j=1:words(inputFile)] word(inputFile,j) u 1:(column(word(cols,i)+0)) w lp t system("basename ".word(inputFile,j)) ls (i+j-1)
system("epstopdf ".outputFile." && rm ".outputFile)
