package components

import (
	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/lipgloss"
	"main/media"
	"main/utils"
)

var (
	column   = lipgloss.NewStyle().Background(lipgloss.Color("#000000")).Align(lipgloss.Center)
	selected = lipgloss.NewStyle().Background(lipgloss.Color("#8B0000")).Align(lipgloss.Center)
)

type FlexBoxComponent struct {
	FlexBox  *flexbox.FlexBox
	Movies   []media.Media
	Selected int
}

// Initialize FlexBox component
func (f *FlexBoxComponent) Init(width, height int) {
	f.FlexBox = flexbox.New(width, height)
}

// Handle scrolling logic
func (f *FlexBoxComponent) Scroll(direction string) {
	var inc int
	if direction == "left" {
		inc = -1
	} else {
		inc = 1
	}

	f.Selected += inc
	if f.Selected > len(f.Movies)-1 {
		f.Selected = 0
	} else if f.Selected < 0 {
		f.Selected = len(f.Movies) - 1
	}

	output := utils.IncrementArray(f.Movies, f.Selected)
	rows := []*flexbox.Row{
		f.FlexBox.NewRow().AddCells(
			flexbox.NewCell(1, 8).SetStyle(column).SetContent(f.Movies[output[0]].View()),
			flexbox.NewCell(1, 8).SetStyle(selected).SetContent(f.Movies[output[1]].View()),
			flexbox.NewCell(1, 8).SetStyle(column).SetContent(f.Movies[output[2]].View()),
		),
	}
	f.FlexBox.SetRows(rows)
}

// Render FlexBox
func (f *FlexBoxComponent) View() string {
	if f.FlexBox == nil {
		return ""
	}
	return f.FlexBox.Render()
}
