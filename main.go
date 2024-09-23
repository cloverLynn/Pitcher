package main

import (
	"encoding/csv"
	"fmt"
	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/bubbletea"
	huh "github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	components "main/Components"
	"main/media"
	"main/utils"
	"os"
	"os/exec"
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
	movies := []media.Media{}
	file, err := os.Open("./data/data.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading records")
	}

	// Loop to iterate through
	// and print each of the string slice
	for _, eachrecord := range records {
		movies = append(movies, media.Media{
			Name:    eachrecord[0],
			Poster:  eachrecord[1],
			Trailer: eachrecord[2],
		})
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
	file.Close()
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
				huh.NewFilePicker().AllowedTypes([]string{".jpg"}).CurrentDirectory(path).
					Title("Poster").
					Key("poster"),
				huh.NewFilePicker().AllowedTypes([]string{".mp4"}).CurrentDirectory(path).
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
	if t.form.State == huh.StateCompleted {
		name := t.form.GetString("name")
		poster := t.form.GetString("poster")
		trailer := t.form.GetString("trailer")
		newDir := "./data/" + utils.ScrubString(name)
		err := os.Mkdir(newDir, 0755)
		command := exec.Command("cp", poster, newDir+"/poster.jpg")
		poster = newDir + "/poster.jpg"
		command.Run()
		command = exec.Command("cp", trailer, newDir+"/trailer.mp4")
		trailer = newDir + "/trailer.mp4"
		command.Run()
		if err != nil {
			panic(err)
		}
		file, err := os.OpenFile("./data/data.csv", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		record := "\n" + name + "," + poster + "," + trailer + "\n"
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {

			}
		}(file)
		_, err = file.WriteString(record)

		err = file.Close()
		if err != nil {
			return nil, nil
		}

	}

	return t, cmd
}

func (t Form) View() string {
	if t.form.State == huh.StateCompleted {
		return fmt.Sprintf("Movie has been saved")
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
