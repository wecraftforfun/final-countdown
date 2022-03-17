package cmds

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wecraftforfun/final-countdown/helpers"
)

func GoBackToList() tea.Msg {
	return helpers.GoBackToList{}
}

func EnterEditMode() tea.Msg {
	return helpers.EnterEditMode{}
}
