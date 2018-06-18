#!/bin/sh
# @author nkhilnani

if [ "$1" == "" ]
then
    echo ""
    echo "     USAGE 1: $0 FILENAME"
    echo ""
    echo "     Creates a tar of a directory's content only"
    echo "     Run from inside the directory to Tar"
    echo "     Tar file created in parent dir"
    echo ""    
	
else
	
	export COPY_EXTENDED_ATTRIBUTES_DISABLE=true
	export COPYFILE_DISABLE=true
	
	find $1 -name \*.DS_Store -exec rm -rf {} \;
	
	find $1 -name \*._* -exec rm -rf {} \;
	
	tar -cvf ../$1.tar *
	
	tar -tvf ../$1.tar

fi


	
