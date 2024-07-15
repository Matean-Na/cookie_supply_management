package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/guregu/null.v4/zero"
	"time"
)

const DATE_FORMAT = "2006-01-02"

type DateOnly struct {
	Date zero.Time
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (d *DateOnly) Scan(value interface{}) error {
	if value == nil {
		*d = DateOnly{Date: zero.Time{}}
		return nil
	}

	time, ok := value.(time.Time)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	*d = DateOnly{
		Date: zero.TimeFrom(time),
	}
	return nil
}

// Value return json value, implement driver.Valuer interface
func (d DateOnly) Value() (driver.Value, error) {
	return d.Date, nil
}

func (d DateOnly) ValueStr() string {
	return d.Date.Time.Format("02.01.2006")
}

func (d DateOnly) Day() int {
	return d.Date.Time.Day()
}

func (d DateOnly) Month() time.Month {
	return d.Date.Time.Month()
}

func (d DateOnly) Year() int {
	return d.Date.Time.Year()
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	var value interface{}
	if d.Date.IsZero() {
		value = ""
	} else {
		value = d.Date.Time.Format(DATE_FORMAT)
	}
	date, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return date, nil
}

// UnmarshalJSON to deserialize []byte
func (d *DateOnly) UnmarshalJSON(b []byte) error {
	var result string
	if err := json.Unmarshal(b, &result); err != nil {
		return err
	}
	if result != "" {
		date, err := time.Parse(DATE_FORMAT, result)
		if err != nil {
			return err
		}
		*d = DateOnly{
			Date: zero.TimeFrom(date),
		}
	} else {
		*d = DateOnly{
			Date: zero.Time{},
		}
	}
	return nil
}

func (d *DateOnly) FromString(s string) error {
	parsed, err := time.Parse(DATE_FORMAT, s)
	if err != nil {
		return err
	}

	*d = DateOnly{
		Date: zero.TimeFrom(parsed),
	}

	return nil
}
