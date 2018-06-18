package toolkit

import (
	"errors"
	"fmt"
	"time"
)

type ResultStatus string

const (
	Status_OK  ResultStatus = "OK"
	Status_NOK ResultStatus = "NOK"
)

type Result struct {
	Status   ResultStatus
	Message  string
	Duration time.Duration
	DurationTxt string
	Data     interface{}

	EncoderID string
	
	time0 time.Time
}

func NewResult() *Result {
	r := new(Result)
	r.Status = Status_OK
	r.time0 = time.Now()
	return r
}

func (r *Result) IsEncoded() bool {
	//fmt.Printf("Encoder ID: %s \n", r.EncoderID)
	return r.EncoderID != ""
}

func (r *Result) SetData(o interface{}) *Result {
	r.Data = o
	r.SetDuration()
	return r
}

func (r *Result) SetBytes(data interface{}, EncoderID string) *Result {
	if EncoderID == "" {
		EncoderID = "json"
	}
	r.EncoderID = EncoderID
	r.Data = ToBytes(data, r.EncoderID)
	//fmt.Println("Encoder ID now is " + r.EncoderID)
	return r
}

func (r *Result) GetFromBytes(out interface{}) error {
	if r.IsEncoded() == false {
		return errors.New("Data is not encoded")
	}
	return FromBytes(r.Data.([]byte), r.EncoderID, out)
}

func (r *Result) SetDuration(){
	r.Duration = time.Since(r.time0)
	r.DurationTxt = r.Duration.String()
}

func (r *Result) SetError(e error) *Result {
	r.Status = Status_NOK
	r.Message = e.Error()
	r.SetDuration()
	return r
}

func (r *Result) SetErrorTxt(e string) *Result {
	r.Status = Status_NOK
	r.Message = e
	r.SetDuration()
	return r
}

func (r *Result) Error() error {
	var e error
	if r.Status == Status_NOK {
		e = errors.New(r.Message)
	}
	return e
}

func (r *Result) Cast(out interface{}, method string) error {
	if method == "" {
		method = "json"
	}

	if r.Data == nil {
		return errors.New("Data is nil")
	}

	if method == "json" {
		bs := Jsonify(r.Data)
		e := Unjson(bs, out)
		if e != nil {
			return errors.New("Can not decode data. " + e.Error())
		}
		return nil
	}

	return errors.New("Unable to cast due to unknown cast method")
}

func (a *Result) Run(f func(data interface{}) (interface{}, error), parm interface{}) *Result {
	t0 := time.Now()
	a.Status = Status_OK
	a.Message = ""
	if f == nil {
		a.Data = nil
	} else {
		r, e := f(parm)
		if e != nil {
			a.Status = Status_NOK
			a.Message = e.Error()
			a.Data = nil
		} else {
			a.Data = r
		}
	}
	a.Duration = time.Since(t0)
	return a
}

func CallResult(url, calltype string, data []byte) (*Result, error) {
	r, e := HttpCall(url, calltype, data, M{}.Set("expectedstatus", 200))
	if e != nil {
		return nil, fmt.Errorf(url + " Call eror: " + e.Error())
	}

	result := NewResult()
	edecode := Unjson(HttpContent(r), &result)
	if edecode != nil {
		return nil, fmt.Errorf(url + " Decode error: " + edecode.Error())
	}
	if result.Status == Status_NOK {
		return result, fmt.Errorf(result.Message)
	}
	return result, nil
}
