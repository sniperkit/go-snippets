package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"github.com/pkg/errors"
	"github.com/sajari/regression"
)

const companyName = "Bobs"
const startDay = "Jan. 1st"
const startYear = 2013
const startMonth = time.January
const firstPredictedYear = 2017
const numberOfDaysSinceStartOfDataSet = 1461.0
const nextPredictedYear = 2018

var daysLater = float64(numberOfDaysSinceStartOfDataSet + 365)

const entity = "Data"
const dataSourceFileName = "./data-createdate.csv"
const numberOfCharactersInCreateDate = 10

func main() {

	// Aggregate the counts of created entity per day over all days
	counts, err := prepareCountData(dataSourceFileName)
	if err != nil {
		log.Fatal(err)
	}

	// Prepare the "observed" count data for plotting.
	xys := preparePlotData(counts)

	// Perform a regression analysis and print the results for inspection.
	r := performRegression(counts)
	fmt.Printf("Regression:\n%s\n", r)

	// Generate the data for the "observed" and "predicted" plot.
	xysPredicted, err := prepareRegPlotData(r)
	if err != nil {
		log.Fatal(err)
	}

	// Create and save the plot.
	if err = makeRegPlots(xys, xysPredicted); err != nil {
		log.Fatal(err)
	}

	// Make predictions for the number of entities that will
	// be created on numberOfDaysSinceStartOfDataSet from the start of our
	// data set and a year later.
	gcValue1, err := r.Predict([]float64{numberOfDaysSinceStartOfDataSet})
	if err != nil {
		log.Fatal(err)
	}
	gcValue2, err := r.Predict([]float64{float64(daysLater)})
	if err != nil {
		log.Fatal(err)
	}

	// Display the prediction results.
	fmt.Printf("%s %s %s: %d\n", companyName, entity, firstPredictedYear, int(gcValue1))
	fmt.Printf("%s %s %s Prediction: %d\n", companyName, entity, nextPredictedYear, int(gcValue2))
}

// prepareCountData prepares the raw time series data for plotting.
func prepareCountData(dataSet string) ([][]int, error) {

	// Store the daily counts of created entities.
	countMap := make(map[int]int)

	// Get the data set
	data, err := getDataSet(dataSet)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get data from file")
	}

	// Extract the records from the data.
	reader := csv.NewReader(bytes.NewReader(data.Bytes()))
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "Could not read in data records.")
	}

	// Create a map of daily created entities where the keys are the days and
	// the values are the counts of created entities on that day.
	startTime := time.Date(startYear, startMonth, 1, 0, 0, 0, 0, time.UTC)
	layout := "2006-01-02"
	for _, each := range records {
		t, err := time.Parse(layout, each[0][0:numberOfCharactersInCreateDate])
		if err != nil {
			return nil, errors.Wrap(err, "Could not parse timestamps")
		}
		interval := int(t.Sub(startTime).Hours() / 24.0)
		countMap[interval]++
	}

	// Sort the day values which is required for plotting.
	var keys []int
	for k := range countMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	var sortedCounts [][]int
	for _, k := range keys {
		sortedCounts = append(sortedCounts, []int{k, countMap[k]})
	}

	return sortedCounts, nil
}

func getDataSet(fileName string) (bytes.Buffer, error) {
	csvBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		errors.Wrap(err, "Could not read input file")
	}
	var b bytes.Buffer
	b.Write(csvBytes)
	return b, nil
}

// preparePlotData prepares the raw input data for plotting.
func preparePlotData(counts [][]int) plotter.XYs {
	pts := make(plotter.XYs, len(counts))
	var i int

	for _, count := range counts {
		pts[i].X = float64(count[0])
		pts[i].Y = float64(count[1])
		i++
	}

	return pts
}

// preformRegression performs a linear regression of create entities counts vs. day.
func performRegression(counts [][]int) *regression.Regression {
	var r regression.Regression
	r.SetObserved("count of created " + companyName + " " + entity)
	var buffer bytes.Buffer
	buffer.WriteString("days since ")
	buffer.WriteString(startDay)
	buffer.WriteString(strconv.Itoa(startYear))
	r.SetVar(0, buffer.String())

	for _, count := range counts {
		r.Train(regression.DataPoint(
			float64(count[1]),
			[]float64{float64(count[0])}))
	}

	r.Run()
	return &r
}

// prepareRegPlotData prepares predicted point for plotting.
func prepareRegPlotData(r *regression.Regression) (plotter.XYs, error) {
	pts := make(plotter.XYs, int64(daysLater))
	i := 1

	for i <= int(daysLater) {
		pts[i-1].X = float64(i)
		value, err := r.Predict([]float64{float64(i)})
		if err != nil {
			return pts, errors.Wrap(err, "Could not calculate predicted value")
		}
		pts[i-1].Y = value
		i++
	}

	return pts, nil
}

// makeRegPlots makes the second plot including the raw input data and the trained function.
func makeRegPlots(xys1, xys2 plotter.XYs) error {

	// Create a plot value.
	p, err := plot.New()
	if err != nil {
		return errors.Wrap(err, "Could not create plot object")
	}

	// Label the plot.
	var bTitle bytes.Buffer
	bTitle.WriteString("Count of ")
	bTitle.WriteString(companyName)
	bTitle.WriteString(" ")
	bTitle.WriteString(entity)
	bTitle.WriteString(" Created")
	bTitle.WriteString(" Daily")
	p.Title.Text = bTitle.String()
	var bXLabel bytes.Buffer
	bXLabel.WriteString("Days from ")
	bXLabel.WriteString(startDay)
	bXLabel.WriteString(" ")
	bXLabel.WriteString(strconv.Itoa(startYear))
	p.X.Label.Text = bXLabel.String()
	p.Y.Label.Text = "Count"

	// Add both sets of points, predicted and actual, to the plot.
	if err := plotutil.AddLinePoints(p, "Actual", xys1, "Predicted", xys2); err != nil {
		return errors.Wrap(err, "Could not add lines to plot")
	}

	// Save the plot.
	if err := p.Save(7*vg.Inch, 4*vg.Inch, entity+"-regression.png"); err != nil {
		return errors.Wrap(err, "Could not output plot")
	}

	return nil
}
