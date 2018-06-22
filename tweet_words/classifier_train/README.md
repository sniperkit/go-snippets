Classifier Trainer
==================
* Run generate_gob.go to create the classifier.gob file
> go run generate_gob.go 
* To improve quality, store some newline demlimited text into a file named training.txt and execute train.go
> go run train.go
  This updates the input files for the generate_gob.go program.  After this is done, you can execute the generate_gob program to create a better classifier.gob file.


