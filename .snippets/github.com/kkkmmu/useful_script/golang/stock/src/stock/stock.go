package stock

import (
	"fmt"
	"history"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"util"
)

type Stock struct {
	Code string
	Name string

	History *history.History
}

func (s *Stock) String() string {
	var res string
	res += fmt.Sprintf("%-10s %8s:\n", s.Name, s.Code)

	res += fmt.Sprintf("   History Summary:\n")
	d, err := s.History.GetHighest()
	if err != nil {
		res += fmt.Sprintf("\t NOHIGHEST NODATE")
		return res
	} else {
		res += fmt.Sprintf("\t %15f %4d/%02d/%02d %10s", d.High, d.Date.Year(), d.Date.Month(), d.Date.Day(), d.Date.Weekday())
	}

	d, err = s.History.GetLowest()
	if err != nil {
		res += fmt.Sprintf(" NOLOWEST NODATE")
		return res
	} else {
		res += fmt.Sprintf(" %15f %4d/%02d/%02d %10s", d.Low, d.Date.Year(), d.Date.Month(), d.Date.Day(), d.Date.Weekday())
	}
	res += fmt.Sprintf("\n")

	res += fmt.Sprintf("   Year Summary:\n")
	for y, ye := range s.History.Year {
		res += fmt.Sprintf("{%d} ", y)
		d, err = ye.GetHighest()
		if err != nil {
			res += fmt.Sprintf("\t NOHIGHEST NODATE")
			return res
		} else {
			res += fmt.Sprintf("\t %15f %4d/%02d/%02d %10s", d.High, d.Date.Year(), d.Date.Month(), d.Date.Day(), d.Date.Weekday())
		}

		d, err = ye.GetLowest()
		if err != nil {
			res += fmt.Sprintf(" NOLOWEST NODATE")
			return res
		} else {
			res += fmt.Sprintf(" %15f %4d/%02d/%02d %10s", d.Low, d.Date.Year(), d.Date.Month(), d.Date.Day(), d.Date.Weekday())
		}

		res += fmt.Sprintf("\n")
	}
	/*
		fmt.Println(time.Now().Year())
		y, err := s.History.GetYear(time.Now().Year())
		if err != nil {
			res += fmt.Sprintf("\t NOYEARE INFO")
			return res
		}
		d, err = y.GetHighest()
		if err != nil {
			res += fmt.Sprintf("\t NOHIGHEST NODATE")
			return res
		} else {
			res += fmt.Sprintf("\t %15f %4d/%02d/%02d %10s", d.High, d.Date.Year(), d.Date.Month(), d.Date.Day(), d.Date.Weekday())
		}

		d, err = s.History.GetLowest()
		if err != nil {
			res += fmt.Sprintf(" NOLOWEST NODATE")
			return res
		} else {
			res += fmt.Sprintf(" %15f %4d/%02d/%02d %10s", d.Low, d.Date.Year(), d.Date.Month(), d.Date.Day(), d.Date.Weekday())
		}
	*/

	return res
}

var HexunDailyK = "http://flashquote.stock.hexun.com/Quotejs/DA/1_%s_DA.html?"
var StockCodeR = regexp.MustCompile("Code:[[:space:]]+(?P<code>[[:alnum:]]{6})[[:space:]]+Name")
var DailyR = regexp.MustCompile(`\[(?P<date>[[:digit:]]{8}),(?P<p1>[[:digit:].]+),(?P<p2>[[:digit:].]+),(?P<p3>[[:digit:].]+),(?P<p4>[[:digit:].]+),(?P<p5>[[:digit:].]+),(?P<v1>[[:digit:]]+),(?P<v2>[[:digit:]]+)\]`)
var StockListURL = "http://quote.eastmoney.com/stocklist.html"

func New(code, name string) (*Stock, error) {
	ns := &Stock{
		Code: code,
		Name: name,
	}

	err := ns.Init()
	if err != nil {
		return nil, err
	}

	return ns, nil
}

func NewFromLocal(code, name string) (*Stock, error) {
	ns := &Stock{
		Code: code,
		Name: name,
	}

	err := ns.InitFromLocal()
	if err != nil {
		return nil, err
	}

	return ns, nil
}

