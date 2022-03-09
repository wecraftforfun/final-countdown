package helpers

import (
	"fmt"

	"github.com/mjehanno/timer/models"
)

type ErrorMsg struct {
	Err error
}

type UpdateListMsg struct {
	Countdowns []models.CountDown
}

type GoBackToList struct{}

type EnterEditMode struct{}

func (e ErrorMsg) Error() string {
	return fmt.Sprintf("Error happened : %s", e.Err)
}
