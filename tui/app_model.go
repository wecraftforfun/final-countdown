package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wecraftforfun/final-countdown/cmds"
	"github.com/wecraftforfun/final-countdown/helpers"
	"github.com/wecraftforfun/final-countdown/models"
)

const timeout = time.Hour * 5

var (
	focusedStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle        = focusedStyle.Copy()
	noStyle            = lipgloss.NewStyle()
	statusMessageStyle = focusedStyle.Copy().Margin(2)
)

type AppModel struct {
	Choices      []models.CountDown
	FocusIndex   int // items on the timer list
	Cursor       int // which timer list item our cursor is pointing at
	Timer        timer.Model
	Inputs       []textinput.Model
	Keys         keyMap
	Help         help.Model
	List         list.Model
	IsInsertMode bool
}

type keyMap struct {
	insertItem key.Binding
	deleteItem key.Binding
	cancel     key.Binding
	quit       key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.insertItem, k.deleteItem, k.cancel, k.quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.insertItem, k.deleteItem, k.cancel, k.quit},
	}
}

func newListKeyMap() keyMap {
	return keyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "Add a new countdown"),
		),
		deleteItem: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "Delete current countdown"),
		),
		cancel: key.NewBinding(
			key.WithKeys("end"),
			key.WithHelp("end", "Return to list"),
			key.WithDisabled(),
		),
		quit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "Quit app"),
		),
	}
}

func InitialModel() AppModel {
	m := AppModel{
		Keys:   newListKeyMap(),
		Help:   help.New(),
		Inputs: make([]textinput.Model, 3),
		List:   list.New(nil, helpers.NewListDelegate(), 1000, 15),
		Timer:  timer.New(0),
	}

	m.List.Styles.StatusBar = statusMessageStyle
	m.List.SetShowHelp(false)
	m.List.KeyMap = list.DefaultKeyMap()
	return m
}

func (m AppModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return cmds.InitListCmd
}

func createForms(m AppModel) []textinput.Model {
	result := []textinput.Model{}
	var t textinput.Model
	for i := range m.Inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32
		switch i {
		case 0:
			t.Placeholder = "Title"
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Description"
			t.CharLimit = 64
		case 2:
			t.Placeholder = "Due Date as YYYY/MM/DD-hh:mm:ss"
		}
		result = append(result, t)
	}
	return result
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case helpers.UpdateListMsg:
		for _, v := range msg.Countdowns {
			m.List.InsertItem(len(m.List.Items()), list.Item(v))
			m.Choices = append(m.Choices, v)
		}
	case helpers.ErrorMsg:
		m.List.NewStatusMessage(msg.Error())
	case helpers.GoBackToList:
		m.IsInsertMode = false
		m.FocusIndex = 0
		m.Cursor = 0
		m.Inputs = make([]textinput.Model, 3)
		m.Keys.cancel.SetEnabled(false)
		m.Keys.insertItem.SetEnabled(true)
		m.Keys.deleteItem.SetEnabled(true)
	case helpers.EnterEditMode:
		m.IsInsertMode = true
		m.Keys.cancel.SetEnabled(true)
		m.Keys.insertItem.SetEnabled(false)
		m.Keys.deleteItem.SetEnabled(false)
		m.Inputs = createForms(m)
		cmd := m.updateInputs(msg)
		m.Inputs[0].Focus()
		return tea.Model(m), cmd
	case timer.TickMsg:
		var cmd tea.Cmd
		m.Timer, cmd = m.Timer.Update(msg)
		return m, cmd
	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.Timer, cmd = m.Timer.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.insertItem) && !m.IsInsertMode:
			// Entering creation mode to add item
			return tea.Model(m), tea.Batch(cmds.EnterEditMode)
		case key.Matches(msg, m.Keys.deleteItem):
			m.Choices = append(m.Choices[:m.List.Cursor()], m.Choices[m.List.Cursor()+1:]...)
			m.List.RemoveItem(m.List.Cursor())
			return tea.Model(m), tea.Batch(cmds.SaveListCmd(m.Choices))
		case !m.IsInsertMode && (msg.String() == "up" || msg.String() == "down"):
			// Moving through the list
			if msg.String() == "up" {
				m.List.CursorUp()
			} else {
				m.List.CursorDown()
			}
		case m.IsInsertMode && (msg.String() == "up" || msg.String() == "down"):
			// Moving through the displayed inputs
			if msg.String() == "up" {
				m.FocusIndex--
			} else {
				m.FocusIndex++
			}
			cmds := make([]tea.Cmd, len(m.Inputs))
			for i := 0; i <= len(m.Inputs)-1; i++ {
				if i == m.FocusIndex {
					// Set focused state
					cmds[i] = m.Inputs[i].Focus()
					m.Inputs[i].PromptStyle = focusedStyle
					m.Inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.Inputs[i].Blur()
				m.Inputs[i].PromptStyle = noStyle
				m.Inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		case m.IsInsertMode && msg.String() == "enter":
			d, err := time.Parse("2006/01/02-15:04:05", m.Inputs[2].Value())
			if err != nil {
				return tea.Model(m), tea.Batch(func() tea.Msg {
					return helpers.ErrorMsg{
						Err: err,
					}
				})
			}
			c := models.CountDown{
				Name:    m.Inputs[0].Value(),
				Desc:    m.Inputs[1].Value(),
				DueDate: models.CustomTime{Time: d},
			}
			m.Choices = append(m.Choices, c)
			m.List.InsertItem(len(m.Choices), c)
			return tea.Model(m), tea.Batch(cmds.SaveListCmd(m.Choices), cmds.GoBackToList)
		case !m.IsInsertMode && msg.String() == "enter":
			m.Timer = timer.New(time.Until(m.Choices[m.List.Cursor()].DueDate.Time))
			return tea.Model(m), m.Timer.Start()
		case key.Matches(msg, m.Keys.quit):
			return tea.Model(m), tea.Quit
		case key.Matches(msg, m.Keys.cancel) && m.IsInsertMode:
			return tea.Model(m), cmds.GoBackToList
		}
	}

	cmd := m.updateInputs(msg)

	return tea.Model(m), cmd
}

func (m *AppModel) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.Inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m AppModel) View() string {
	m.List.Title = "Which timer do you want to look at ?"
	var s string
	//s :=
	helpView := m.Help.View(m.Keys)
	if m.IsInsertMode {
		for _, input := range m.Inputs {
			s += "\n"
			s += input.View()
		}
	} else {
		s += m.List.View()
		if m.Timer.Timeout != time.Hour*5 {
			days := int(m.Timer.Timeout.Seconds()) / (24 * 3600)
			hours := (int(m.Timer.Timeout.Seconds()) - days*(24*3600)) / 3600
			min := (int(m.Timer.Timeout.Seconds()) - days*(24*3600) - hours*3600) / 60
			sec := (int(m.Timer.Timeout.Seconds()) - days*(24*3600) - hours*3600 - min*60)
			s += fmt.Sprintf("%vd%vh%vm%vs", days, hours, min, sec)
		}
	}
	return s + "\n" + helpView
}
