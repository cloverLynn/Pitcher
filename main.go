package main

import (
	"fmt"
	"github.com/76creates/stickers/flexbox"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"os"
	"os/exec"
)

var (
	columnStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).AlignVertical(lipgloss.Center)
	focusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))
)
var style = lipgloss.NewStyle().
	Height(0).Width(100)
var (
	column   = lipgloss.NewStyle().Background(lipgloss.Color("#000000")).Inherit(columnStyle)
	selected = lipgloss.NewStyle().Background(lipgloss.Color("#8B0000")).Inherit(columnStyle)
)

type Media struct {
	name    string
	trailer string
	poster  string
	width   int
}

func (m *Media) Init() tea.Cmd {
	m.width = -1
	return nil
}
func (m *Media) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = 100
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	}
	return m, nil
}

func (m *Media) View() string {
	//width := string(rune(m.width))
	cmd := exec.Command("catimg", "-w", "60", m.poster)
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	poster := out
	name := style.Inline(true).Render(m.name)
	return lipgloss.JoinVertical(lipgloss.Top, name, string(poster))

}

type Model struct {
	selected int
	flexBox  *flexbox.FlexBox
	movies   []Media
	shows    []Show
	err      error
	loaded   bool
	quitting bool
	width    int
}

type Movie struct {
	name string
	path string
}

type Show struct {
	name    string
	picture string
	path    string
	episode []Episode
}

type Episode struct {
	name   string
	path   string
	number int
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			m.InitializeData()
			m.loaded = true
		}
		m.flexBox.SetWidth(msg.Width)
		m.flexBox.SetHeight(msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			m.Scroll(msg.String())
		case "right":
			m.Scroll(msg.String())
		case "enter":
			m.movies[m.selected].PlayTrailer()

		case "ctrl+c", "q":
			return m, tea.Quit
		}

	}
	return m, nil
}

func (m *Model) Scroll(direction string) {
	var inc int
	var rows []*flexbox.Row
	if direction != "left" {
		inc = 1
	} else {
		inc = -1
	}
	m.selected += inc
	if m.selected > len(m.movies)-1 {
		m.selected = 0
	} else if m.selected < 0 {
		m.selected = len(m.movies) - 1
	}
	output := IncrementArray(m.movies, m.selected)
	rows = []*flexbox.Row{
		m.flexBox.NewRow().AddCells(
			flexbox.NewCell(1, 8).SetStyle(column).SetContent(m.movies[output[0]].View()),
			flexbox.NewCell(1, 8).SetStyle(selected).SetContent(m.movies[output[1]].View()),
			flexbox.NewCell(1, 8).SetStyle(column).SetContent(m.movies[output[2]].View()),
		),
	}
	m.flexBox.SetRows(rows)
}

func IncrementArray(arr []Media, selected int) []int {
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

func (m *Model) View() string {
	return m.flexBox.Render()
}

func (m *Media) PlayTrailer() {
	app := "mpv"

	arg0 := "-fs"
	arg1 := m.trailer
	cmd := exec.Command(app, arg0, arg1)
	_, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

}

func (m *Model) InitializeData() {
	m.movies = make([]Media, 0)
	m.movies = append(m.movies, Media{
		name:    "Scream",
		poster:  "./media/scream/poster.jpg",
		trailer: "./media/scream/trailer.mp4",
	})
	m.movies = append(m.movies, Media{
		name:    "The Craft",
		poster:  "./media/craft/poster.jpg",
		trailer: "./media/craft/trailer.mp4",
	})
	m.movies = append(m.movies, Media{
		name:    "Dracula",
		poster:  "./media/dracula/poster.jpg",
		trailer: "./media/dracula/trailer.mp4",
	})
	m.movies = append(m.movies, Media{
		name:    "Silence Of The Lambs",
		poster:  "./media/sotl/poster.jpg",
		trailer: "./media/sotl/trailer.mp4",
	})
	m.selected = 1
	rows := []*flexbox.Row{
		m.flexBox.NewRow().AddCells(
			flexbox.NewCell(1, 6).SetStyle(column).SetContent(m.movies[m.selected-1].View()),
			flexbox.NewCell(1, 6).SetStyle(selected).SetContent(m.movies[m.selected].View()),
			flexbox.NewCell(1, 6).SetStyle(column).SetContent(m.movies[m.selected+1].View()),
		),
	}
	m.flexBox.AddRows(rows)
}

func main() {
	m := Model{
		flexBox: flexbox.New(0, 0),
	}
	p := tea.NewProgram(&m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
