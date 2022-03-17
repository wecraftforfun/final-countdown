package cmds

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wecraftforfun/final-countdown/helpers"
	"github.com/wecraftforfun/final-countdown/models"
)

func readFromFile() ([]models.CountDown, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(path.Join(currentUser.HomeDir, ".countdown.json")); err != nil {
		return nil, nil
	}

	content, err := ioutil.ReadFile(path.Join(currentUser.HomeDir, ".countdown.json"))
	if err != nil {
		return nil, err
	}
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

func SaveListCmd(list []models.CountDown) tea.Cmd {
	return func() tea.Msg {
		if err := writeToFile(list); err != nil {
			return helpers.ErrorMsg{
				Err: err,
			}
		}
		return nil
	}
}

func writeToFile(list []models.CountDown) error {
	currentUser, err := user.Current()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path.Join(currentUser.HomeDir, ".countdown.json"), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	data, err := json.Marshal(list)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	err = file.Sync()
	if err != nil {
		return err
	}
	return nil
}
