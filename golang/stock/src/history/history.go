package history

import (
	"errors"
	"fmt"
	"time"
)

var ErrorAPINotImplemented = errors.New("The API is not implemented")

type Daily struct {
	Date     *time.Time
	High     float64
	Low      float64
	Open     float64
	Close    float64
	Everage  float64
	Exchange int64
	Total    int64
}

type History struct {
	Info []*Daily
	Year map[int]*Year
}

type Year struct {
	Year  int
	Info  []*Daily
	Month map[time.Month]*Month
}

type Month struct {
	Month time.Month
	Info  []*Daily
}

func (hi History) GetHighest() (*Daily, error) {
	if len(hi.Info) == 0 {
		return nil, fmt.Errorf("Invalid History info")
	}

	var di *Daily
	for i, h := range hi.Info {
		if i == 0 {
			di = h
		}

		if h.High > di.High {
			di = h
		}
	}

	return di, nil
}

func (hi History) GetLowest() (*Daily, error) {
	if len(hi.Info) == 0 {
		return nil, fmt.Errorf("Invalid History info")
	}

	var di *Daily
	for i, h := range hi.Info {
		if i == 0 {
			di = h
		}

		if h.High < di.High {
			di = h
		}
	}

	return di, nil
}

func (y *Year) GetHighest() (*Daily, error) {
	if len(y.Info) == 0 {
		return nil, fmt.Errorf("Invalid History info")
	}

	var di *Daily
	for i, h := range y.Info {
		if i == 0 {
			di = h
		}

		if h.High > di.High {
			di = h
		}
	}

	return di, nil
}

func (y *Year) GetLowest() (*Daily, error) {
	if len(y.Info) == 0 {
		return nil, fmt.Errorf("Invalid History info")
	}

	var di *Daily
	for i, h := range y.Info {
		if i == 0 {
			di = h
		}

		if h.High < di.High {
			di = h
		}
	}

	return di, nil
}

func (m *Month) GetHighest() (*Daily, error) {
	if len(m.Info) == 0 {
		return nil, fmt.Errorf("Invalid History info")
	}

	var di *Daily
	for i, h := range m.Info {
		if i == 0 {
			di = h
		}

		if h.High > di.High {
			di = h
		}
	}

	return di, nil
}

func (m *Month) GetLowest() (*Daily, error) {
	if len(m.Info) == 0 {
		return nil, fmt.Errorf("Invalid History info")
	}

	var di *Daily
	for i, h := range m.Info {
		if i == 0 {
			di = h
		}

		if h.High < di.High {
			di = h
		}
	}

	return di, nil
}

func (h *History) GetYear(year int) (*Year, error) {
	if h.Year == nil {
		return nil, fmt.Errorf("Year database is not initialized")
	}

	y, ok := h.Year[year]
	if !ok {
		return nil, fmt.Errorf("Entry not exist for year: %d", year)
	}

	return y, nil
}

func (h *History) SplitByYear() (map[int]*Year, error) {
	if h.Info == nil {
		return nil, fmt.Errorf("History database is not intialized")
	}

	years := make(map[int]*Year, 10)
	for _, d := range h.Info {
		if _, ok := years[d.Date.Year()]; !ok {
			years[d.Date.Year()] = &Year{
				Year: d.Date.Year(),
				Info: make([]*Daily, 0, 100),
			}
		}

		years[d.Date.Year()].Info = append(years[d.Date.Year()].Info, d)
	}

	return years, nil
}

func (h *History) Parse() error {
	ys, err := h.SplitByYear()
	if err != nil {
		return fmt.Errorf("Cannot parse years info with: %s", err)
	}

	for _, y := range ys {
		err = y.ParseMonth()
		if err != nil {
			return fmt.Errorf("Cannot parse month of year %d with %s", y.Year, err)
		}
	}

	h.Year = ys

	return nil
}

func (h *History) Init() error {
	return h.Parse()
}

func (y *Year) SplitByMonth() (map[time.Month]*Month, error) {
	if y.Info == nil {
		return nil, fmt.Errorf("Year database is not intialized")
	}

	months := make(map[time.Month]*Month, 10)
	for _, d := range y.Info {
		if _, ok := months[d.Date.Month()]; !ok {
			months[d.Date.Month()] = &Month{
				Month: d.Date.Month(),
				Info:  make([]*Daily, 0, 31),
			}
		}

		months[d.Date.Month()].Info = append(months[d.Date.Month()].Info, d)
	}

	return months, nil
}

func (y *Year) ParseMonth() error {
	ms, err := y.SplitByMonth()
	if err != nil {
		return fmt.Errorf("Cannot parse month info with: %s", err)
	}

	y.Month = ms

	return nil
}

func (y Year) GetMonth(month time.Month) (*Month, error) {
	if y.Month == nil {
		return nil, fmt.Errorf("Month database is not initialized")
	}

	m, ok := y.Month[month]
	if !ok {
		return nil, fmt.Errorf("Entry not exist for month: %d", month)
	}

	return m, nil
}
