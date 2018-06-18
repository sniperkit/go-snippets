#!/bin/sh
# @author nkhilnani

if [ "$1" == "" ]
then
    echo ""
    echo "     USAGE 1: $0 DIRNAME"
    echo ""    
    echo "     Removes .svn dirs" 
    echo "     Creates a tar of a directory"   
    echo "     Run from outside the directory to Tar"
    echo ""    
	
else
	
	export COPY_EXTENDED_ATTRIBUTES_DISABLE=true
	export COPYFILE_DISABLE=true
	
	find $1 -name \*.DS_Store -exec rm -rf {} \;
	
	find $1 -name \*._* -exec rm -rf {} \;

  find $1 -name \*.svn -exec rm -rf {} \;
	
	tar -cvf $1.tar $1
	
	tar -tvf $1.tar

fi


	
