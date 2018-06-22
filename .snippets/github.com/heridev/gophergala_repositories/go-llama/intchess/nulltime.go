// this is used to give the ability to NULL times in the database
// source: https://github.com/jinzhu/gorm/issues/10

package intchess

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
// Can scan from both time.Time and NullTime interfaces
func (nt *NullTime) Scan(value interface{}) error {
	if value == nil {
		nt.Valid = false
		return nil
	}

	switch iface := value.(type) {
	case time.Time:
		nt.Time, nt.Valid = iface, true
		return nil
	case NullTime:
		nt.Time, nt.Valid = iface.Time, true //should this be nt.Valid instead of true?
		return nil
	default:
		return nil
	}
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

//this function is used when JSON tries to marshal this struct.
//Without it it will be marshalled into a JSON object with fields Time and Valid.
//This is far more elegant.
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return nt.Time.MarshalJSON()
	} else {
		return json.Marshal(nil)
	}
}

func (nt *NullTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" || string(b) == "" {
		nt.Valid = false
		return nil
	}
	v, err := time.Parse(time.RFC3339, string(b[1:len(b)-1]))
	if err != nil {
		return err
	}
	nt.Time = v
	nt.Valid = true
	return nil
}
