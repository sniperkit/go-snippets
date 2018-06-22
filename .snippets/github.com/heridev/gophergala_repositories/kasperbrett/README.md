# Kasperbrett
A dead simple dashboard app created in the course of the [Gopher Gala](http://gophergala.com/).

## Demo

A demo can be found [here](http://mycloset.dlinkddns.com:8080/).

## Features

* Persistent storage of all time series data
* Pluggable data source architecture
* Web-based UI with realtime updates
* Data is accessible via REST API
* Minimally scriptable

## Limitations

* Due to Gopher Gala's time limit there is currently only one data source implemented (URL Scraper).


## Use Cases

* Metrics and analytics for all kind of data accessible on the Web


## Data Sources

* URL Scraper




## Caveats
* Time is limited at the Gopher Gala and I don't want to deal with managing own packages. That's why all the Go code for Kasperbrett is in one package (and probably even within one file). Code quality might also be very rough because the goal is to create a running prototype.