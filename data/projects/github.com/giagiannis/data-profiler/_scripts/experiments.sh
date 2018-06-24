#!/bin/bash
# script used to execute the experiment


# variables that need to be defined before execution
echo $GOPATH
[ "$DATASET_PATH" == "" ] && echo "DATASET_PATH variable needed" && exit 0
[ "$OPERATORS" == "" ] && echo "OPERATORS variable needed" && exit 0

[ "$MDSSCRIPT" == "" ] && export MDSSCRIPT="$GOPATH/src/github.com/giagiannis/data-profiler/_rscripts/mdscaling.R"
[ "$SVMSCRIPT" == "" ] && export SVMSCRIPT="$GOPATH/src/github.com/giagiannis/data-profiler/_rscripts/svm-appx.R"

[ "$NO_PARTITIONS" == ""  ] && NO_PARTITIONS=32 && echo "Using default NO_PARTITIONS=$NO_PARTITIONS" 
[ "$REPETITIONS" == ""  ] && REPETITIONS=10 && echo "Using default REPETITIONS=$REPETITIONS" 
[ "$THREADS" == ""  ] && THREADS=10 && echo "Using default THREADS=$THREADS" 
[ "$TREE_SR" == ""  ] && TREE_SR=0.10 && echo "Using default TREE_SR=$TREE_SR" 
[ "$MDS_K" == ""  ] && MDS_K=5 && echo "Using default MDS_K=$MDS_K" 
[ "$SR" == ""  ] && SR="$(seq -s, 0.05 0.01 0.2)" && echo "Using default SR=$SR"
MDS_K_MAX=30

[ "$SCORES_PATH" == "" ] && echo "SCORES_PATH is empty - the operators will be executed exhaustively"


# initialization of en
b="$(basename $DATASET_PATH)/$NO_PARTITIONS-$TREE_SR-$MDS_K"
[ ! -d "$b" ] && mkdir -p $b && echo "Creating $b"
BASE="$b/base"
[ ! -d "$BASE" ]  && mkdir -p $BASE 2>/dev/null && echo "Creating $BASE"
RESULTS_DIR="$b/experiments"
[ ! -d "$RESULTS_DIR" ] && mkdir -p $RESULTS_DIR 2>/dev/null && echo "Creating $RESULTS_DIR"
PLOTS_DIR="$b/plots"
[ ! -d "$PLOTS_DIR" ] && mkdir -p $PLOTS_DIR 2>/dev/null && echo "Creating $PLOTS_DIR"
echo -e "##########################################\n\n"

function prepare_dataset() {
		echo "Preparing the dataset"
		LOGFILE="$BASE/log"
		echo "Calculating similarities"
		data-profiler-utils similarities -i $DATASET_PATH -o $BASE/matrix.sim -opt concurrency=$THREADS,partitions=$NO_PARTITIONS,tree.sr=$TREE_SR -l $LOGFILE
		echo "Calculating MDS (2 dimensions)"
		data-profiler-utils mds -k 2 -l $LOGFILE -m coords -o $BASE/mds-2 -sc $MDSSCRIPT -sim $BASE/matrix.sim
		echo "Calculating MDS"
		data-profiler-utils mds -k $MDS_K -l $LOGFILE -m coords -o $BASE/mds -sc $MDSSCRIPT -sim $BASE/matrix.sim
		echo "Calculating MDS (evaluating goodness of fit)"
		data-profiler-utils mds -k $MDS_K_MAX -l $LOGFILE -m gof -o $BASE/mds -sc $MDSSCRIPT -sim $BASE/matrix.sim

		echo "Running operators"
		for op in $OPERATORS; do 
			b="$(basename $op)"
			SCORESFILENAME="$b-scores"
			echo -e "\tOperator $op"
			if [ -f "$SCORES_PATH/$SCORESFILENAME" ]; then
					ln -v  $SCORES_PATH/$SCORESFILENAME $BASE/$SCORESFILE
			else
					data-profiler-utils train -i $DATASET_PATH -l $LOGFILE -p $THREADS -o $BASE/$SCORESFILENAME -s $op -t "nothing"
			fi
		done
		echo -e "##########################################\n\n"
}

function execute_experiments() {
		LOGFILE="$RESULTS_DIR/log"

		echo "Executing accuracy and ordering experiments"
		for op in $OPERATORS; do 
				echo -e "\tOperator $op"
				b="$(basename $op)"
				CMD="data-profiler-utils exp-accuracy -c $BASE/mds-coords -i $DATASET_PATH -l $LOGFILE -ml $SVMSCRIPT -o $RESULTS_DIR/$b-accuracy.csv -r $REPETITIONS -s $BASE/$b-scores -sr $SR -t $THREADS"
				$CMD 1>/dev/null
		done
		echo -e "##########################################\n\n"

}

function plot_experiments() {
		GNUPLOT_ACCURACY="/home/giannis/Projects/data-profiler/_plots/exp-accuracy.gnu"
		GNUPLOT_GOF="/home/giannis/Projects/data-profiler/_plots/mds-gof.gnu"
		GNUPLOT_COORDS="/home/giannis/Projects/data-profiler/_plots/mds-coords.gnu"

		echo "Generating mds-coords plot"
		gnuplot -e "inputFile='$BASE/mds-2-coords';outputFile='$PLOTS_DIR/mds-coords'" $GNUPLOT_COORDS
		echo "Generating mds-gof plot"
		gnuplot -e "inputFile='$BASE/mds-gof';outputFile='$PLOTS_DIR/mds'" $GNUPLOT_GOF 
		echo "Generating experiments plots"
		for op in $OPERATORS; do
				echo -e "\tOperator $op"
				b=$(basename $op)
				gnuplot -e "inputFile='$RESULTS_DIR/$b-accuracy.csv';outputFile='$PLOTS_DIR/$b'" $GNUPLOT_ACCURACY 
		done
		echo -e "##########################################\n\n"
}

prepare_dataset
execute_experiments
plot_experiments
