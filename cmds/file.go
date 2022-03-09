package cmds

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mjehanno/timer/helpers"
	"github.com/mjehanno/timer/models"
)

func readFromFile() ([]models.CountDown, error) {
	content, err := ioutil.ReadFile(path.Clean("~/.countdown.json"))
	if err != nil {
		return nil, err
	}
	fmt.Println(string(content))
	result := make([]models.CountDown, 0)
	json.Unmarshal(content, &result)
	return result, nil
}

func InitListCmd() tea.Msg {
	if countdowns, err := readFromFile(); err != nil {
		return helpers.ErrorMsg{
			Err: err,
		}
	} else {
		return helpers.UpdateListMsg{
			Countdowns: countdowns,
		}
	}
}

func writeToFile([]models.CountDown) error {
	file, err := os.Open(path.Clean("~/.countdown.json"))
	if err != nil {
		return err
	}

	err = file.Sync()
	if err != nil {
		return err
	}

	return nil
}
