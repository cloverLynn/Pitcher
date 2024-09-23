package utils

import (
	tea "github.com/charmbracelet/bubbletea"
	"main/media"
	"strings"
	"time"
)

// IncrementArray handles selection of previous, current, and next items
func IncrementArray(arr []media.Media, selected int) []int {
	length := len(arr) - 1
	output := make([]int, 3)
	switch selected {
	case length:
		output[0] = selected - 1
		output[1] = selected
		output[2] = 0
	case 0:
		output[0] = length
		output[1] = selected
		output[2] = selected + 1
	default:
		output[0] = selected - 1
		output[1] = selected
		output[2] = selected + 1
	}

	return output
}

type ClearErrorMsg struct{}

func ClearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return ClearErrorMsg{}
	})
}

func ScrubString(s string) string {
	s = strings.Replace(s, " ", "", -1)
	s = strings.ToLower(s)
	return s
}
