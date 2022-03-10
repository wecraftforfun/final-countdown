package models

import (
	"encoding/json"
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

func (t *CustomTime) UnmarshalJSON(buf []byte) error {
	tt, err := time.Parse("2006/01/02-15:04:05", strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

func (t *CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Time.Format("2006/01/02-15:04:05") + `"`), nil
}

type CountDown struct {
	Name    string     `json:"name"`
	Desc    string     `json:"description"`
	DueDate CustomTime `json:"dueDate"`
}

func (c CountDown) String() string {
	js, _ := json.Marshal(c)
	return string(js)
}

func (c CountDown) FilterValue() string {
	return c.Name
}
func (c CountDown) Title() string {
	return c.Name
}

func (c CountDown) Description() string {
	return c.Desc
}
