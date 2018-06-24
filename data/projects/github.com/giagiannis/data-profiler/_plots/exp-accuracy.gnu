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

set terminal postscript eps size 7,5.0 enhanced color font 'Arial,34'
set xlabel "Sampling Rate"
set grid

set style line 1 lt 1 pt 7 lc rgb "#000000" lw 4 ps 4
set style line 2 lt 1 pt 2 lc rgb "#333333" lw 4 ps 4
set style line 3 lt 1 pt 4 lc rgb "#777777" lw 4 ps 4
set style line 4 lt 1 pt 10 lc rgb "#999999" lw 4 ps 4
set style line 5 lt 1 pt 13 lc rgb "#aaaaaa" lw 4 ps 4
set style line 6 lt 1 pt 14 lc rgb "#cccccc" lw 4 ps 4
set key top right font ",15" samplen 1 outside horizontal

colNames=system("head -n 1 ".word(inputFile,1))
do for [t=2:67] {
	colName=word(colNames,t)
	set title colName." vs Sampling Rate".titleComment
	set ylabel colName
	set output outputFile
	if (t%3 == 1) {
		plot for [i=1:words(inputFile)] word(inputFile,i) u 1:t w lp t system("basename ".word(inputFile,i)) ls i
	} 
	if (t%3 == 2) {
		plot for [i=1:words(inputFile)] word(inputFile,i) u 1:t:t+1 w e t system("basename ".word(inputFile,i)) ls i
	}
	if (t%3 > 0 ) {
		system("epstopdf --outfile=".outputFile."-".colName.".pdf ".outputFile." && rm ".outputFile)
	}
}
