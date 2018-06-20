![IMAGE](https://raw.github.com/ctava/linearregression-go/master/Data-regression.png)

0. [Install Golang](https://golang.org/doc/install), git, setup `$GOPATH`, and `PATH=$PATH:$GOPATH/bin`
1. go gets
   ```
   go get github.com/gocarina/gocsv
   go get github.com/gonum/plot
   go get github.com/gonum/plot/plotter
   go get github.com/gonum/plot/plotutil
   go get github.com/gonum/plot/vg
   go get github.com/pkg/errors
   go get github.com/sajari/regression
   ```
2. run the source
`go run main.go`
3. review the results
```
N = 1551
Variance observed = 5.138608046119863e+06
Variance Predicted = 2.623877776398886e+06
R2 = 0.5106203378131093

Bobs Data (int=2017): 4401
Bobs Data (int=2018) Prediction: 5707
```
Notes:
This script predicts the amount of data that will be created based on fitting a sample data set to a linear regression model.

data-createdate.csv was created from relational database that simply has a date timestamp of when data was created. the creation of data is - in theory - in a linear relationship between the dates that data has been created and time.

