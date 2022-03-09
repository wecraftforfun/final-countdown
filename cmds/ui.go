package cmds

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mjehanno/timer/helpers"
)

func GoBackToList() tea.Msg {
	return helpers.GoBackToList{}
}

func EnterEditMode() tea.Msg {
	return helpers.EnterEditMode{}
}
