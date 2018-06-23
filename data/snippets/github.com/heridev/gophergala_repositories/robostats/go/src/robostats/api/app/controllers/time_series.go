package controllers

import (
	"strconv"
)

type TimeEvent struct {
	LocalTime int
	Event     map[string]interface{}
}

type TimeEvents []TimeEvent

type timeSeries struct {
	Steps  []string                 `json:"steps"`
	Values map[string][]interface{} `json:"values"`
	Label  map[string]string        `json:"label"`
}

type timeSeriesEnvelope struct {
	TimeSeries timeSeries `json:"time_serie"`
}

func (te *TimeEvents) Envelope() timeSeriesEnvelope {
	ts := timeSeries{}

	teLen := len(*te)

	ts.Steps = make([]string, 0, teLen)
	ts.Values = make(map[string][]interface{})
	ts.Label = make(map[string]string)

	for _, ev := range *te {
		ts.Steps = append(ts.Steps, strconv.Itoa(ev.LocalTime))
		for k := range ev.Event {
			if _, ok := ts.Values[k]; !ok {
				ts.Values[k] = make([]interface{}, 0, teLen)
			}
			if _, ok := ts.Label[k]; !ok {
				ts.Label[k] = k
			}
			ts.Values[k] = append(ts.Values[k], ev.Event[k])
		}
	}

	return timeSeriesEnvelope{ts}
}
