package models

import (
	"encoding/json"
	"time"
)

type CountDown struct {
	Name    string    `json:"name"`
	Desc    string    `json:"description"`
	DueDate time.Time `json:"dueDate"`
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
