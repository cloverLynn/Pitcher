package main

import (
	"fmt"
	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/bubbletea"
	huh "github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	components "main/Components"
	"main/media"
	"os"
)

type status int

var models []tea.Model

const (
	movies status = iota
	form
)

var (
	column   = lipgloss.NewStyle().Background(lipgloss.Color("#000000")).Align(lipgloss.Center)
	selected = lipgloss.NewStyle().Background(lipgloss.Color("#8B0000")).Align(lipgloss.Center)
)

type Model struct {
	FlexBoxComponent components.FlexBoxComponent
	Loaded           bool
	Quitting         bool
}

func (m *Model) Init() tea.Cmd { return nil }

func New() *Model { return &Model{} }

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
		case "ctrl+n":
			models[movies] = m
			newForm := NewForm()
			models[form] = newForm
			return models[form].Update("First")
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
			flexbox.NewCell(1, 6).SetStyle(column).SetContent(m.FlexBoxComponent.Movies[m.FlexBoxComponent.Selected-1].View()),
			flexbox.NewCell(1, 6).SetStyle(selected).SetContent(m.FlexBoxComponent.Movies[m.FlexBoxComponent.Selected].View()),
			flexbox.NewCell(1, 6).SetStyle(column).SetContent(m.FlexBoxComponent.Movies[m.FlexBoxComponent.Selected+1].View()),
		),
	}
	m.FlexBoxComponent.FlexBox.AddRows(rows)
}

func (m *Model) View() string {
	return m.FlexBoxComponent.View()
}

type Form struct {
	form    *huh.Form // huh.Form is just a tea.Model
	loaded  bool
	name    string
	poster  string
	trailer string
}

func NewForm() Form {
	path, _ := os.UserHomeDir()
	return Form{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("name").
					Title("Media Name"),
				huh.NewFilePicker().CurrentDirectory(path).
					Title("Poster").
					Key("poster"),
				huh.NewFilePicker().CurrentDirectory(path).
					Title("Trailer").
					Key("trailer"),
			)),
	}
}

func (t Form) Init() tea.Cmd {
	return t.form.Init()
}

func (t Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !t.loaded {
		t.form.NextField()
		t.form.PrevField()
		t.loaded = true
	}
	form, cmd := t.form.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return models[movies].Update(nil)
		case "enter":
		}
	}
	if f, ok := form.(*huh.Form); ok {
		t.form = f
	}

	return t, cmd
}

func (t Form) View() string {
	if t.form.State == huh.StateCompleted {
		name := t.form.GetString("name")
		poster := t.form.GetString("poster")
		trailer := t.form.GetString("trailer")
		return fmt.Sprintf("You selected: %s, Lvl. %s %s", name, poster, trailer)
	}
	return t.form.View()
}

func main() {
	models = []tea.Model{New(), NewForm()}
	m := models[movies]
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
