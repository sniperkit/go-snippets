#!/usr/bin/gnuplot
# mds-stress visualizes the results of the mds util (stress module)

if(!exists("inputFile") || !exists("outputFile")) {
	print("The following params must be defined: inputFile, outputFile");
	exit(0);
}


if(!exists("titleComment")) {
	titleComment="" 
} else {
	titleComment=" (".titleComment.")"
}

set style line 1 lt 1 pt 7 lc rgb "#000000" lw 4 ps 4
set style line 2 lt 1 pt 2 lc rgb "#333333" lw 4 ps 4
set style line 3 lt 1 pt 4 lc rgb "#777777" lw 4 ps 4
set style line 4 lt 1 pt 10 lc rgb "#999999" lw 4 ps 4
set style line 5 lt 1 pt 13 lc rgb "#aaaaaa" lw 4 ps 4



set title "Goodness of fit vs MDS dimensionality".titleComment
set terminal postscript eps size 7,5.0 enhanced color font 'Arial,34'
set output outputFile."-gof"
set xlabel "Dimensions"
set ylabel "GOF"
set grid
set key bottom right

plot for [i=1:words(inputFile)] word(inputFile,i) u 1:2 w lp t system("basename ".word(inputFile,i)) ls i
system("epstopdf --outfile=".outputFile."-gof.pdf ".outputFile."-gof && rm ".outputFile."-gof")

set title "Sammon Stress vs MDS dimensionality".titleComment
set terminal postscript eps size 7,5.0 enhanced color font 'Arial,34'
set output outputFile."-stress"
set xlabel "Dimensions"
set ylabel "Sammon Stress"
set grid
set key top right


plot for [i=1:words(inputFile)] word(inputFile,i) u 1:3 w lp t system("basename ".word(inputFile,i)) ls i
system("epstopdf --outfile=".outputFile."-stress.pdf ".outputFile."-stress && rm ".outputFile."-stress")
