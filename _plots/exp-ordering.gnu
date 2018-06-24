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

set title "Kendall {/Symbol t} vs Sampling Rate".titleComment
set terminal postscript eps size 7,5.0 enhanced color font 'Arial,34'
set output outputFile
set xlabel "Sampling Rate"
set grid
set key top right

set style line 1 lt 1 pt 7 lc rgb "#000000" lw 4 ps 4
set style line 2 lt 1 pt 2 lc rgb "#333333" lw 4 ps 4
set style line 3 lt 1 pt 4 lc rgb "#777777" lw 4 ps 4
set style line 4 lt 1 pt 10 lc rgb "#999999" lw 4 ps 4
set style line 5 lt 1 pt 13 lc rgb "#aaaaaa" lw 4 ps 4

set ylabel "Kendall {/Symbol t}"
plot inputFile u 1:2 w lp t system("basename ".inputFile) ls 1
system("epstopdf --outfile=".outputFile."-tau.pdf ".outputFile." && rm ".outputFile)


set title "top-k hits vs Sampling Rate".titleComment
set ylabel "Percentage of top-k hits"
set output outputFile
plot inputFile u 1:8 w lp t "top-10%" ls 1, \
'' u 1:14 w lp t "top-25%" ls 2, \
'' u 1:20 w lp t "top-50%" ls 3
system("epstopdf --outfile=".outputFile."-topk.pdf ".outputFile." && rm ".outputFile)


set title "% needed to contain appx % score vs Sampling Rate".titleComment
set ylabel "% of actual list"
set output outputFile
plot inputFile u 1:26 w lp t "top-2%" ls 1, \
'' u 1:32 w lp t "top-5%" ls 2, \
'' u 1:38 w lp t "top-10%" ls 3
system("epstopdf --outfile=".outputFile."-topk-perc.pdf ".outputFile." && rm ".outputFile)


set ylabel "Pearson {/Symbol r}"
set title "Pearson {/Symbol r} vs Sampling Rate".titleComment
set output outputFile
plot inputFile u 1:56 w lp t system("basename ".inputFile) ls 1
system("epstopdf --outfile=".outputFile."-rho.pdf ".outputFile." && rm ".outputFile)
