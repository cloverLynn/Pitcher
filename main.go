package main

import (
	"fmt"
	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/bubbletea"
	components "main/Components"
	"main/media"
	"os"
)

type Model struct {
	FlexBoxComponent components.FlexBoxComponent
	Loaded           bool
	Quitting         bool
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.Loaded {
			m.InitializeData(msg.Width, msg.Height)
			m.Loaded = true
		}
		m.FlexBoxComponent.FlexBox.SetWidth(msg.Width)
		m.FlexBoxComponent.FlexBox.SetHeight(msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "left", "right":
			m.FlexBoxComponent.Scroll(msg.String())
		case "enter":
			m.FlexBoxComponent.Movies[m.FlexBoxComponent.Selected].PlayTrailer()
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *Model) InitializeData(width, height int) {
	m.FlexBoxComponent.Init(width, height)
	movies := []media.Media{
		{Name: "Scream", Poster: "./data/scream/poster.jpg", Trailer: "./data/scream/trailer.mp4"},
		{Name: "The Craft", Poster: "./data/craft/poster.jpg", Trailer: "./data/craft/trailer.mp4"},
		{Name: "Dracula", Poster: "./data/dracula/poster.jpg", Trailer: "./data/dracula/trailer.mp4"},
		{Name: "Silence of The Lambs", Poster: "./data/sotl/poster.jpg", Trailer: "./data/sotl/trailer.mp4"},
	}
	m.FlexBoxComponent.Movies = movies
	m.FlexBoxComponent.Selected = 1

	rows := []*flexbox.Row{
		m.FlexBoxComponent.FlexBox.NewRow().AddCells(
			flexbox.NewCell(1, 6).SetContent(m.FlexBoxComponent.Movies[m.FlexBoxComponent.Selected-1].View()),
			flexbox.NewCell(1, 6).SetContent(m.FlexBoxComponent.Movies[m.FlexBoxComponent.Selected].View()),
			flexbox.NewCell(1, 6).SetContent(m.FlexBoxComponent.Movies[m.FlexBoxComponent.Selected+1].View()),
		),
	}
	m.FlexBoxComponent.FlexBox.AddRows(rows)
}

func (m *Model) View() string {
	return m.FlexBoxComponent.View()
}

func main() {
	m := &Model{}
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