func (s *Stock) Init() error {
	h, err := s.GetHistory()
	if err != nil {
		return fmt.Errorf("Cannot init stock %s with: %s", s.Code, err)
	}

	return h.Init()
}

func (s *Stock) InitFromLocal() error {
	h, err := s.GetHistoryFromLocal()
	if err != nil {
		return fmt.Errorf("Cannot init From local stock %s with: %s", s.Code, err)
	}

	return h.Init()
}

func (s *Stock) GetHistory() (*history.History, error) {
	if s.History != nil {
		return s.History, nil
	}

	h, err := s.GetHistoryFromLocal()
	if err != nil {
		h, err = s.GetHistoryFromNetwork()
		if err != nil {
			return nil, fmt.Errorf("Cannot get hisotry for %s with %s", s.Code, err)
		}
	}

	s.History = h

	return h, nil
}

func (s *Stock) GetHistoryFromNetwork() (*history.History, error) {
	url := fmt.Sprintf(HexunDailyK, s.Code)

	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	util.SaveToFile("asset/"+s.Code, body)

	if !IsHistoryFileValid("asset/" + s.Code) {
		os.Remove("asset/" + s.Code)
		return nil, fmt.Errorf("History info from hexun is invalid")
	}

	return s.GetHistoryFromLocal()
}

func (s *Stock) GetHistoryFromLocal() (*history.History, error) {
	content, err := ioutil.ReadFile("asset/" + s.Code)
	if err != nil {
		return nil, err
	}

	days := DailyR.FindAllStringSubmatch(string(content), -1)

	dis := make([]*history.Daily, 0, len(days))
	for _, day := range days {
		di, err := GetDailyInfoFromString(day)
		if err != nil {
			panic(err)
		}
		dis = append(dis, di)
	}

	h := &history.History{
		Info: dis,
	}

	s.History = h

	return s.History, nil
}

func GetDailyInfoFromString(day []string) (*history.Daily, error) {
	t, err := GetTimeFromString(day[1])
	if err != nil {
		return nil, fmt.Errorf("Invalid daily info: %+v", day)
	}

	//[19970704,0.00,6.66,7.45,6.66,6.85,53921321,375882768]
	dopen, err := strconv.ParseFloat(day[2], 62)
	if err != nil {
		return nil, fmt.Errorf("Invalid daily info: %+v: %s", day, err)
	}

	dclose, err := strconv.ParseFloat(day[3], 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid daily info: %+v: %s", day, err)
	}

	dhigh, err := strconv.ParseFloat(day[4], 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid daily info: %+v: %s", day, err)
	}

	dlow, err := strconv.ParseFloat(day[5], 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid daily info: %+v: %s", day, err)
	}

	deva, err := strconv.ParseFloat(day[6], 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid daily info: %+v: %s", day, err)
	}

	dchange, err := strconv.ParseInt(day[7], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid daily info: %+v: %s", day, err)
	}

	dtotal, err := strconv.ParseInt(day[8], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Invalid daily info: %+v: %s", day, err)
	}

	return &history.Daily{
		Date:     t,
		Open:     dopen,
		Close:    dclose,
		High:     dhigh,
		Low:      dlow,
		Everage:  deva,
		Exchange: dchange,
		Total:    dtotal,
	}, nil
}

func IsHistoryFileValid(name string) bool {
	content, err := ioutil.ReadFile(name)
	if err != nil {
		return false
	}

	if strings.Contains(string(content), "DOCTYPE") {
		return false
	}

	return true
}

func GetCurrentDateString() string {
	y, m, d := time.Now().Date()
	return fmt.Sprintf("%d%02d%02d", y, m, d)
}

func GetTimeFromString(str string) (*time.Time, error) {
	if len(str) != 8 {
		return nil, fmt.Errorf("Invalid date: %s", str)
	}

	year, err := strconv.ParseInt(str[:4], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Invalid date: %s", str)
	}

	month, err := strconv.ParseInt(str[4:6], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Invalid date: %s", str)
	}

	date, err := strconv.ParseInt(str[6:8], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Invalid date: %s", str)
	}

	t := time.Date(int(year), time.Month(int(month)), int(date), 0, 0, 0, 0, time.UTC)

	return &t, nil
}
